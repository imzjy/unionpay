package gounionpay

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AppTrans struct {
	Config *UnionpayConfig
}

func NewAppTrans(cfg *UnionpayConfig) *AppTrans {
	return &AppTrans{Config: cfg}
}

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

	Sign(this.Config.SignKeyPath, this.Config.SignCertPath, param)

	response, err := doTrans(param, this.Config.AppTransUrl)
	if err != nil {
		return "", err
	}


	respValue := parseResponse(response)
	respCode, ok := respValue["respCode"]
	if !ok {
		return "", errors.New("respCode field not found")
	}
	if respCode != "00" {
		printMap(respValue)
		return "", errors.New("respCode:" + respCode)
	}

	err = Validate(this.Config.VerifyCertPath, respValue)
	if err != nil {
		return "", err
	}

	tn, ok := respValue["tn"]
	if ok {
		// fmt.Println("TN:", tn)
		return tn, nil
	} else {
		printMap(respValue)
		return "", errors.New("tn field not found")
	}
}

func parseResponse(resp []byte) map[string]string {

	retMap := make(map[string]string)
	content := strings.Split(string(resp), "&")

	for _, item := range content {

		//strings.Split(s, "=") will cause error when signature has padding(that is something like "==")
		idx := strings.IndexAny(item, "=")
		if idx < 0 {
			panic("response value parse error:" + item)
		}

		k := item[:idx]
		v := item[idx+1:]
		retMap[k] = v
	}

	return retMap
}

func doTrans(param map[string]string, appTransUrl string) ([]byte, error) {
	datagram := concatParam(param)

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

func concatParam(param map[string]string) bytes.Buffer {
	var sortedParam []string
	for k, v := range param {
		// fmt.Println(k, "=", UrlEncoded(v))
		sortedParam = append(sortedParam, k+"="+urlEncode(v))
	}

	return *bytes.NewBufferString(strings.Join(sortedParam, "&"))
}

func urlEncode(str string) string {
	// fmt.Println("in:", str)
	encodedUrl := url.QueryEscape(str)
	// fmt.Println("out:", encodedUrl)

	return encodedUrl
}

func printMap(m map[string]string) {
	for k, v := range m {
		fmt.Println(k, "=", v)
	}
}
