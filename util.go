package unionpay

import (
	"sort"
	"strings"
	"bytes"
	"net/url"
	"fmt"
)

// SortAndConcat sort the map by key in ASCII order,
// and concat it in form of "k1=v1&k2=2"
func SortAndConcat(param map[string]string) string {
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

// ConcatWithUrlEncode concat the map to form of "k1=v1&k2=v2" and ensure "v1,v2"
// is Url encoded
func ConcatWithUrlEncode(param map[string]string) bytes.Buffer {
	var sortedParam []string
	for k, v := range param {
		// fmt.Println(k, "=", UrlEncode(v))
		sortedParam = append(sortedParam, k+"="+UrlEncode(v))
	}

	return *bytes.NewBufferString(strings.Join(sortedParam, "&"))
}

// ParseResponseMsg parse the response message in form of "k1=v1&k2=v2" to
// a map
func ParseResponseMsg(resp []byte) (map[string]string, error) {

	retMap := make(map[string]string)
	content := strings.Split(string(resp), "&")

	for _, item := range content {

		//strings.Split(s, "=") will cause error when signature has padding(that is something like "==")
		idx := strings.IndexAny(item, "=")
		if idx < 0 {
			return retMap, fmt.Errorf("parse error for value of %s", item)
		}

		k := item[:idx]
		v := item[idx+1:]
		retMap[k] = v
	}

	return retMap, nil
}

// UrlEncode encode the string for url passing
func UrlEncode(str string) string {
	// fmt.Println("in:", str)
	encodedUrl := url.QueryEscape(str)
	// fmt.Println("out:", encodedUrl)

	return encodedUrl
}

// UrlDecode decode the url encoded string to original string
func UrlDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// PrintMap print map in form of "k=v"
func PrintMap(m map[string]string) {
	for k, v := range m {
		fmt.Println(k, "=", v)
	}
}