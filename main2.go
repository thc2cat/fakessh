package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/crypto/ssh"
)

// Names Suggesting Sensitive Data:
// FinancialServer
// CustomerDatabase
// BankingData
// ConfidentialDocs
// CompanySecrets

// Names Suggesting Important Resources:
// MainServer
// DomainController
// PrimaryFirewall
// CriticalRouter
// MassStorage
// Names Suggesting Vulnerabilities:

// TestServer
// WebDevelopment
// LegacyMachine
// FirewallDisabled
// NoUpdates

// Names Suggesting Security Measures:
// SecureServer
// EncryptedData
// TwoFactorAuth
// SecureTunnel
// SecureConnection

// Names Suggesting Security Testing:
// PenTestServer
// VulnerabilityScanner
// SecurityAudit
// SecurityTesting
// SecurityAssessment

// Names Suggesting Security Incidents:
// DataBreach
// SecurityIncident
// UnauthorizedAccess
// CompromisedServer
// SecurityBreach

// Names Suggesting Security Monitoring:
// SecurityLogs
// IntrusionDetection
// SecurityMonitoring
// SecurityAlerts
// SecurityDashboard

func main() {
	log.SetFlags(log.LstdFlags) // Add timestamps to logs

	// Load private key path from environment variable or use default
	keyPath := os.Getenv("SSH_KEY_PATH")
	if keyPath == "" {
		keyPath = "id_rsa"
		log.Printf("Warning: No key path specified, using default '%s'", keyPath)
	}

	privateBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// SSH server configuration
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			log.Printf("fakessh attempt: user=%s, pass=%s, ip=%s", c.User(), string(pass), c.RemoteAddr())
			return nil, fmt.Errorf("access denied")
		},
	}
	config.AddHostKey(private)

	// Load port from environment variable or use default
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2022"
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port: %s. Port must be an integer between 1 and 65535.", port)
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %s...\n", port)

	// Signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("Shutting down server...")
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Server stopped gracefully.")
			return
		default:
			nConn, err := listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Server stopped gracefully.")
					return
				}
				log.Printf("Error accepting connection: %v", err)
				continue
			}

			go func(conn net.Conn) {
				defer conn.Close()

				_, _, _, err := ssh.NewServerConn(conn, config)
				if err != nil {
					// We don't need to log this error
					// log.Printf("Handshake failed with %s: %v", conn.RemoteAddr(), err)
					return
				}
			}(nConn)
		}
	}
}
