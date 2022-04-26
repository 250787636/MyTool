package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"net/http"
)

// 发送json请求
// m  map[sting]interface{}
// url 请求地址
func NewJsonRequest(url string, m interface{}) (*http.Response, error) {
	data, err := jsoniter.Marshal(m)
	if err != nil {
		//log.Error(err.Error())
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(data))
	return response, nil
}