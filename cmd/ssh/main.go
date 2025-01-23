package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	config "github.com/NesterovYehor/txtnest-cli/configs"
	"github.com/NesterovYehor/txtnest-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/google/uuid"
)

func main() {
	// Load configuration from environment variables or use defaults
	cfg, err := config.LoadSSHConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	// Initialize the AppModel with the backend URL
	appModel := ui.NewAppModel(cfg.BackendURL)

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	// SSH server setup
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("0.0.0.0", cfg.SSHPort)),
		wish.WithHostKeyPEM(cfg.PrivateKey), // Use your SSH private key here
		wish.WithMiddleware(bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
			return appModel, []tea.ProgramOption{tea.WithAltScreen()}
		}), activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			anonymousMiddleware,
			logging.Middleware()), wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
	)
	if err != nil {
		log.Fatalf("Could not start SSH server: %v", err)
	}

	// Start the SSH server in a goroutine
	go func() {
		log.Println("Starting SSH server on port 2222...")
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("Could not start SSH server: %v", err)
		}
	}()

	// Wait for a signal to cancel the context
	<-ctx.Done()
	s.Shutdown(ctx)
	log.Println("SSH server shutting down...")
}

func anonymousMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		// Generate a unique session ID
		sessionID := uuid.New().String()

		// Log the session ID (optional)
		log.Printf("New anonymous session: %s", sessionID)

		// Welcome the user
		fmt.Fprintln(s, "Welcome to the anonymous SSH app!")
		fmt.Fprintf(s, "Your session ID: %s\n", sessionID)
		fmt.Fprintln(s, "Use this app to create and manage pastes anonymously.")

		// Pass the session to the next handler
		if next != nil {
			next(s)
		}
	}
}
