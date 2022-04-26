package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"fmt"
	"sort"
	"strings"
)

func ApiAuth(apiKey, sign string, param interface{}) error {

	if apiKey == "" {
		return errors.New("api key is empty")
	}
	if sign == "" {
		return errors.New("sign is empty")
	}

	strRet, err := json.Marshal(param)
	if err != nil {
		log.Error("json.Marshal(param) error(%v)", err)
		return errors.New("接口验证失败")
	}

	// json转map
	var mRet map[string]interface{}
	err1 := json.Unmarshal(strRet, &mRet)
	if err1 != nil {
		log.Error("json.UnMarshal(param) error(%v)", err)
		return errors.New("接口验证失败")
	}

	rightApiKey := app.Conf.JiaGu.ApiKey
	apiSecret := app.Conf.JiaGu.ApiSecret

	// 1.验证apiKey
	if apiKey != rightApiKey {
		return errors.New("apiKey 信息有误")
	}

	// 2.验证签名
	rightSign := hmacSha1(apiSecret, concatParam(mRet, apiKey))
	if sign != rightSign {
		return errors.New("签名信息有误")
	}
	return nil
}

func hmacSha1(secret, text string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}

// 获取 sign
func concatParam(m map[string]interface{}, apiKey string) string {
	result := apiKey
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		result = result + fmt.Sprintf("%v", m[k])
	}
	return strings.Trim(result, "&")
}
