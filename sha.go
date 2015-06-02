package gounionpay

import (
	"crypto/rsa"
	"crypto"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	// "fmt"
	"io/ioutil"
	"log"
)

func Sha1DigestFromString(data string) []byte{
	return Sha1Digest([]byte(data))
}

func Sha1Digest(data []byte) []byte {

	h := sha1.New()
	h.Write(data)

	return h.Sum(nil)
}

func sha1RsaSign(keypath string, in []byte) []byte {

	// Read the private key
	pemData, err := ioutil.ReadFile(keypath)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}

	// fmt.Println("priv:", priv)
	// fmt.Println("in:", in, "len(in):", len(in))

	encData, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA1, Sha1Digest(in)) 
	if err != nil {
		panic(err)
	}

	// rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), signedData) 

	// Write data to output file
	return encData
}

func sha1RsaVerify(cerpath string, signature, in []byte) error {
	// Read the verify sign certification key
	pemData, err := ioutil.ReadFile(cerpath)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "CERTIFICATE"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	// Decode the certification
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}
	// fmt.Println(cert)

	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, Sha1Digest(in), signature)
	if err != nil {
		log.Fatalf("VerifyPKCS1v15 fail: %s", err)	
	}

	return nil
}

func base64String(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func base64Bytes(in string) []byte{
	data, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	return data
}
