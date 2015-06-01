package sign

import (
	"cert"
)

func Sign(param map[string]string) error {
	//设置签名证书序列号 ？

	param["certId"] = CertUtil.GetSignCertId();

	//将Dictionary信息转换成key1=value1&key2=value2的形式
	string stringData = CoverDictionaryToString(param);

	string stringSign = null;

	byte[] signDigest = SecurityUtil.Sha1X16(stringData, encoder);

	string stringSignDigest = BitConverter.ToString(signDigest).Replace("-", "").ToLower();
	
	byte[] byteSign = SecurityUtil.SignBySoft(CertUtil.GetSignProviderFromPfx(), encoder.GetBytes(stringSignDigest));

	stringSign = Convert.ToBase64String(byteSign);

	//设置签名域值
	param["signature"] = stringSign;

	return true;
}
