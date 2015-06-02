package gounionpay

import (
	"fmt"
)

//Sign the submit data
//More info refer to https://open.unionpay.com/ajweb/help/faq/detail?id=38
func Sign(keypath string, certpath string, param map[string]string) error {

	//证书序列号
	certSn, err := CertSerialNumberFromFile(certpath)
	if err != nil {
		return err
	}

	param["certId"] = certSn.String()

	sortedPairStr := sortAndConcat(param)
	signedDigest := Sha1DigestFromString(sortedPairStr)
	hexSignedDigest := fmt.Sprintf("%x", signedDigest)

	byteSign, err := rsaSignBySha1(keypath, []byte(hexSignedDigest))
	if err != nil {
		return err
	}

	//设置签名
	param["signature"] = base64String(byteSign)

	return nil
}

//Validate the response message with verfy certificate
//More info refer to https://open.unionpay.com/ajweb/help/faq/detail?id=38
func Validate(certpath string, param map[string]string) error{
	//获取签名
	signature := param["signature"]
	// fmt.Println(signature)
	signByte := base64Bytes(signature+"==")

	delete(param, "signature")

	stringData := sortAndConcat(param)
	signedDigest := Sha1DigestFromString(stringData)
	hexSignedDigest := fmt.Sprintf("%x", signedDigest)


	//TODO: check serial number of certifcate
	return rsaVerifyBySha1(certpath, signByte, []byte(hexSignedDigest))
}
