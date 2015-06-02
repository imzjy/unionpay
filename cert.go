// Parse certification
// With some handy functions, su as get public key, serial number of a certifcaton
package gounionpay

import (
	"math/big"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"errors"
)


func ParseCertificate(pemData []byte) (*x509.Certificate, error){
	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("cannot decode the pem file")
	}
	if got, want := block.Type, "CERTIFICATE"; got != want {
		return nil, fmt.Errorf("unknown key type %q, want %q", got, want)
	}

	// Decode the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("bad private key: %s", err)
	}

	return cert, nil
}

func CertSerialNumber(pemData []byte) (*big.Int, error) {
	cert, err := ParseCertificate(pemData)
	if err != nil {
		return big.NewInt(0), err
	}

	return cert.SerialNumber, nil
}

func CertPublickey(pemData []byte)(interface{}, error) {
	cert, err := ParseCertificate(pemData)
	if err != nil {
		return nil, err
	}

	return cert.PublicKey, nil
}