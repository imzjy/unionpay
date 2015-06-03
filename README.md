# gounionpay

中国银联手机控件支付后端实现

Implemenation of China unionpay backend for mobile app transaction with [golang](http://golang.org)

# usage

```go
//初始化
cfg := &gounionpay.UnionpayConfig{
	SignKeyPath:    "sign key file path",
	SignCertPath:   "sign cert file path",
	VerifyCertPath: "verify cert file path",

	CallbackUrl: "服务端回调URL",
	MerId:       "商户号",
	AppTransUrl: "移动应用交易URL",
}
appTrans = gounionpay.NewAppTrans(cfg)

//获取tn，手机端得到tn后就可以使用这个tn发起支付调用
tn, err := appTrans.Submit(orderId, orderAmount, orderDescription)
if err != nil {
	log(err)
	return
}

//回调接口校验
respVal, err := gounionpay.ParseResponseMsg(respBody)
if err != nil {
	log(err)
	return
}

for rk, rv := range respVal {
	decVal, err := gounionpay.UrlDecode(rv)
	if err == nil {
		respVal[rk] = decVal
	}
}

err = appTrans.Validate(respVal)
if err != nil {
	log(err)
}
```

# documentation

Please refer to [gowalker](https://gowalker.org/github.com/imzjy/gounionpay)
