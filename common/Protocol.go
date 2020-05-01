package common

import (
	"encoding/json"
	"net/http"
)

//job的结构体
type Job struct {
	Name     string `json:"name"`      // name
	Command  string `json:"command"`   //shell命令
	CronExpr string `json:"cron_expr"` //cron表达式
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ApiResponse(w http.ResponseWriter,code int,msg string,data interface{}) error {
	var (
		res Response
		ret []byte
		err error
	)

	res.Code = code
	res.Msg  = msg
	res.Data = data

	ret,err = json.Marshal(res)
	if err != nil {
		return err
	}
	//设置响应头
	w.Header().Set("Content-Type","application/json")
	_,err = w.Write(ret)
	if err != nil {
		return err
	}

	return nil
}
