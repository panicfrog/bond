package client

import (
	"bond/common"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"github.com/santhosh-tekuri/jsonschema"
	"log"
	"strings"
	"sync"
)

var river actionRiver
var metadataLock sync.RWMutex
var metadataPushImpl MetadataPush

func init() {
	river = actionRiver{
		reqResPipes:     sync.Map{},
		fireForgetPipes: sync.Map{},
		reqStreamPipes:  sync.Map{},
	}
	metadataLock = sync.RWMutex{}
}

type Entity struct {
	client   rsocket.Client
	identify string
	host     string
	port     int
}

type actionRiver struct {
	reqResPipes     sync.Map
	fireForgetPipes sync.Map
	reqStreamPipes  sync.Map
}

type reqResAction struct {
	requestSchema  *jsonschema.Schema
	responseSchema *jsonschema.Schema
	pipe           ReqResPipe
}

type reqStreamAction struct {
	requestSchema  *jsonschema.Schema
	responseSchema *jsonschema.Schema
	errorChan      chan error
	payloadChan    chan payload.Payload
	pipe           RequestStreamPipe
}

type fireForgetAction struct {
	requestSchema *jsonschema.Schema
	pipe          FireForgetPipe
}

func (e *Entity) Close() error {
	defer func() {
		river.reqResPipes.Range(func(key, value any) bool {
			river.reqResPipes.Delete(key)
			return true
		})
		river.fireForgetPipes.Range(func(key, value any) bool {
			river.fireForgetPipes.Delete(key)
			return true
		})
	}()
	return e.client.Close()
}

type ReqResPipe interface {
	ReqResFlow(input []byte) ([]byte, error)
}

type FireForgetPipe interface {
	FireForgetFlow(input []byte)
}

type RequestStreamPipe interface {
	RequestStream(input []byte) ([]byte, error)
}

type MetadataPush interface {
	ServiceRegister(input []byte)
	ServiceUnRegister(input []byte)
}

