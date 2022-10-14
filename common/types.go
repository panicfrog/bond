package common

import (
	"encoding/json"
	"errors"
	"github.com/rsocket/rsocket-go/payload"
)

type ApiType int

type ApiSchema struct {
	Name     string  `json:"name"`
	Request  string  `json:"request,omitempty"`
	Response string  `json:"response,omitempty"`
	Type     ApiType `json:"type"`
}

func NewApiSchema(name string, request string, response string, apiType ApiType) *ApiSchema {
	return &ApiSchema{
		Name:     name,
		Request:  request,
		Response: response,
		Type:     apiType,
	}
}

func NewApiSchemaData(data []byte) (*ApiSchema, error) {
	api := new(ApiSchema)
	err := json.Unmarshal(data, api)
	return api, err
}

func (api *ApiSchema) Data() ([]byte, error) {
	return json.Marshal(api)
}

type ApiEvent struct {
	ApiSchema
	EventType EventType
	From      string
}

func NewApiEvent(api ApiSchema, eventType EventType) *ApiEvent {
	return &ApiEvent{
		ApiSchema: api,
		EventType: eventType,
	}
}

type ApiService struct {
	ApiSchema `json:"apiSchema"`
	From      string `json:"from"`
}

func NewApiServiceFormData(data []byte) (*ApiService, error) {
	api := new(ApiService)
	err := json.Unmarshal(data, api)
	if err != nil {
		return nil, err
	} else {
		return api, nil
	}
}

type EventType int
type EventMetadata struct {
	From      string    `json:"from,omitempty"`
	EventType EventType `json:"eventType"`
	EventData string    `json:"eventData,omitempty"`
}

func NewEventMetadata(from string, eventType EventType) *EventMetadata {
	return &EventMetadata{
		From:      from,
		EventType: eventType,
		EventData: "",
	}
}

func NewEventMetadetaData(data []byte) (*EventMetadata, error) {
	metadata := new(EventMetadata)
	err := json.Unmarshal(data, metadata)
	return metadata, err
}

func (d *EventMetadata) Data() ([]byte, error) {
	return json.Marshal(d)
}

type ApiMetadata struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to"`
	PipeName string `json:"pipeName"`
}

func NewApiMetadataTo(to string, pipeName string) *ApiMetadata {
	return &ApiMetadata{
		To:       to,
		PipeName: pipeName,
	}
}

func GetApiMetadata(payload2 payload.Payload) (*ApiMetadata, error) {
	infoData, ok := payload2.Metadata()
	if !ok {
		return nil, errors.New("metadata not exit")
	}
	info, err := NewApiMetadata(infoData)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func NewApiMetadata(data []byte) (*ApiMetadata, error) {
	var info ApiMetadata
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (m *ApiMetadata) Data() ([]byte, error) {
	return json.Marshal(m)
}

type SoundOutData struct {
	SeverId     string `json:"severId"`
	ServiceName string `json:"serviceName"`
}

func (s *SoundOutData) Data() ([]byte, error) {
	return json.Marshal(s)
}

func NewSoundOutData(severId string, serviceName string) *SoundOutData {
	return &SoundOutData{
		SeverId:     severId,
		ServiceName: serviceName,
	}
}

func NewSoundOutDataFromData(data []byte) (*SoundOutData, error) {
	s := new(SoundOutData)
	err := json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	} else {
		return s, nil
	}

}
