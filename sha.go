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

func sha1Digest(data string) []byte {
	
	h := sha1.New()
	h.Write([]byte(data))

	return h.Sum(nil)
}

func sha1RsaSign(in []byte) []byte {

	// Read the private key
	pemData, err := ioutil.ReadFile("/Users/zjy/Downloads/key.pem")
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

	encData, err := rsa.SignPKCS1v15(nil, priv, crypto.SHA1, sha1Digest(string(in))) 
	if err != nil {
		panic(err)
	}

	// rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), signedData) 

	// Write data to output file
	return encData
}

func base64String(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}
