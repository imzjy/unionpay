package unionpay

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	// "fmt"
	"fmt"
	"io/ioutil"
)

func sha1DigestFromString(data string) []byte {
	return sha1Digest([]byte(data))
}

func sha1Digest(data []byte) []byte {

	h := sha1.New()
	h.Write(data)

	return h.Sum(nil)
}

func rsaSignBySha1(keypath string, in []byte) ([]byte, error) {

	// Read the private key
	pemData, err := ioutil.ReadFile(keypath)
	if err != nil {
		return []byte(""), fmt.Errorf("read key file: %s", err)
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return []byte(""), fmt.Errorf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return []byte(""), fmt.Errorf("unknown key type %q, want %q", got, want)
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return []byte(""), fmt.Errorf("bad private key: %s", err)
	}

	// fmt.Println("priv:", priv)
	// fmt.Println("in:", in, "len(in):", len(in))

	encData, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA1, sha1Digest(in))
	if err != nil {
		return []byte(""), err
	}

	return encData, nil
}

func rsaVerifyBySha1(cerpath string, signature, in []byte) error {
	// Read the verify sign certification key
	pemData, err := ioutil.ReadFile(cerpath)
	if err != nil {
		return err
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		return fmt.Errorf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "CERTIFICATE"; got != want {
		return fmt.Errorf("unknown key type %q, want %q", got, want)
	}

	// Decode the certification
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("bad private key: %s", err)
	}
	// fmt.Println(cert)

	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, sha1Digest(in), signature)
	if err != nil {
		return fmt.Errorf("VerifyPKCS1v15 fail: %s", err)
	}

	return nil
}

func base64String(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func base64Bytes(in string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(in)
}