func (e *Entity) RegisterFireForgetPipe(name string, requestSchema string, action FireForgetPipe) error {
	compiler := jsonschema.NewCompiler()
	var (
		requestSch *jsonschema.Schema
	)
	if len(requestSchema) != 0 {
		if err := compiler.AddResource("request.json", strings.NewReader(requestSchema)); err != nil {
			return err
		}
		_requestSch, err := compiler.Compile("request.json")
		if err != nil {
			return err
		}
		requestSch = _requestSch
	}
	pipeAction := fireForgetAction{
		requestSchema: requestSch,
		pipe:          action,
	}
	if _, ok := river.fireForgetPipes.LoadOrStore(name, pipeAction); ok {
		return errors.New(fmt.Sprintf("'%s' is exit", name))
	}

	defer func() {
		api := common.NewApiSchema(name, requestSchema, "", common.ApiTypeFireForget)
		if err := e.notifyRegister(api); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (e *Entity) RegisterReqResPipe(name string, requestSchema string, responseSchema string, action ReqResPipe) error {
	compiler := jsonschema.NewCompiler()
	var (
		requestSch  *jsonschema.Schema
		responseSch *jsonschema.Schema
	)
	if len(requestSchema) != 0 {
		if err := compiler.AddResource("request.json", strings.NewReader(requestSchema)); err != nil {
			return err
		}
		_requestSch, err := compiler.Compile("request.json")
		if err != nil {
			return err
		}
		requestSch = _requestSch
	}

	if len(responseSchema) != 0 {
		if err := compiler.AddResource("response.json", strings.NewReader(responseSchema)); err != nil {
			return err
		}
		_responseSch, err := compiler.Compile("response.json")
		if err != nil {
			return err
		}
		responseSch = _responseSch
	}

	pipeAction := reqResAction{requestSch, responseSch, action}
	if _, ok := river.reqResPipes.LoadOrStore(name, pipeAction); ok {
		return errors.New(fmt.Sprintf("'%s' is exit", name))
	}

	defer func() {
		api := common.NewApiSchema(name, requestSchema, responseSchema, common.ApiTypeRequestResponse)
		if err := e.notifyRegister(api); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (e *Entity) OnMetadataPush(callback MetadataPush) {
	metadataLock.Lock()
	defer metadataLock.Unlock()
	metadataPushImpl = callback
	if metadataPushImpl == nil {
		log.Println("set error metadataPushImpl")
	}
	log.Println("set metadatapush function")
}

func (e *Entity) notifyRegister(api *common.ApiSchema) error {
	data, err := api.Data()
	if err != nil {
		return err
	}
	m := common.NewEventMetadata("", common.EventTypeRegister)
	m.EventData = string(data)
	metadata, err := m.Data()
	if err != nil {
		return err
	}
	e.client.MetadataPush(payload.New(data, metadata))
	return nil
}

func (e *Entity) notifyUnRegister(api *common.ApiSchema) error {
	data, err := api.Data()
	if err != nil {
		return err
	}
	m := common.NewEventMetadata("", common.EventTypeUnRegister)
	m.EventData = string(data)
	metadata, err := m.Data()
	if err != nil {
		return err
	}

	e.client.MetadataPush(payload.New(data, metadata))
	return nil
}

type StreamSink struct {
	action *reqStreamAction
}

func (s *StreamSink) Validate(data []byte) error {
	if s.action.requestSchema != nil {
		return s.action.requestSchema.Validate(bytes.NewReader(data))
	}
	return nil
}

func (s *StreamSink) Sink(data []byte) error {
	if s.action.responseSchema != nil {
		if err := s.action.responseSchema.Validate(bytes.NewReader(data)); err != nil {
			return err
		}
	}
	s.action.payloadChan <- payload.New(data, nil)
	return nil
}

func (s *StreamSink) Error(err error) {
	s.action.errorChan <- err
}

func (e *Entity) RegisterReqStreamPipe(name string, requestSchema string, responseSchema string, action RequestStreamPipe) (*StreamSink, error) {
	compiler := jsonschema.NewCompiler()
	var (
		requestSch  *jsonschema.Schema
		responseSch *jsonschema.Schema
	)
	if len(requestSchema) != 0 {
		if err := compiler.AddResource("request.json", strings.NewReader(requestSchema)); err != nil {
			return nil, err
		}
		_requestSch, err := compiler.Compile("request.json")
		if err != nil {
			return nil, err
		}
		requestSch = _requestSch
	}

	if len(responseSchema) != 0 {
		if err := compiler.AddResource("response.json", strings.NewReader(responseSchema)); err != nil {
			return nil, err
		}
		_responseSch, err := compiler.Compile("response.json")
		if err != nil {
			return nil, err
		}
		responseSch = _responseSch
	}
	errChan := make(chan error)
	payloadChan := make(chan payload.Payload)
	pipeAction := reqStreamAction{
		requestSchema:  requestSch,
		responseSchema: responseSch,
		errorChan:      errChan,
		payloadChan:    payloadChan,
		pipe:           action,
	}

	if _, ok := river.reqStreamPipes.LoadOrStore(name, pipeAction); ok {
		return nil, errors.New(fmt.Sprintf("'%s' is exit", name))
	}

	sink := &StreamSink{
		action: &pipeAction,
	}

	defer func() {
		api := common.NewApiSchema(name, requestSchema, requestSchema, common.ApiTypeRequestStream)
		if err := e.notifyRegister(api); err != nil {
			log.Println(err)
		}
	}()

	return sink, nil
}

func (e *Entity) RequestResponse(name string, to string, parameters []byte) ([]byte, error) {
	data, err := common.NewApiMetadataTo(to, name).Data()
	if err != nil {
		return nil, err
	}
	response, err := e.client.RequestResponse(payload.New(parameters, data)).Block(context.Background())
	if err != nil {
		log.Println("RequestResponse " + err.Error())
		return nil, err
	}
	return response.Data(), nil
}

func (e *Entity) FireForget(name string, to string, parameters []byte) {
	data, err := common.NewApiMetadataTo(to, name).Data()
	if err != nil {
		log.Println(err)
		return
	}
	e.client.FireAndForget(payload.New(parameters, data))
}

type StreamResponse interface {
	Response([]byte)
	Error(error)
}

func (e *Entity) RequestStream(name string, to string, parameters []byte, responseF StreamResponse) error {
	data, err := common.NewApiMetadataTo(to, name).Data()
	if err != nil {
		return err
	}
	e.client.RequestStream(payload.New(parameters, data)).
		DoOnNext(func(input payload.Payload) error {
			responseF.Response(input.Data())
			return nil
		}).
		DoOnError(func(e error) {
			responseF.Error(e)
		})
	return nil
}

func (e *Entity) SoundOutService(sever string, service string) error {
	data, err := common.NewSoundOutData(sever, service).Data()
	if err != nil {
		return err
	}
	m := common.NewEventMetadata(e.identify, common.EventSoundOut)
	m.EventData = string(data)
	metadata, err := m.Data()
	if err != nil {
		return err
	}
	e.client.MetadataPush(payload.New(data, metadata))
	return nil
}

func reqRes(request payload.Payload) (response mono.Mono) {
	info, err := common.GetApiMetadata(request)
	if err != nil {
		return mono.Error(errors.Wrap(err, "reqRes metadata error"))
	}
	value, ok := river.reqResPipes.Load(info.PipeName)
	if !ok {
		return mono.Error(errors.New("action not exit"))
	}
	action, ok := value.(reqResAction)
	if !ok {
		return mono.Error(errors.New("pipe type error"))
	}
	if action.requestSchema != nil {
		if err := action.requestSchema.Validate(bytes.NewReader(request.Data())); err != nil {
			return mono.Error(errors.Wrap(err, "reqRes request validate error"))
		}
	}
	res, err := action.pipe.ReqResFlow(request.Data())
	if err != nil {
		return mono.Error(err)
	}
	if action.responseSchema != nil {
		if err := action.responseSchema.Validate(bytes.NewReader(res)); err != nil {
			return mono.Error(errors.Wrap(err, "reqRes response validate error"))
		}
	}
	responseMeta, err := common.NewApiMetadataTo(info.From, info.PipeName).Data()
	if err != nil {
		return mono.Error(errors.Wrap(err, "reqRes response metadata error"))
	}
	return mono.Just(payload.New(res, responseMeta))
}

func requestStream(request payload.Payload) flux.Flux {
	info, err := common.GetApiMetadata(request)
	if err != nil {
		return flux.Error(err)
	}
	value, ok := river.reqStreamPipes.Load(info.PipeName)
	if !ok {
		return flux.Error(errors.New("action not exit"))
	}
	action, ok := value.(reqStreamAction)
	if !ok {
		return flux.Error(errors.New("pipe type error"))
	}
	if action.responseSchema != nil {
		if err := action.responseSchema.Validate(bytes.NewReader(request.Data())); err != nil {
			return flux.Error(err)
		}
	}
	return flux.CreateFromChannel(action.payloadChan, action.errorChan)
}

func fireForget(request payload.Payload) {
	info, err := common.GetApiMetadata(request)
	if err != nil {
		log.Println(err)
		return
	}
	value, ok := river.fireForgetPipes.Load(info.PipeName)
	if !ok {
		log.Println(errors.New("action not exit"))
		return
	}
	action, ok := value.(fireForgetAction)
	if !ok {
		log.Println(errors.New("action type error"))
		return
	}
	if action.requestSchema != nil {
		if err := action.requestSchema.Validate(bytes.NewReader(request.Data())); err != nil {
			log.Println(err)
			return
		}
	}
	action.pipe.FireForgetFlow(request.Data())
}

func metadataPush(request payload.Payload) {
	metadata, ok := request.Metadata()
	if !ok {
		return
	}
	eventMatadata, err := common.NewEventMetadetaData(metadata)
	if err != nil {
		log.Println(err)
		return
	}
	metadataLock.RLock()
	defer metadataLock.RUnlock()
	if metadataPushImpl == nil {
		return
	}
	api, err := common.NewApiSchemaData([]byte(eventMatadata.EventData))
	if err != nil {
		log.Println(err)
		return
	}
	service := common.ApiService{
		ApiSchema: *api,
		From:      eventMatadata.From,
	}
	data, err := json.Marshal(service)
	if err != nil {
		return
	}
	if eventMatadata.EventType == common.EventTypeRegister {
		metadataPushImpl.ServiceRegister(data)
	} else if eventMatadata.EventType == common.EventTypeUnRegister {
		metadataPushImpl.ServiceUnRegister(data)
	}
}

func Setup(identify string, host string, port int, onMetadataPush MetadataPush) (*Entity, error) {
	v, err := common.AesEncryptCBC([]byte(identify), []byte(common.EncryptKey))
	if err != nil {
		return nil, err
	}
	setupPayload := payload.New(v, []byte("auth"))
	service := rsocket.TCPClient().SetHostAndPort(host, port).Build()

	metadataLock.Lock()
	metadataPushImpl = onMetadataPush
	metadataLock.Unlock()

	acceptor := func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(reqRes),
			rsocket.FireAndForget(fireForget),
			rsocket.RequestStream(requestStream),
			rsocket.MetadataPush(metadataPush),
		)
	}

	cli, err := rsocket.Connect().
		SetupPayload(setupPayload).
		Acceptor(acceptor).
		Transport(service).
		Start(context.Background())

	return &Entity{cli, identify, host, port}, err
}
