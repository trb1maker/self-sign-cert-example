package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	keySize          = 2048
	keyName          = "pravo.tech.key"
	certName         = "pravo.tech.pem"
	organizationName = "Pravo.Tech"
)

type Generator struct {
	privateKey  *rsa.PrivateKey
	certificate *x509.Certificate
}

// CreateCertificate генерирует приватный ключ и сертификат публичного ключа
// для указанных domens и сохраняет их в outDirName.
func (g *Generator) CreateCertificate(outDirName string, domens ...string) error {
	if err := g.createPrivrateKey(); err != nil {
		return err
	}

	if err := g.createCertificate(domens...); err != nil {
		return err
	}

	if err := g.store(outDirName); err != nil {
		return err
	}

	return nil
}

// createPrivrateKey генерирует приватный RSA-ключ.
func (g *Generator) createPrivrateKey() (err error) { //nolint:nonamedreturns
	g.privateKey, err = rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return fmt.Errorf("can't generate private key: %w", err)
	}

	return nil
}

// createCertificate создает сертификат открытого ключа.
func (g *Generator) createCertificate(domens ...string) error {
	serianNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return fmt.Errorf("can't create certificate: %w", err)
	}

	//nolint:lll
	g.certificate = &x509.Certificate{
		SerialNumber:          serianNumber,
		Subject:               pkix.Name{Organization: []string{organizationName}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),                                   // Действие сертификата 1 год
		KeyUsage:              x509.KeyUsageDataEncipherment | x509.KeyUsageDigitalSignature, // Основное назначение: шифрование данных и подписание сертификатов
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},                // Дополнительное назначение: авторизация сервера
		BasicConstraintsValid: true,
		DNSNames:              append(domens, "localhost"), // Добавляю localhost в список разрешенных доменов - основное назначение
	}

	return nil
}

// store сохраняет созданные ключ и сертификат в dir.
func (g *Generator) store(dir string) error {
	certDER, err := x509.CreateCertificate(
		rand.Reader,
		g.certificate,
		g.certificate,
		&g.privateKey.PublicKey,
		g.privateKey,
	)
	if err != nil {
		return fmt.Errorf("can't create DER certificate: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(dir, certName),
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}),
		0600, //nolint:gofumpt
	); err != nil {
		return fmt.Errorf("can't write certificate: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(dir, keyName),
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(g.privateKey)}),
		0600, //nolint:gofumpt
	); err != nil {
		return fmt.Errorf("can't write private key: %w", err)
	}

	return nil
}
