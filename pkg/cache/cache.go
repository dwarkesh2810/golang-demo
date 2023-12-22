package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/memcache"
)

var memoryCacheProvider cache.Cache
var fileCacheProvider *cache.FileCache
var memcacheProvider cache.Cache

func init() {
	memoryCacheProvider, _ = cache.NewCache("memory", `{"interval":60}`)
	fileCacheProvider = &cache.FileCache{
		CachePath:  "caches",
		FileSuffix: ".cache",
	}
	memcacheProvider, _ = cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)
}

func Set(cacheType, key string, value interface{}, expire time.Duration) error {
	switch cacheType {
	case "memory":
		return memoryCacheProvider.Put(context.Background(), key, value, expire)

	case "file":
		return fileCacheProvider.Put(context.Background(), key, value, expire)

	case "memcache":
		val, err := InterfaceToBytes(value)
		if err != nil {
			return err
		}
		return memcacheProvider.Put(context.Background(), key, val, expire)

	default:
		return errors.New("unsupported Method")
	}
}

func Get(cacheType, key string) (interface{}, error) {
	switch cacheType {
	case "memory":
		return memoryCacheProvider.Get(context.Background(), key)

	case "file":
		return fileCacheProvider.Get(context.Background(), key)

	case "memcache":
		return memcacheProvider.Get(context.Background(), key)

	default:
		return nil, errors.New("unsupported Method")
	}
}

func Delete(cacheType, key string) error {
	switch cacheType {
	case "memory":
		return memoryCacheProvider.Delete(context.Background(), key)

	case "file":
		return fileCacheProvider.Delete(context.Background(), key)

	case "memcache":
		return memcacheProvider.Delete(context.Background(), key)

	default:
		return errors.New("unsupported Method")
	}
}

func IsExist(cacheType, key string) (bool, error) {
	switch cacheType {
	case "memory":
		return memoryCacheProvider.IsExist(context.Background(), key)

	case "file":
		return fileCacheProvider.IsExist(context.Background(), key)

	case "memcache":
		return memcacheProvider.IsExist(context.Background(), key)

	default:
		return false, errors.New("unsupported Method")
	}
}

func ClearAll() error {
	var err error
	err = memcacheProvider.ClearAll(context.Background())
	if err != nil {
		return err
	}
	err = fileCacheProvider.ClearAll(context.Background())
	if err != nil {
		return err
	}
	err = memcacheProvider.ClearAll(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func InterfaceToBytes(data interface{}) ([]byte, error) {

	switch v := data.(type) {
	case string:
		return []byte(v), nil

	default:
		return json.Marshal(data)
	}
}
