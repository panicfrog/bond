package main

import (
	"bond/common"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

var (
	memebers *common.MemberMap
	borkerKV clientv3.KV
)

func init() {
	memebers = common.NewMemberMap()
}

func insertService(kv clientv3.KV, user string, api *common.ApiSchema) error {
	data, err := api.Data()
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s/%s/%s", common.RootPrefix, user, api.Name)
	_, err = kv.Put(context.Background(), key, string(data))
	return err
}

func deleteService(kv clientv3.KV, user string, api *common.ApiSchema) error {
	key := fmt.Sprintf("%s/%s/%s", common.RootPrefix, user, api.Name)
	_, err := kv.Delete(context.Background(), key)
	return err
}

func deleteUserService(kv clientv3.KV, user string) error {
	key := fmt.Sprintf("%s/%s/", common.RootPrefix, user)
	_, err := kv.Delete(context.Background(), key, clientv3.WithPrefix())
	return err
}

func getService(kv clientv3.KV, user string, apiName string) ([]string, error) {
	key := fmt.Sprintf("%s/%s/%s", common.RootPrefix, user, apiName)
	response, err := kv.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	if len(response.Kvs) > 0 {
		for _, res := range response.Kvs {
			result = append(result, string(res.Value))
		}
	}
	return result, nil
}

func getUserService(kv clientv3.KV, user string) ([]string, error) {
	keyPrefix := fmt.Sprintf("%s/%s/", common.RootPrefix, user)
	response, err := kv.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	if len(response.Kvs) > 0 {
		for _, res := range response.Kvs {
			result = append(result, string(res.Value))
		}
	}
	return result, nil
}

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"etcd1:2379"},
		//Endpoints:   []string{"172.16.232.77:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatalln(err)
	}

	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			log.Println(errors.Wrap(err, "close etcd client error"))
		}
	}(cli)

	kv := clientv3.NewKV(cli)
	borkerKV = kv
	rootPrefix := fmt.Sprintf("%s/", common.RootPrefix)
	watchChan := cli.Watch(context.Background(), rootPrefix, clientv3.WithPrefix())

	apiEventChan := make(chan *common.ApiEvent, 500)

	go func() {
		for {
			select {
			case response, ok := <-watchChan:
				if ok {
					for _, event := range response.Events {
						if event.Type == clientv3.EventTypePut {
							api, err := common.NewApiSchemaData(event.Kv.Value)
							if err == nil {
								var evt *common.ApiEvent
								if clientv3.EventTypePut == event.Type {
									evt = common.NewApiEvent(*api, common.EventTypeRegister)
								} else if clientv3.EventTypeDelete == event.Type {
									evt = common.NewApiEvent(*api, common.EventTypeUnRegister)
								}
								slice := strings.Split(string(event.Kv.Key), "/")
								evt.From = slice[1]
								if evt != nil {
									apiEventChan <- evt
								}
							} else {
								log.Println(errors.Wrap(err, "create api schema error"))
							}
						} else if event.Type == clientv3.EventTypeDelete {
							// TODO: 处理删除
						}
					}
				} else {
					log.Fatalln("watchChan close")
				}
			case evt, ok := <-apiEventChan:
				if ok {
					memebers.ForEachSocket(func(s *common.BondSocket) {
						data, err := evt.ApiSchema.Data()
						if err != nil {
							log.Println(errors.Wrap(err, "encode api schema error"))
							return
						}
						m := common.NewEventMetadata(evt.From, evt.EventType)
						m.EventData = string(data)
						m.From = evt.From
						metadata, err := m.Data()
						if err != nil {
							log.Println(errors.Wrap(err, "create event metadata error"))
							return
						}
						log.Printf("api event push %s", string(metadata))
						s.Socket.MetadataPush(payload.New(data, metadata))
					})
				}
			}
		}
	}()

	err = rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			data := setup.Data()
			_id, err := common.AesDecryptCBC(data, []byte(common.EncryptKey))
			id := string(_id)
			if err != nil {
				return nil, errors.Wrap(err, "unAuthorization")
			}
			log.Println(id + " is connected")
			sendingSocket.OnClose(func(err error) {
				log.Println(id + " is disconnected")
				if _err := deleteUserService(borkerKV, id); _err != nil {
					log.Println(errors.Wrap(_err, "delete user service error"))
				}
				memebers.RemoveById(id)
			})
			s := &common.BondSocket{
				Id:     id,
				Socket: sendingSocket,
			}
			if err := memebers.Add(s); err != nil {
				return nil, err
			}
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(request payload.Payload) mono.Mono {
					info, err := common.GetApiMetadata(request)
					if err != nil {
						return mono.Error(errors.Wrap(err, "get api metadata error"))
					}
					to, ok := memebers.GetSocket(info.To)
					if !ok {
						return mono.Error(errors.New(fmt.Sprintf("%s is disconnet", info.To)))
					}
					info.From = s.Id
					responseMetadata, err := info.Data()
					if err != nil {
						return mono.Error(errors.New("inner error"))
					}
					return to.Socket.RequestResponse(payload.New(request.Data(), responseMetadata))
				}),
				rsocket.FireAndForget(func(request payload.Payload) {
					info, err := common.GetApiMetadata(request)
					if err != nil {
						log.Println(errors.Wrap(err, "create api metadata error"))
						return
					}
					to, ok := memebers.GetSocket(info.To)
					if !ok {
						log.Println(fmt.Sprintf("%s is disconnected", info.To))
						return
					}
					info.From = s.Id
					responseMetadata, err := info.Data()
					if err != nil {
						return
					}
					to.Socket.FireAndForget(payload.New(request.Data(), responseMetadata))
				}),
				rsocket.MetadataPush(func(request payload.Payload) {
					md, ok := request.Metadata()
					if !ok {
						return
					}
					metadata, merr := common.NewEventMetadetaData(md)
					if merr != nil {
						log.Println(errors.Wrap(merr, "metadata push create event metadata error"))
						return
					}
					if metadata.EventType == common.EventSoundOut {
						if soundOutValue, serr := common.NewSoundOutDataFromData([]byte(metadata.EventData)); serr == nil {
							var service []string
							if len(soundOutValue.ServiceName) > 0 { // 获取用户一个服务
								service, serr = getService(borkerKV, soundOutValue.SeverId, soundOutValue.ServiceName)
								if serr != nil {
									return
								}
							} else { // 获取用户所有服务
								service, serr = getUserService(borkerKV, soundOutValue.SeverId)
								if serr != nil {
									return
								}
							}
							event := common.NewEventMetadata(soundOutValue.SeverId, common.EventTypeRegister)

							for _, sv := range service {
								event.EventData = sv
								metadata, merr := event.Data()
								if merr != nil {
									log.Println(errors.Wrap(merr, "event sound out, create event metadata error"))
									return
								}
								s.Socket.MetadataPush(payload.New([]byte(sv), metadata))
							}
						}
						return
					} else {
						api, err := common.NewApiSchemaData([]byte(metadata.EventData))
						if err != nil {
							log.Print(errors.Wrap(err, "event sound out, create api schema error"))
							log.Println(metadata.EventData)
							return
						}

						if metadata.EventType != common.EventTypeRegister && metadata.EventType != common.EventTypeUnRegister {
							return
						}
						if metadata.EventType == common.EventTypeRegister {
							if err := insertService(borkerKV, id, api); err != nil {
								log.Println(errors.Wrap(err, "event sound out, insert service error"))
								return
							}
						} else if metadata.EventType == common.EventTypeUnRegister {
							if err := deleteService(borkerKV, id, api); err != nil {
								log.Println(errors.Wrap(err, "event sound out, delete service error"))
								return
							}
						}
					}
				}),
				rsocket.RequestStream(func(request payload.Payload) flux.Flux {
					info, err := common.GetApiMetadata(request)
					if err != nil {
						return flux.Error(err)
					}
					to, ok := memebers.GetSocket(info.To)
					if !ok {
						return flux.Error(errors.New(fmt.Sprintf("%s is disconnet", info.To)))
					}
					info.From = s.Id
					responseMetadata, err := info.Data()
					if err != nil {
						return flux.Error(errors.New("inner error"))
					}
					return to.Socket.RequestStream(payload.New(request.Data(), responseMetadata))
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
		Serve(context.Background())
	log.Fatalln(err)
}
