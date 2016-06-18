package rpcproxy

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	result_error_format   = `{"code":"%s","msg":"%s"}`
	result_success_format = `{"code":"100","msg":"success"}`
	result_data_format    = `{"code":"100","msg":"success","data":"%s"}`
)

//ResultEntity 结果实体
type ResultEntity struct {
	Code string
	Msg  string
}

//ResultIsSuccess 检查当前result是否为成功
func ResultIsSuccess(content string) bool {
	entity := &ResultEntity{}
	err := json.Unmarshal([]byte(content), &entity)
	if err != nil {
		return false
	}
	return strings.EqualFold(entity.Msg, "success")
}

func GetErrorResult(code string, msg ...interface{}) string {
	return fmt.Sprintf(result_error_format, code, fmt.Sprint(msg...))
}
func GetSuccessResult() string {
	return result_success_format
}

func GetDataResult(data string) string {
	if strings.EqualFold(data, "") || strings.EqualFold(data, "nil") {
		return result_success_format
	}
	if strings.HasSuffix(data, "}") &&
		strings.HasPrefix(data, "{") {
		return data
	}
	return fmt.Sprintf(result_data_format, data)
}
