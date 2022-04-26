package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"example.com/m/model"
	"example.com/m/pkg/log"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"unicode"
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

// formData请求
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

// 发送json请求
func NewJsonRequest(url string, m interface{}) (*http.Response, error) {
	data, err := jsoniter.Marshal(m)
	if err != nil {
		log.Error(err.Error())
	}
	response, err := http.Post(url, "application/json", bytes.NewReader(data))
	return response, nil
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
func GetURLDataPageSizeOffNum(c *gin.Context) (int, int, int) {
	num, pageOk := c.GetQuery("page")  // 页数
	num2, sizeOK := c.GetQuery("size") // 每页条数
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

// 定时轮询清理数据
func TickerDel(orgPath, delPath string, clock time.Duration) {
	// 创建文件路径
	err := os.MkdirAll(orgPath, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	err = os.MkdirAll(delPath, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	// 设置定时器进行轮询清理满足条件的文件
	ticker := time.NewTicker(clock)
	for range ticker.C {
		log.Info("-----开启本轮清理-----")
		err := DelDataByTime(orgPath,
			delPath)
		if err != nil {
			log.Error(err)
		}
		log.Info("-----完成本轮清理-----")
	}
}

// 定时删除服务器中的数据方法 movpath初始位置  deldir待删除位置
func DelDataByTime(movpath, deldir string) error {
	now := time.Now()
	// 超过30天放入待删除路径
	err := filepath.Walk(movpath, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// 文件真实路径
		relpath, _ := filepath.Rel(movpath, path)
		//movTime := time.Second
		movTime := 30 * time.Hour * 24
		if isMov := now.Sub(info.ModTime()) > movTime; isMov {
			afterPath := deldir + "/" + relpath
			err := os.MkdirAll(filepath.Dir(afterPath), os.ModePerm)
			if err != nil {
				log.Error(err)
			}
			err = os.Rename(path, afterPath)
			if err != nil {
				log.Error(err)
			}
			// 删除初始文件夹
			err = os.Remove(filepath.Dir(path))
			if err != nil {
				log.Error(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}

	// 超过365天进行清除
	err = filepath.Walk(deldir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		//delTime := 1 * time.Second
		delTime := 365 * time.Hour * 24
		if isMov := now.Sub(info.ModTime()) > delTime; isMov {
			err := os.Remove(path)
			if err != nil {
				log.Error(err)
			}
		}
		err = os.Remove(filepath.Dir(path))
		if err != nil {
			log.Error(err)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func SpecialString(str string) string {
	var chars []rune
	for _, letter := range str {
		ok, letters := SpecialLetters(letter)
		if ok {
			chars = append(chars, letters...)
		} else {
			chars = append(chars, letter)
		}
	}
	return string(chars)
}
func SpecialLetters(letter rune) (bool, []rune) {
	if unicode.IsPunct(letter) || unicode.IsSymbol(letter) || unicode.Is(unicode.Han, letter) {
		var chars []rune
		chars = append(chars, '\\', letter)
		return true, chars
	}
	return false, nil
}
