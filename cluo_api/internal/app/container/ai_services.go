package container

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	ollamaClient "github.com/hengadev/cluo_api/internal/infrastructure/ollama"
	aiChatRepo "github.com/hengadev/cluo_api/internal/infrastructure/postgres/ai_chat"
	aiTranscriptAnalysisRepo "github.com/hengadev/cluo_api/internal/infrastructure/postgres/ai_transcript_analysis"
	aiTranscriptionJobRepo "github.com/hengadev/cluo_api/internal/infrastructure/postgres/ai_transcription_job"
	aiTranscriptionRepo "github.com/hengadev/cluo_api/internal/infrastructure/postgres/ai_transcription"
	whisperClient "github.com/hengadev/cluo_api/internal/infrastructure/whisper"
	"github.com/hengadev/cluo_api/internal/infrastructure/webhook"
	"github.com/hengadev/cluo_api/internal/infrastructure/worker"
	aiChatApp "github.com/hengadev/cluo_api/internal/application/ai_chat"
	aiTextTransformationApp "github.com/hengadev/cluo_api/internal/application/ai_text_transformation"
	aiSpeechToTextApp "github.com/hengadev/cluo_api/internal/application/ai_speech_to_text"
	aiTranscriptAnalysisApp "github.com/hengadev/cluo_api/internal/application/ai_transcript_analysis"
)

// initAIServices initializes AI services (Ollama, Whisper, Workers).
func (c *Container) initAIServices(ctx context.Context) error {
	c.analysisRepo = aiTranscriptAnalysisRepo.New(ctx, c.dbPool)

	// Initialize Whisper first: it sets transcriptionRepo, which the transcript
	// analysis service (constructed during Ollama init below) depends on.
	if err := c.initWhisper(ctx); err != nil {
		if c.config.AI.Whisper.Enabled {
			c.logger.WarnContext(ctx, "Whisper initialization failed", "error", err)
		}
	}

	// Initialize Ollama client
	if err := c.initOllama(ctx); err != nil {
		if c.config.AI.Ollama.Enabled {
			c.logger.WarnContext(ctx, "Ollama initialization failed", "error", err)
		}
	}

	return nil
}

func (c *Container) initOllama(ctx context.Context) error {
	cfg := c.config.AI.Ollama

	if !cfg.Enabled {
		c.logger.InfoContext(ctx, "Ollama disabled, skipping initialization")
		return nil
	}

	// Create Ollama client
	ollama, err := ollamaClient.New(cfg)
	if err != nil {
		return fmt.Errorf("create ollama client: %w", err)
	}

	// Create text transformation service
	c.textTransformationService = aiTextTransformationApp.New(
		ollama,
		cfg.MaxInputLen,
		c.logger,
	)

	// Create transcript analysis service
	c.transcriptAnalysisService = aiTranscriptAnalysisApp.New(
		ollama,
		c.transcriptionRepo,
		c.analysisRepo,
		c.crypto,
		c.logger,
	)

	// Create chat repository and service
	chatRepo := aiChatRepo.New(ctx, c.dbPool)
	c.chatService = aiChatApp.New(chatRepo, ollama, c.crypto, c.logger)

	c.logger.InfoContext(ctx, "Ollama client initialized",
		"baseURL", cfg.BaseURL,
		"model", cfg.Model,
	)

	return nil
}

func (c *Container) initWhisper(ctx context.Context) error {
	cfg := c.config.AI.Whisper

	if !cfg.Enabled {
		c.logger.InfoContext(ctx, "Whisper disabled, skipping initialization")
		return nil
	}

	// Create Whisper client
	whisper, err := whisperClient.New(cfg)
	if err != nil {
		return fmt.Errorf("create whisper client: %w", err)
	}

	// Create audio processor
	tempDir := os.TempDir()
	audioProcessor := whisperClient.NewAudioProcessor(tempDir)

	// Create transcription repositories
	workerID := "worker-" + uuid.New().String()[:8]
	c.transcriptionJobRepo = aiTranscriptionJobRepo.New(ctx, c.dbPool, workerID)
	c.transcriptionRepo = aiTranscriptionRepo.New(ctx, c.dbPool)

	// Create webhook client
	webhookClient := webhook.New(c.logger)

	// Create transcription worker
	c.transcriptionWorker = worker.NewTranscriptionWorker(
		whisper,
		c.transcriptionJobRepo,
		c.transcriptionRepo,
		c.mediaRepo,
		webhookClient,
		c.crypto,
		c.config.AI.Worker.Concurrency,
		c.config.AI.Worker.QueueSize,
		cfg.DeleteAudio,
		c.logger,
	)

	// Create speech-to-text service
	c.speechToTextService = aiSpeechToTextApp.New(
		c.transcriptionJobRepo,
		c.transcriptionRepo,
		c.transcriptionWorker,
		audioProcessor,
		c.crypto,
		tempDir,
		c.logger,
	)

	c.logger.InfoContext(ctx, "Whisper client initialized",
		"binary", cfg.Binary,
		"model", cfg.Model,
		"workerConcurrency", c.config.AI.Worker.Concurrency,
	)

	return nil
}
