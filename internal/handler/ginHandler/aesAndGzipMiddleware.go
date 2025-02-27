package ginHandler

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"encoding/json"
	"fmt"
	"io"
	"ipfast_server/internal/handler/aesHandler"
	"ipfast_server/pkg/util/log"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
AES加解密密钥和向量
*/
const (
	key = "0edafd4bcb0da1ba"
	iv  = "e5b2a779b34464df"
)

/*
检查数据长度是否是 AES 块大小的整数倍

	param:
		data []byte 数据
	return:
		bool 是否是 AES 块大小的整数倍
*/
func checkDataLength(data []byte) bool {
	blockSize := aes.BlockSize // AES 块大小为 16 字节
	return len(data)%blockSize == 0
}

/*
压缩 Gzip 数据

	param:
		data []byte 数据
	return:
		[]byte 压缩后的数据
		error 错误信息
*/
func compressGzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

/*
解压缩 Gzip 数据

	param:
		data []byte Gzip 数据
	return:
		[]byte 解压缩后的数据
		error 错误信息
*/
func decompressGzip(data []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	return io.ReadAll(gz)
}

/*
AES加解密中间件

	param:
		NeedEncryption: 响应是否需要加密
		NeedDecryption: 请求是否需要解密
	return:
		gin.HandlerFunc Gin中间件
*/
func AESAndGzipMiddleware(c *gin.Context) {
	var err error
	var bodyBytes []byte
	isGzip := strings.Contains(c.GetHeader("Content-Encoding"), "gzip")
	encryption := c.Request.Header.Get("Encryption") == "true"
	// 使用自定义响应写入器（如果需要）
	if isGzip || encryption {
		c.Writer = &CustomResponseWriter{Body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	}
	// 解密请求体（如果需要）
	if encryption {
		bodyBytes, err = getRequestBodyBytes(c)
		if err != nil {
			Failed("001", c, isGzip, encryption, err) // 读取请求体失败
			return
		}
		if !checkDataLength(bodyBytes) {
			Failed("002", c, isGzip, encryption, fmt.Errorf("aes data length is error:%d", len(bodyBytes))) // 数据长度不是 AES 块大小的整数倍
			return
		}
		bodyBytes, err = aesHandler.AesDecrypt(bodyBytes, []byte(key), []byte(iv))
		if err != nil {
			Failed("003", c, isGzip, encryption, err) // 解密失败
			return
		}
	} else {
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			Failed("004", c, isGzip, encryption, err) // 读取请求体失败
			return
		}
	}

	// 解压缩（如果需要）
	if isGzip {
		bodyBytes, err = decompressGzip(bodyBytes)
		if err != nil {
			Failed("005", c, isGzip, encryption, err) // 解压缩失败
			return
		}
		c.Request.Header.Del("Content-Encoding")
	}

	// 替换请求体
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	c.Next()
	handleResponse(c, isGzip, encryption, nil)
}

/*
处理响应

	param:
		c *gin.Context Gin上下文
		isGzip bool 是否压缩
		encryption bool 是否加密
*/
func handleResponse(c *gin.Context, isGzip, encryption bool, responseData []byte) {
	if writer, ok := c.Writer.(*CustomResponseWriter); ok {
		if responseData == nil {
			responseData = writer.Body.Bytes()
		}
		var err error
		// 压缩响应（如果需要）
		if isGzip {
			responseData, err = compressGzip(responseData)
			if err != nil {
				_, err1 := writer.ResponseWriter.Write([]byte("Failed to compress response data"))
				if err1 != nil {
					log.Error("Failed to write response data:", err)
				}
				c.Abort()
				return
			}
		}

		// 加密响应（如果需要）
		if encryption {
			responseData, err = aesHandler.AesEncrypt(responseData, []byte(key), []byte(iv))
			if err != nil {
				_, err1 := writer.ResponseWriter.Write([]byte("Failed to compress response data"))
				if err1 != nil {
					log.Error("Failed to write response data:", err)
				}
				c.Abort()
				return
			}
			writer.Header().Set("Encryption", "true")
		}

		_, err1 := writer.ResponseWriter.Write(responseData)
		if err1 != nil {
			log.Error("Failed to write response data:", err)
		}
		c.Abort()
		return
	} else {
		if responseData != nil {
			c.JSON(200, responseData)
			c.Abort()
			return
		}
	}
}

/*
获取请求体字节
*/
func getRequestBodyBytes(c *gin.Context) (bodyBytes []byte, err error) {
	bodyBytes, err = io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return
}

/*
返回失败响应 code 1 和 Response中的Res数据
*/
func Failed(errstr string, c *gin.Context, isGzip, encryption bool, err error) {
	var res = make(map[string]interface{})
	res["code"] = 1
	res["data"] = map[string]string{
		"error": errstr,
	}
	resBytes, err1 := json.Marshal(res)
	if err1 != nil {
		log.Error("Failed to marshal response data:", err)
		return
	}
	log.Error(err.Error() + ":" + errstr)
	handleResponse(c, isGzip, encryption, []byte(resBytes))
}
