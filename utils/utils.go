package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"example.com/m/model"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
)

// 获取MD5sting
func MD5String(data []byte) string {
	return hex.EncodeToString(MD5(data))
}

// 获取MD5[]byte
func MD5(data []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	return md5Ctx.Sum(nil)
}

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

// 获取分页的 page size offNum
func GetFromDataPageSizeOffNum(c *gin.Context) (int, int, int) {
	num, pageOk := c.GetPostForm("page")  // 页数
	num2, sizeOK := c.GetPostForm("size") // 每页条数
	if !(pageOk && sizeOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return -1, -1, -1
	}
	page, _ := strconv.Atoi(num)
	size, _ := strconv.Atoi(num2)
	// 跳过条数
	offNum := size * (page - 1)
	return page, size, offNum
}

// 判断数据是否为空字符
func IsStringEmpty(data string, isExist bool) bool {
	if isExist && data == "" {
		isExist = false
	}
	return isExist
}

// string 数组去重
func RemoveDuplicatesAndEmpty(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
