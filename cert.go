package gounionpay

import "math/big"

func signCertId() string {
	certNo := new(big.Int)
	certNo.SetString("5DF269CB0583CA185F34728C852C61EE", 16)

	return certNo.String()
}
