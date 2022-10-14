package client

import (
	"bond/common"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var listers sync.Map

func init() {
	listers = sync.Map{}
}

type consolePipe struct{}

func (*consolePipe) FireForgetFlow(s []byte) {
	// cli
	fmt.Println(string(s))
}

func (e *Entity) RegisterConsole() error {
	// cli
	return e.RegisterFireForgetPipe(common.ConsolePipeName, "", &consolePipe{})
}

func (e *Entity) Console(str string) {
	// client
	listers.Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			e.FireForget(common.ConsolePipeName, k, []byte(str))
		}
		return true
	})
}

type ConsoleResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ConsoleRequest struct {
	User string `json:"user"`
}

func (*consolePipe) ReqResFlow(req []byte) ([]byte, error) {
	// TODO: 改成从请求本身去获取
	request := new(ConsoleRequest)
	err := json.Unmarshal(req, request)
	if err != nil {
		return nil, err
	}
	if len(request.User) == 0 {
		return nil, errors.New("user should not empty")
	}
	listers.Store(request.User, "")
	response := ConsoleResponse{
		Code:    "S0001",
		Message: "success",
	}
	return json.Marshal(&response)
}

func (e *Entity) RegisterConsoleListener() error {
	var (
		// https://transform.tools/json-to-json-schema
		resScheme = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"title\": \"response.json\",\"type\": \"object\",\"properties\": {\"message\": {\"type\": \"string\"},\"code\": {\"type\": \"string\"}},\"required\": [\"message\",\"code\"]}"
	)
	return e.RegisterReqResPipe(common.ConsoleListenerPipeName, "", resScheme, &consolePipe{})
}
