package rpcclient

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

func getErrorResult(code string, msg string) string {
	return fmt.Sprintf(result_error_format, code, msg)
}
func getSuccessResult() string {
	return result_success_format
}

func getDataResult(data string) string {
	return fmt.Sprintf(result_data_format, data)
}
