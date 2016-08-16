package base

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ERR_NOT_FIND_SRVS = "999"
)

//ResultEntity 结果实体
type ResultEntity struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    string `json:"data"`
}

func (r ResultEntity) ToJson() string {
	buffer, err := json.Marshal(r)
	if err != nil {
		return err.Error()
	}
	return string(buffer)
}

//ResultIsSuccess 检查当前result是否为成功
func ResultIsSuccess(content string) bool {
	entity := &ResultEntity{}
	err := json.Unmarshal([]byte(content), &entity)
	if err != nil {
		return false
	}
	return strings.EqualFold(entity.Code, "success")
}
func GetResult(content string) (entity ResultEntity) {
	entity = ResultEntity{}
	err := json.Unmarshal([]byte(content), &entity)
	if err != nil {
		entity = ResultEntity{Code: "-1"}
	}
	return
}

func GetErrorResult(code string, msg ...interface{}) string {
	message := fmt.Sprint(msg...)
	r := ResultEntity{Code: code, Message: message}
	return r.ToJson()
}
func GetSuccessResult() string {
	r := ResultEntity{Code: "success", Message: "成功"}
	return r.ToJson()
}

func IsRaw(types map[string]string) bool {
	return types["_original"] == "true"
}

func GetDataResult(cdata string, redirect bool) string {
	if redirect {
		return cdata
	}
	data := strings.Trim(cdata, " ")
	if strings.EqualFold(data, "") || strings.EqualFold(data, "nil") || strings.EqualFold(data, "null") {
		return GetSuccessResult()
	}
	if strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}") {
		return data
	}
	if strings.HasPrefix(data, "[{") && strings.HasSuffix(data, "}]") {
		return data
	}
	if strings.HasPrefix(data, "<?xml") || strings.HasPrefix(data, "<html>") || strings.HasPrefix(data, "<xml") {
		return data
	}
	r := ResultEntity{Code: "success", Message: "成功", Data: data}
	return r.ToJson()
}
