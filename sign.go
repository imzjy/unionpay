package gounionpay

import (
	"fmt"
	"sort"
	"strings"
)

func Sign(keypath string, param map[string]string) error {

	//证书序列号
	param["certId"] = signCertId()

	//将Dictionary信息转换成key1=value1&key2=value2的形式
	// sortedPairStr := "accNo=6225682141000002950&accessType=0&backUrl=https://101.231.204.80:5000/gateway/api/backTransReq.do&bizType=000201&certId=124876885185794726986301355951670452718&channelType=07&currencyCode=156&encoding=UTF-8&merId=898340183980105&orderId=2014110600007615&signMethod=01&txnAmt=000000010000&txnSubType=01&txnTime=20150109135921&txnType=01&version=5.0.0"
	sortedPairStr := sortAndConcat(param)
	// fmt.Println(sortedPairStr)

	signedDigest := Sha1DigestFromString(sortedPairStr)
	// fmt.Println(signedDigest)

	hexSignedDigest := fmt.Sprintf("%x", signedDigest)
	// fmt.Println("sha1:", hexSignedDigest)

	byteSign := sha1RsaSign(keypath, []byte(hexSignedDigest))
	// fmt.Println(byteSign)

	// fmt.Println("sign:", base64String(byteSign))
	//设置签名
	param["signature"] = base64String(byteSign)


	return nil

}

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

func sortAndConcat(param map[string]string) string {
	var keys []string
	for k := range param {
		keys = append(keys, k)
	}

	var sortedParam []string
	sort.Strings(keys)
	for _, k := range keys {
		// fmt.Println(k, "=", param[k])
		sortedParam = append(sortedParam, k+"="+param[k])
	}

	return strings.Join(sortedParam, "&")
}
