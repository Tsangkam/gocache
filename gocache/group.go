package gocache

import (
	"fmt"
	"log"
	"sync"
)

// cache namespace
type Group struct {
	name      string
	getter    Getter
	mainBlock block
}

var (
	rwmu   sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	rwmu.Lock()
	defer rwmu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainBlock: block{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	rwmu.RLock()
	g := groups[name]
	rwmu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainBlock.get(key); ok {
		log.Println("[gocache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainBlock.add(key, value)
}
