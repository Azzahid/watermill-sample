package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/IBM/sarama"
)

type Credentials struct {
	Protocol   string // SASL_SSL, SASL_PLAINTEXT, SSL, PLAINTEXT
	Mechanism  string // PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
	Username   string
	Password   string
	CACertPath string
	CertPath   string
	KeyPath    string
}

func AddCredentials(saramaConfig *sarama.Config, securityConfig Credentials) error {
	// Configure base SASL settings if needed
	if isSASLProtocol(securityConfig.Protocol) {
		if err := configureSASLBase(saramaConfig, securityConfig); err != nil {
			return fmt.Errorf("failed to configure SASL: %w", err)
		}
	}

	// Configure SSL if needed
	if isSSLProtocol(securityConfig.Protocol) {
		if err := configureSSL(saramaConfig, securityConfig); err != nil {
			return fmt.Errorf("failed to configure SSL: %w", err)
		}
	}

	return nil
}

func isSASLProtocol(protocol string) bool {
	return protocol == "SASL_SSL" || protocol == "SASL_PLAINTEXT"
}

func isSSLProtocol(protocol string) bool {
	return protocol == "SASL_SSL" || protocol == "SSL"
}

func configureSASLBase(config *sarama.Config, securityConfig Credentials) error {
	config.Net.SASL.Enable = true
	config.Net.SASL.User = securityConfig.Username
	config.Net.SASL.Password = securityConfig.Password

	return configureSASLMechanism(config, securityConfig.Mechanism)
}

func configureSASLMechanism(config *sarama.Config, mechanism string) error {
	switch mechanism {
	case "PLAIN":
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	case "SCRAM-SHA-256":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA256}
		}
	case "SCRAM-SHA-512":
		config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
		}
	default:
		return fmt.Errorf("unsupported SASL mechanism: %s", mechanism)
	}
	return nil
}

func configureSSL(config *sarama.Config, securityConfig Credentials) error {
	config.Net.TLS.Enable = true
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if err := configureCertificates(tlsConfig, securityConfig); err != nil {
		return err
	}

	config.Net.TLS.Config = tlsConfig
	return nil
}

func configureCertificates(tlsConfig *tls.Config, securityConfig Credentials) error {
	// Configure CA certificate if provided
	if securityConfig.CACertPath != "" {
		if err := configureRootCA(tlsConfig, securityConfig.CACertPath); err != nil {
			return err
		}
	}

	// Configure client certificates if both cert and key are provided
	if securityConfig.CertPath != "" && securityConfig.KeyPath != "" {
		if err := configureClientCert(tlsConfig, securityConfig.CertPath, securityConfig.KeyPath); err != nil {
			return err
		}
	}

	return nil
}

func configureRootCA(tlsConfig *tls.Config, caCertPath string) error {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("failed to parse CA certificate")
	}

	tlsConfig.RootCAs = caCertPool
	return nil
}

func configureClientCert(tlsConfig *tls.Config, certPath, keyPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("failed to load client certificate: %w", err)
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	return nil
}
