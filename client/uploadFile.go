package client

import (
	"bond/common"
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type uploadFilePipe struct {
	handler UploadFileHandler
}
type UploadFileHandler interface {
	HandleUploadFile(request *FileUploadRequest, error error)
}

type FileUploadRequest struct {
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	Flag     string `json:"flag"`
	FileData []byte
}

type FileInfo struct {
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	Flag     string `json:"flag"`
}

type FileUploadResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Data    FileInfo `json:"data"`
}

func (r *FileUploadRequest) Encode() ([]byte, error) {
	var result bytes.Buffer
	enc := gob.NewEncoder(&result)
	err := enc.Encode(r)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func (u *uploadFilePipe) ReqResFlow(input []byte) ([]byte, error) {
	dec := gob.NewDecoder(bytes.NewReader(input))
	var request = new(FileUploadRequest)
	err := dec.Decode(&request)
	if err != nil {
		if u.handler != nil {
			u.handler.HandleUploadFile(nil, err)
		}
		return nil, err
	}
	if u.handler != nil {
		u.handler.HandleUploadFile(request, nil)
	}
	response := FileUploadResponse{
		Code:    "S0001",
		Message: "success",
		Data: FileInfo{
			Flag:     request.Flag,
			FileName: request.FileName,
			FileSize: request.FileSize},
	}
	return json.Marshal(response)
}

func (e *Entity) RegisterUploadFile(handler UploadFileHandler) error {
	return e.RegisterReqResPipe(common.UploadFilePipeName, "", "", &uploadFilePipe{handler: handler})
}
