package gounionpay

type UnionpayConfig struct {
	SignKeyPath    string //加密密钥路径(convert from *.pfx)
	SignCertPath   string //加密证书路径(convert from *.pfx)
	VerifyCertPath string //验证证书路径

	CallbackUrl string //回调地址
	MerId       string //商户号
	AppTransUrl string //App方式交易提交地址
}
