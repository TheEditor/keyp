package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/TheEditor/keyp/internal/server"
)

var (
	servePort    = 8080
	serveBind    = "127.0.0.1"
	serveTimeout = 15 * time.Minute
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP API server",
	Long:  "Start the HTTP API server for remote access to the vault via REST endpoints.",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().IntVar(&servePort, "port", 8080, "HTTP server port (default: 8080)")
	serveCmd.Flags().StringVar(&serveBind, "bind", "127.0.0.1", "Address to bind to (default: 127.0.0.1)")
	serveCmd.Flags().DurationVar(&serveTimeout, "timeout", 15*time.Minute, "Session timeout duration (default: 15m)")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	vaultPath := getVaultPath()
	address := fmt.Sprintf("%s:%d", serveBind, servePort)

	// Create server
	srv := server.NewServer(address, vaultPath)
	srv.SetSessionTimeout(serveTimeout)

	// Start server in goroutine
	errs := make(chan error, 1)
	go func() {
		errs <- srv.Start()
	}()

	fmt.Printf("Server listening on http://%s\n", address)
	fmt.Println("Press Ctrl+C to shutdown...")

	// Wait for interrupt or error
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errs:
		return err
	case sig := <-sigChan:
		fmt.Printf("\nReceived signal: %v\n", sig)

		// Graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown failed: %w", err)
		}

		fmt.Println("Server shutdown complete")
		return nil
	}
}
