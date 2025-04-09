package main

// Description:
//
// This program implements a fake SSH server that logs connection attempts and
// authentication failures. It uses the SSH protocol to accept incoming connections
// and checks the provided credentials against a predefined set of users and passwords.
//
// The server listens on a specified port and logs the type of authentication attempt
// (legitimate, local, admin, dictionary, or brute force) based on the provided credentials.
//
// It also includes a configuration file for managing users, passwords, and other settings.
// The server can be used for security testing and monitoring purposes, providing insights
// into potential vulnerabilities and unauthorized access attempts.
//
// The program is designed to be run as a standalone application and can be configured
// through a YAML configuration file. It includes error handling for various scenarios,
// such as invalid port numbers, failed key parsing, and connection acceptance errors.
//
// The server can be extended to include additional features, such as signal handling
// for graceful shutdown, logging to syslog, and more advanced authentication mechanisms.
//
// It is important to note that this program is intended for educational and testing
// purposes only and should not be used for malicious activities or unauthorized access
// to systems. Always ensure that you have the necessary permissions and legal authority
// before conducting any security testing or monitoring activities on a network or system.

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

// DNS names suggesting worth attacking:
//
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

var (
	configFile        = "config.yaml"
	Version    string = "0.2"
)

func main() {
	log.SetFlags(log.LstdFlags) // Add timestamps to logs

	cfg := readconfig(configFile)

	privateBytes, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	portNum, err := strconv.Atoi(cfg.Port)
	if err != nil || portNum < 1 || portNum > 65535 {
		log.Fatalf("Invalid port: %s. Port must be an integer between 1 and 65535.", cfg.Port)
	}

	// SSH server configuration
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {

			ip := ipcut(c.RemoteAddr().String())
			_, legit := cfg.Hosts[ip] // We know the IP is in the list of hosts

			switch {

			case legit:
				log.Printf("attempt type Legit '%s/'%s'/%s'",
					cfg.Hosts[ip], c.User(), string(pass))

			case ipinrange(ip) && !legit:
				log.Printf("attempt type Local '%s/'%s'/%s'",
					c.User(), string(pass), ip)
				inform("local", c.User(), string(pass), ip)

			case contains(cfg.NoPassDump, c.User()):
				log.Printf("attempt type Admin '%s/'*'/%s'", c.User(), ip)
				inform("admin", c.User(), "", ip)

			case cfg.Users[c.User()] == string(pass) && len(pass) > 0:
				log.Printf("attempt type Dict '%s/'%s'/%s'",
					c.User(), string(pass), ip)
				inform("dict", c.User(), string(pass), ip)

			default:
				log.Printf("attempt type Bruteforce '%s/'%s'/%s'",
					c.User(), string(pass), ip)
			}

			return nil, fmt.Errorf("access denied")
		},
	}
	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
	}
	defer listener.Close()

	fmt.Printf("FakeSSH %s server listening on port %s...\n", Version, cfg.Port)

	initSyslog("fakessh")

	// // Signal handling for graceful shutdown
	// ctx, cancel := context.WithCancel(context.Background())
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// go func() {
	// 	<-sigs
	// 	fmt.Println("Shutting down server...")
	// 	cancel()
	// }()

	for {
		select {
		// case <-ctx.Done():
		// 	log.Println("Server stopped gracefully.")
		// 	return
		default:
			nConn, err := listener.Accept()
			if err != nil {
				// Check if the context is canceled to exit the loop
				// if ctx.Err() != nil {
				// 	log.Println("Server stopped gracefully.")
				// 	return
				// }
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
