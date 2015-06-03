# gounionpay

中国银联手机控件支付后端实现

Implemenation of China unionpay backend for mobile app transaction with [golang](http://golang.org)

# usage

```go
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

appTrans := newAppTrans()
err = appTrans.Validate(respVal)
if err != nil {
	log(err)
}
```

# documentation

Please refer to [gowalker](https://gowalker.org/github.com/imzjy/gounionpay)
