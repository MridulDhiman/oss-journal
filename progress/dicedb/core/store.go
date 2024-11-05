package core

import (
	"time"

	"github.com/MridulDhiman/dice/config"
)

var store map[string]*Obj


type Obj struct {
	Value interface{};
	ExpiresAt int64;
}

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMs int64) *Obj {
	var expiresAt int64 = -1
	if(durationMs > 0) {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	return &Obj{
		Value: value,
		ExpiresAt: expiresAt,
	}
}


func evictFirst() {
	for k:= range store {
		delete(store, k)
		return
	}
}

func evict () {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	}
}

func Put(k string, obj *Obj) {
	if len(store) >= config.KeysLimit {
		evict()
	}
	store[k] = obj;
}

func Get(k string) *Obj {
	v:= store[k];
	if v != nil {
		// check if obj is expired or not:
			// check if expiry set and if yes, is it expired
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli() {
			// if key is expired: delete the key
				delete(store, k)
				return nil
		}
	}
	return v
}

func Del(k string) bool {
	if _, ok:= store[k]; ok {
		delete(store, k)
		return true
	}
	return false
}

