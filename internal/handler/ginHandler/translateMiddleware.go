package ginHandler

import (
	"bytes"
	"encoding/json"
	"ipfast_server/internal/config/i18n"

	"github.com/gin-gonic/gin"
)

// 翻译中间件
func TranslateMiddleware(c *gin.Context) {
	lang := c.GetHeader("Accept-Language")
	I18n := i18n.NewLocalizer(lang)
	c.Writer = &CustomResponseWriter{c.Writer, &bytes.Buffer{}}
	c.Next()
	var resdata = struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{}

	if writer, ok := c.Writer.(*CustomResponseWriter); ok {
		err := json.Unmarshal(writer.Body.Bytes(), &resdata)
		if err != nil {
			writer.ResponseWriter.Write(writer.Body.Bytes())
		}
		resdata.Msg = I18n.F(resdata.Msg)
		data, err := json.Marshal(resdata)
		if err != nil {
			writer.ResponseWriter.Write(writer.Body.Bytes())
		}
		writer.ResponseWriter.Write(data)
	}
}
