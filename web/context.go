package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Context struct {
	Resp       http.ResponseWriter
	Req        *http.Request
	PathParams map[string]string
}

func (c *Context) BindJSON(val any) error {
	if val == nil {
		return errors.New("web:输入不能为nil")
	}
	//bs,_:=io.ReadAll(c.Req.Body)
	//json.Unmarshal(bs,val)
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}
