package unionpay

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// AppTrans is client for submit order and get back the TN(transaction number) from
// china unionpay
type AppTrans struct {
	Config *UnionpayConfig
}

// NewApppTrans initial the AppTrans with specific configuration
func NewAppTrans(cfg *UnionpayConfig) *AppTrans {
	return &AppTrans{Config: cfg}
}

// Submit the order to china unionpay and return the TN(transaction number) if success,
// TN is used by mobile app.
// If fail, error is not nil, check error for more information
func (this *AppTrans) Submit(orderId string, amount float64, desc string) (string, error) {

	if this.Config.SignKeyPath == "" ||
		this.Config.SignCertPath == "" ||
		this.Config.VerifyCertPath == "" ||
		this.Config.CallbackUrl == "" ||
		this.Config.MerId == "" ||
		this.Config.AppTransUrl == "" {
		return "", errors.New("Please set key and cert")
	}

	txnTime := time.Now().Format("20060102030405")
	txnAmt := fmt.Sprintf("%.0f", amount) //单位为分，不可以有小数

	param := make(map[string]string)
	param["version"] = "5.0.0"
	param["encoding"] = "UTF-8"                                                //编码方式
	param["txnType"] = "01"                                                    //交易类型
	param["txnSubType"] = "01"                                                 //交易子类
	param["bizType"] = "000201"                                                //业务类型
	param["frontUrl"] = "http://not.exist.com/demo/utf8/FrontRcvResponse.aspx" //前台通知地址 ，控件接入方式无作用
	param["backUrl"] = this.Config.CallbackUrl                                 //后台通知地址
	param["signMethod"] = "01"                                                 //签名方法
	param["channelType"] = "08"                                                //渠道类型，07-PC，08-手机
	param["accessType"] = "0"                                                  //接入类型
	param["merId"] = this.Config.MerId                                         //商户号，请改成自己的商户号
	param["orderId"] = orderId                                                 //商户订单号
	param["txnTime"] = txnTime                                                 //订单发送时间
	param["txnAmt"] = txnAmt                                                   //交易金额，单位分
	param["currencyCode"] = "156"                                              //交易币种
	param["orderDesc"] = desc                                                  //订单描述，可不上送，上送时控件中会显示该信息
	param["reqReserved"] = "透传字段"                                              //请求方保留域，透传字段，查询、通知、对账文件中均会原样出现

	this.Sign(param)

	respMsg, err := doTrans(param, this.Config.AppTransUrl)
	if err != nil {
		return "", err
	}

	respValue, err := ParseResponseMsg(respMsg)
	if err != nil {
		return "", err
	}

	respCode, ok := respValue["respCode"]
	if !ok {
		return "", errors.New("respCode field not found")
	}
	if respCode != "00" {
		PrintMap(respValue)
		return "", errors.New("respCode:" + respCode)
	}

	err = this.Validate(respValue)
	if err != nil {
		return "", err
	}

	tn, ok := respValue["tn"]
	if ok {
		// fmt.Println("TN:", tn)
		return tn, nil
	} else {
		PrintMap(respValue)
		return "", errors.New("tn field not found")
	}
}


// Sign the data to comform with specs,
// more info refer to https://open.unionpay.com/ajweb/help/faq/detail?id=38
func (this *AppTrans) Sign(param map[string]string) error {

	//证书序列号
	certSn, err := certSerialNumberFromFile(this.Config.SignCertPath)
	if err != nil {
		return err
	}

	param["certId"] = certSn.String()

	sortedPairStr := SortAndConcat(param)
	signedDigest := sha1DigestFromString(sortedPairStr)
	hexSignedDigest := fmt.Sprintf("%x", signedDigest)

	byteSign, err := rsaSignBySha1(this.Config.SignKeyPath, []byte(hexSignedDigest))
	if err != nil {
		return err
	}

	//设置签名
	param["signature"] = base64String(byteSign)

	return nil
}


// Validate the response message with verfy certificate,
// more info refer to https://open.unionpay.com/ajweb/help/faq/detail?id=38
func (this *AppTrans) Validate(param map[string]string) error {
	//获取签名
	signature := param["signature"]
	// fmt.Println(signature)
	signByte, err := base64Bytes(signature)
	if err != nil {
		return err
	}

	delete(param, "signature")

	stringData := SortAndConcat(param)
	signedDigest := sha1DigestFromString(stringData)
	hexSignedDigest := fmt.Sprintf("%x", signedDigest)

	//TODO: check serial number of certifcate
	return rsaVerifyBySha1(this.Config.VerifyCertPath, signByte, []byte(hexSignedDigest))
}



func doTrans(param map[string]string, appTransUrl string) ([]byte, error) {
	datagram := ConcatWithUrlEncode(param)

	req, err := http.NewRequest("POST", appTransUrl, &datagram)
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}
