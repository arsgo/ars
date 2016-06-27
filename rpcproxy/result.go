package rpcproxy

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	result_error_format   = `{"code":"%s","msg":"%s"}`
	result_success_format = `{"code":"success","msg":"成功"}`
	result_data_format    = `{"code":"success","msg":"成功","data":"%s"}`
)

//ResultEntity 结果实体
type ResultEntity struct {
	Code string `json:"code"`
}

//ResultIsSuccess 检查当前result是否为成功
func ResultIsSuccess(content string) bool {
	entity := &ResultEntity{}
	err := json.Unmarshal([]byte(content), &entity)
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.ToLower(entity.Code), "success")
}

func GetErrorResult(code string, msg ...interface{}) string {
	return fmt.Sprintf(result_error_format, code, fmt.Sprint(msg...))
}
func GetSuccessResult() string {
	return result_success_format
}

func GetDataResult(data string) string {
	if strings.EqualFold(data, "") || strings.EqualFold(strings.ToLower(data), "nil") || strings.EqualFold(strings.ToLower(data), "null") {
		return result_success_format
	}
	if strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}") {
		return data
	}
	if strings.HasPrefix(data, "[{") && strings.HasSuffix(data, "}]") {
		return data
	}
	return fmt.Sprintf(result_data_format, data)

}
