package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

// formData请求
// url 请求路径
// params 非文件请求参数
// fileParams os.open之后的 *File文件
func NewFormDataRequest(url string, params map[string]interface{}, fileParams map[string]interface{}) (*http.Request, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	// 创建 string类型
	for filed, val := range params {
		switch v := val.(type) {
		case string:
			err := writer.WriteField(filed, v)
			if err != nil {
				return nil, err
			}
		case interface{}:

		}
	}
	// 创建file类型
	for filed, val := range fileParams {
		switch v := val.(type) {
		case *multipart.FileHeader:
			part, err := writer.CreateFormFile(filed, v.Filename)
			if err != nil {
				return nil, err
			}
			file, err := v.Open()
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		case *os.File:
			part, err := writer.CreateFormFile(filed, path.Base(v.Name()))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, v)
			if err != nil {
				return nil, err
			}
		}
	}

	_ = writer.Close()
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}
