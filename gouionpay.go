package gounionpay

import (
	"fmt"
)

//Sign the submit data
//More info refer to https://open.unionpay.com/ajweb/help/faq/detail?id=38
func Sign(keypath string, param map[string]string) error {

	//证书序列号
	param["certId"] = signCertId()

	sortedPairStr := sortAndConcat(param)
	signedDigest := Sha1DigestFromString(sortedPairStr)
	hexSignedDigest := fmt.Sprintf("%x", signedDigest)
	byteSign := sha1RsaSign(keypath, []byte(hexSignedDigest))

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
	return sha1RsaVerify(certpath, signByte, []byte(hexSignedDigest))
}
