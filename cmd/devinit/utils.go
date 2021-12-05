package devinit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/nrc-no/core/pkg/utils/files"
	"math/big"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"
)

func getOrCreatePrivateKey(fileName string) (*rsa.PrivateKey, error) {
	exists, err := files.FileExists(fileName)
	if err != nil {
		return nil, err
	}
	if exists {
		keyBytes, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		privPem, _ := pem.Decode(keyBytes)
		var parsedKey interface{}
		if parsedKey, err = x509.ParsePKCS1PrivateKey(privPem.Bytes); err != nil {
			if parsedKey, err = x509.ParsePKCS8PrivateKey(privPem.Bytes); err != nil {
				return nil, err
			}
		}
		var privateKey *rsa.PrivateKey
		var ok bool
		privateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("could not decode private key")
		}
		return privateKey, nil
	}
	newKey, err := genPrivateKey()
	if err != nil {
		return nil, err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(newKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(fileName, pubBytes, os.ModePerm); err != nil {
		return nil, err
	}
	return newKey, nil
}

func genCert(template, parent *x509.Certificate, publicKey *rsa.PublicKey, caKey *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	certPEM := pem.EncodeToMemory(&b)

	return cert, certPEM, nil
}

func genPrivateKey() (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func getOrCreateCaRoot(path string, privateKey *rsa.PrivateKey) (*x509.Certificate, error) {
	exists, err := files.FileExists(path)
	if err != nil {
		return nil, err
	}
	if exists {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(fileBytes)
		return x509.ParseCertificate(block.Bytes)
	}
	rootCert, rootPem, err := genCARoot(privateKey)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, rootPem, os.ModePerm); err != nil {
		return nil, err
	}
	return rootCert, nil
}

func genCARoot(caKey *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(99999999999999999))
	if err != nil {
		return nil, nil, err
	}
	var rootTemplate = x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Company Co."},
			Country:      []string{"DE"},
			Province:     []string{"Berlin"},
			Locality:     []string{"Berlin"},
			CommonName:   "Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(2, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		IPAddresses: []net.IP{net.ParseIP(
			"127.0.0.1",
		)},
	}
	return genCert(&rootTemplate, &rootTemplate, &caKey.PublicKey, caKey)
}

func getOrCreateServerCert(
	path string,
	privateKey *rsa.PrivateKey,
	parentCert *x509.Certificate,
	parentKey *rsa.PrivateKey,
	dsnNames []string,
	ipAddresses []net.IP,
) (*x509.Certificate, error) {
	exists, err := files.FileExists(path)
	if err != nil {
		return nil, err
	}
	if exists {
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(fileBytes)
		return x509.ParseCertificate(block.Bytes)
	}
	rootCert, rootPem, err := genServerCert(privateKey, parentCert, parentKey, dsnNames, ipAddresses)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, rootPem, os.ModePerm); err != nil {
		return nil, err
	}
	return rootCert, nil
}

func genServerCert(
	privateKey *rsa.PrivateKey,
	parentCert *x509.Certificate,
	caKey *rsa.PrivateKey,
	dsnNames []string,
	ipAddresses []net.IP,
) (*x509.Certificate, []byte, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(99999999999999999))
	if err != nil {
		return nil, nil, err
	}
	var ServerTemplate = x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{"Company Co."},
			Country:      []string{"DE"},
			Province:     []string{"Berlin"},
			Locality:     []string{"Berlin"},
			CommonName:   "Server Certificate",
		},
		SerialNumber: serial,
		NotBefore:    time.Now().Add(-10 * time.Second),
		NotAfter:     time.Now().AddDate(2, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		IsCA:         false,
		DNSNames:     dsnNames,
		IPAddresses:  ipAddresses,
	}
	return genCert(&ServerTemplate, parentCert, &privateKey.PublicKey, caKey)
}

func getOrCreateRandomSecretStr(length int, filePaths ...string) (string, error) {
	secretBytes, err := getOrCreateRandomSecret(length, filePaths...)
	return string(secretBytes), err
}

func getOrCreateRandomSecret(length int, filePaths ...string) ([]byte, error) {
	filePath := path.Join(filePaths...)
	exists, err := files.FileExists(filePath)
	if err != nil {
		return nil, err
	}

	if exists {
		return os.ReadFile(filePath)
	}

	if _, err := files.CreateDirectoryIfNotExists(filepath.Dir(filePath)); err != nil {
		return nil, err
	}
	value := []byte(randomStringBase64(length))
	if err := os.WriteFile(filePath, value, os.ModePerm); err != nil {
		return nil, err
	}
	return value, nil
}

func createRandomSecretIfNotExists(length int, filePaths ...string) error {
	filePath := path.Join(filePaths...)
	if _, err := files.CreateDirectoryIfNotExists(filepath.Dir(filePath)); err != nil {
		return err
	}
	if err := os.WriteFile(filePath, []byte(randomStringBase64(length)), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func randomStringBase64(length int) string {
	return base64.StdEncoding.EncodeToString([]byte(randomString(length)))
}

func randomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
