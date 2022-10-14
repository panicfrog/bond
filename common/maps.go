package common

import (
	"errors"
	"github.com/rsocket/rsocket-go"
	"sync"
)

type Topic struct {
	OwnerId    string
	Name       string
	Parameters []byte
}

type BondSocket struct {
	Id     string
	Socket rsocket.CloseableRSocket
}

type MemberMap struct {
	id2socket map[string]*BondSocket
	socket2id map[*BondSocket]string
	lock      sync.RWMutex
}

func NewMemberMap() *MemberMap {
	return &MemberMap{
		id2socket: make(map[string]*BondSocket),
		socket2id: make(map[*BondSocket]string),
		lock:      sync.RWMutex{},
	}
}

func (m *MemberMap) Add(socket *BondSocket) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.id2socket[socket.Id]; ok {
		return errors.New("is exit")
	}
	if _, ok := m.socket2id[socket]; ok {
		return errors.New("is exit")
	}
	m.id2socket[socket.Id] = socket
	m.socket2id[socket] = socket.Id
	return nil
}

func (m *MemberMap) Remove(socket *BondSocket) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if id, ok := m.socket2id[socket]; ok {
		delete(m.id2socket, id)
	}
	delete(m.socket2id, socket)
}

func (m *MemberMap) RemoveById(id string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if socket, ok := m.id2socket[id]; ok {
		delete(m.socket2id, socket)
	}
	delete(m.id2socket, id)
}

func (m *MemberMap) GetID(socket *BondSocket) (string, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	id, ok := m.socket2id[socket]
	return id, ok
}

func (m *MemberMap) GetSocket(id string) (*BondSocket, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	socket, ok := m.id2socket[id]
	return socket, ok
}

func (m *MemberMap) DealSocket(id string, f func(*BondSocket)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if socket, ok := m.id2socket[id]; ok {
		f(socket)
	}
}

func (m *MemberMap) ForEachSocket(f func(socket *BondSocket)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, s := range m.id2socket {
		f(s)
	}
}

func (m *MemberMap) DealId(socket *BondSocket, f func(string)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if id, ok := m.socket2id[socket]; ok {
		f(id)
	}
}

type TopicMap struct {
	topic2sockets  map[*Topic]map[*BondSocket]interface{}
	sockets2topics map[*BondSocket]map[*Topic]interface{}
	lock           sync.RWMutex
}

func NewTopicMap() *TopicMap {
	t2s := make(map[*Topic]map[*BondSocket]interface{})
	s2t := make(map[*BondSocket]map[*Topic]interface{})
	return &TopicMap{
		topic2sockets:  t2s,
		sockets2topics: s2t,
		lock:           sync.RWMutex{},
	}
}

func (m *TopicMap) Add(topic *Topic, socket *BondSocket) {
	m.lock.Lock()
	defer m.lock.Unlock()
	socketsSet, ok := m.topic2sockets[topic]
	if !ok {
		socketsSet = make(map[*BondSocket]interface{})
		m.topic2sockets[topic] = socketsSet
	}
	socketsSet[socket] = struct{}{}
	topicSet, ok := m.sockets2topics[socket]
	if !ok {
		topicSet = make(map[*Topic]interface{})
		m.sockets2topics[socket] = topicSet
	}
	topicSet[topic] = struct{}{}
}

func (m *TopicMap) RemoveSocket(topic *Topic, socket *BondSocket) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if socketSet, ok := m.topic2sockets[topic]; ok {
		delete(socketSet, socket)
	}
	delete(m.sockets2topics, socket)
}

func (m *TopicMap) Remove(socket *BondSocket, topic *Topic) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if topicSet, ok := m.sockets2topics[socket]; ok {
		delete(topicSet, topic)
	}
	delete(m.topic2sockets, topic)
}

func (m *TopicMap) DealSockets(topic *Topic, f func(socket *BondSocket)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if socketSet, ok := m.topic2sockets[topic]; ok {
		for s, _ := range socketSet {
			f(s)
		}
	}
}

func (m *TopicMap) DealTopics(socket *BondSocket, f func(topic *Topic)) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if topicSet, ok := m.sockets2topics[socket]; ok {
		for t, _ := range topicSet {
			f(t)
		}
	}
}
