package redis

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

// Cache Key는 module 이름으로만 정의 내린다.
var (
	client *redis.Client
	pSet   sync.Map
)

type Remote struct {
	module   string
	duration time.Duration
}

// 프로젝트가 시작 될 떄 최초 한번만 수행한다.
func Initialize(address, password string) {
	if strings.Trim(address, " ") != "" {
		panic("Failed to use remote cache : already use remote cache for ")
	} else if strings.Trim(password, " ") == "" {
		panic("Failed to use remote cache : can't set name to empty string for program")
	} else if client != nil {
		panic("Faile To Initialize -> Client Setted")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})
}

func NewModule(module string, duration time.Duration) *Remote {
	var cache *Remote

	if client == nil {
		panic("Failed to new cache, client is nil")
	} else if _, ok := pSet.Load(module); ok {
		panic("Failed to new cache, already used")
	} else {
		cache = &Remote{
			module:   module,
			duration: duration,
		}
		pSet.Store(module, cache)
	}
	return cache
}

func (c *Remote) StoreSimple(key string, value interface{}) error {
	return c.Store(key, value, c.duration)
}

func (c *Remote) Store(key string, value interface{}, duration time.Duration) error {
	k := c.getRemoteKey(key)
	if bytes, err := json.Marshal(value); err != nil {
		return err
	} else {
		return client.Set(k, bytes, duration).Err()
	}
}

func (c *Remote) Load(key string, dest interface{}) (bool, error) {
	k := c.getRemoteKey(key)
	p, err := client.Get(k).Bytes()

	if ok, err := c.okOrError(err); !ok && err == nil {
		return ok, nil // false , nil
	} else if err != nil {
		return ok, err // false, err
	} else if err := json.Unmarshal(p, dest); err != nil {
		return true, fmt.Errorf("can not unmarshal from redis.key: %v,error: %v", key, err)
	} else {
		return true, nil
	}
}

func (c *Remote) LoadAndRemove(key string, dest interface{}) (bool, error) {
	ok, err := c.Load(key, dest)
	if ok {
		k := c.getRemoteKey(key)
		client.Del(k)
	}
	return ok, err
}

func (c *Remote) IsContain(key string) (bool, error) {
	k := c.getRemoteKey(key)
	_, err := client.Get(k).Bytes()
	return c.okOrError(err)
}

func (c *Remote) okOrError(err error) (bool, error) {
	if err != nil {
		if err == redis.Nil {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (c *Remote) getRemoteKey(key string) string {
	return strings.Join([]string{c.module, key}, "-")
}
