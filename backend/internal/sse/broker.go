package sse

import "sync"

type Broker struct {
	Clients map[string]chan string
	Mu      sync.RWMutex
}

var MainBroker = &Broker{
	Clients: make(map[string]chan string),
}

func (b *Broker) AddClient(userID string) chan string {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	ch := make(chan string)
	b.Clients[userID] = ch
	return ch
}

func (b *Broker) RemoveClient(userID string) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	if ch, ok := b.Clients[userID]; ok {
		close(ch)
		delete(b.Clients, userID)
	}
}
