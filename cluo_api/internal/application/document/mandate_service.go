package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Mandate operations

func (s *Service) CreateMandate(ctx context.Context, mandate *domain.Mandate) (*domain.Mandate, error) {
	// TODO: Validate user permissions to create mandates
	// TODO: Verify that case and client exist and are accessible

	// Validate mandate
	if err := mandate.Validate(); err != nil {
		return nil, fmt.Errorf("mandate validation failed: %w", err)
	}

	// Generate mandate number if not provided
	if mandate.MandateNumber == "" {
		// TODO: Implement mandate number generation
		mandate.MandateNumber = fmt.Sprintf("MND-%d", uuid.New().ID())
	}

	// Save to repository
	if err := s.repo.CreateMandate(ctx, mandate); err != nil {
		return nil, fmt.Errorf("failed to create mandate: %w", err)
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, mandate, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

func (s *Service) SignMandate(ctx context.Context, mandateID string, req *domain.SignDocumentRequest) (*domain.Mandate, error) {
	// TODO: Validate user permissions to sign mandates
	// TODO: Verify that signer has authority to sign for the specified role

	// Get mandate
	mandate, err := s.repo.GetMandateByID(ctx, mandateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	// Check if mandate can be signed by this role
	if !mandate.CanBeSigned(req.SignerRole) {
		return nil, fmt.Errorf("mandate cannot be signed by %s as %s", req.SignerName, req.SignerRole)
	}

	// Create signature
	signature := domain.NewSignature(
		req.SignerName,
		req.SignerRole,
		req.Method,
		req.SignatureFileURL,
		nil, // TODO: Get signer ID from context
	)

	// Add additional fields if provided
	if req.IPAddress != nil {
		signature.IPAddress = req.IPAddress
	}
	if req.UserAgent != nil {
		signature.UserAgent = req.UserAgent
	}

	// Add signature to mandate
	var signerUUID uuid.UUID
	var err error

	switch req.SignerRole {
	case "client":
		if err := mandate.AddClientSignature(signature); err != nil {
			return nil, fmt.Errorf("failed to add client signature: %w", err)
		}
		// TODO: Get actual signer ID from context
		signerUUID = uuid.New()

	case "investigator":
		if err := mandate.AddInvestigatorSignature(signature); err != nil {
			return nil, fmt.Errorf("failed to add investigator signature: %w", err)
		}
		// TODO: Get actual signer ID from context
		signerUUID = uuid.New()

	default:
		return nil, fmt.Errorf("unsupported signer role: %s", req.SignerRole)
	}

	// Validate updated mandate
	if err := mandate.Validate(); err != nil {
		return nil, fmt.Errorf("updated mandate validation failed: %w", err)
	}

	// Update mandate
	if err := s.repo.UpdateMandate(ctx, mandate); err != nil {
		return nil, fmt.Errorf("failed to update mandate: %w", err)
	}

	// Create version record
	reason := fmt.Sprintf("Signed by %s as %s", req.SignerName, req.SignerRole)
	if err := s.createDocumentVersion(ctx, mandate, &signerUUID, stringPtr(reason)); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

func (s *Service) ActivateMandate(ctx context.Context, mandateID string) (*domain.Mandate, error) {
	// TODO: Validate user permissions to activate mandates

	// Get mandate
	mandate, err := s.repo.GetMandateByID(ctx, mandateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	// Activate mandate
	if err := mandate.Activate(); err != nil {
		return nil, fmt.Errorf("failed to activate mandate: %w", err)
	}

	// Update mandate
	if err := s.repo.UpdateMandate(ctx, mandate); err != nil {
		return nil, fmt.Errorf("failed to update mandate: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	if err := s.createDocumentVersion(ctx, mandate, nil, stringPtr("Mandate activated")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

func (s *Service) CreateContractFromMandate(ctx context.Context, mandateID string, contract *domain.Contract) (*domain.Contract, error) {
	// TODO: Validate user permissions to create contracts
	// TODO: Verify that case and client exist and are accessible

	// Get mandate
	mandate, err := s.repo.GetMandateByID(ctx, mandateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	// Verify mandate is in appropriate state
	if mandate.Status != domain.DocumentStatusSigned && mandate.Status != domain.DocumentStatusActive {
		return nil, fmt.Errorf("mandate must be signed or active to create contract")
	}

	// Link contract to mandate
	if err := contract.LinkToMandate(mandate.ID); err != nil {
		return nil, fmt.Errorf("failed to link contract to mandate: %w", err)
	}

	// Ensure contract belongs to same case and client
	if contract.CaseID != mandate.CaseID {
		return nil, fmt.Errorf("contract must belong to same case as mandate")
	}
	if contract.ClientID != mandate.ClientID {
		return nil, fmt.Errorf("contract must belong to same client as mandate")
	}

	// Generate contract number if not provided
	if contract.ContractNumber == "" {
		// TODO: Implement contract number generation
		contract.ContractNumber = fmt.Sprintf("CNT-%d", uuid.New().ID())
	}

	// Validate contract
	if err := contract.Validate(); err != nil {
		return nil, fmt.Errorf("contract validation failed: %w", err)
	}

	// Save contract
	if err := s.repo.CreateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	reason := fmt.Sprintf("Created from mandate %s", mandate.MandateNumber)
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr(reason)); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}