package email

import (
	"bufio"
	"context"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hengadev/cluo_api/internal/app/config"
)

// mockSMTPServer is a minimal SMTP server for testing.
type mockSMTPServer struct {
	ln      net.Listener
	mu      sync.Mutex
	lastMsg string
}

func newMockSMTPServer(t *testing.T) *mockSMTPServer {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	s := &mockSMTPServer{ln: ln}
	go s.serve()
	t.Cleanup(func() { ln.Close() })
	return s
}

func (s *mockSMTPServer) addr() string { return s.ln.Addr().String() }

func (s *mockSMTPServer) message() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastMsg
}

func (s *mockSMTPServer) serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *mockSMTPServer) handle(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)

	write := func(msg string) {
		bw.WriteString(msg + "\r\n")
		bw.Flush()
	}
	readLine := func() string {
		line, _ := br.ReadString('\n')
		return strings.TrimSpace(line)
	}

	write("220 test ESMTP")
	for {
		line := readLine()
		if line == "" {
			return
		}
		cmd := strings.ToUpper(line)

		switch {
		case strings.HasPrefix(cmd, "EHLO"):
			write("250-test Hello")
			write("250 OK")
		case strings.HasPrefix(cmd, "HELO"):
			write("250 Hello")
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			write("250 OK")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			write("250 OK")
		case cmd == "DATA":
			write("354 Go ahead")
			var data strings.Builder
			for {
				l, err := br.ReadString('\n')
				if err != nil {
					return
				}
				trimmed := strings.TrimRight(l, "\r\n")
				if trimmed == "." {
					break
				}
				data.WriteString(l)
			}
			write("250 OK: queued")

			s.mu.Lock()
			s.lastMsg = data.String()
			s.mu.Unlock()
			// Don't return — let the client send QUIT.
		case cmd == "QUIT":
			write("221 Bye")
			return
		case cmd == "RSET":
			write("250 OK")
		case cmd == "NOOP":
			write("250 OK")
		default:
			write("502 Error")
		}
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestSMTPAdapter_Send_Success(t *testing.T) {
	srv := newMockSMTPServer(t)
	host, port, _ := strings.Cut(srv.addr(), ":")

	cfg := config.SMTPConfig{Host: host, Port: port, From: "sender@example.com"}
	adapter := NewSMTPAdapter(cfg, slog.Default())

	err := adapter.Send(context.Background(), "recipient@example.com", "Test Subject", "<h1>Hello</h1>")
	require.NoError(t, err)

	msg := srv.message()
	assert.Contains(t, msg, "To: recipient@example.com")
	assert.Contains(t, msg, "Subject: Test Subject")
	assert.Contains(t, msg, "<h1>Hello</h1>")
}

func TestSMTPAdapter_Send_VerifyFromAndContentType(t *testing.T) {
	srv := newMockSMTPServer(t)
	host, port, _ := strings.Cut(srv.addr(), ":")

	cfg := config.SMTPConfig{Host: host, Port: port, From: "noreply@cluo.com"}
	adapter := NewSMTPAdapter(cfg, slog.Default())

	err := adapter.Send(context.Background(), "user@test.com", "Welcome", "<p>Welcome!</p>")
	require.NoError(t, err)

	msg := srv.message()
	assert.Contains(t, msg, "From: noreply@cluo.com")
	assert.Contains(t, msg, "Content-Type: text/html")
}

func TestSMTPAdapter_Send_InvalidHost(t *testing.T) {
	cfg := config.SMTPConfig{Host: "invalid.host.that.does.not.exist", Port: "25", From: "sender@example.com"}
	adapter := NewSMTPAdapter(cfg, slog.Default())
	err := adapter.Send(context.Background(), "recipient@example.com", "Test", "body")
	assert.Error(t, err)
}

func TestSMTPAdapter_Send_ConnectionRefused(t *testing.T) {
	cfg := config.SMTPConfig{Host: "127.0.0.1", Port: "19999", From: "sender@example.com"}
	adapter := NewSMTPAdapter(cfg, slog.Default())
	err := adapter.Send(context.Background(), "to@example.com", "Test", "body")
	assert.Error(t, err)
}

func TestNoOpAdapter_Send(t *testing.T) {
	adapter := NewNoOpAdapter(slog.Default())
	err := adapter.Send(context.Background(), "to@example.com", "Subject", "<html></html>")
	assert.NoError(t, err)
}
