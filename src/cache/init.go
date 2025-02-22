package cache

import (
	"github.com/trancecho/open-sdk/cache/types"
	"github.com/trancecho/open-sdk/config"
	"log"
	"sync"
)

var (
	dbs = make(map[string]types.Cache)
	mux sync.RWMutex
)

func InitCache() {
	sources := config.GetConfig().Caches
	for _, source := range sources {
		setCacheByKey(source.Key, mustCreateCache(source))
		if source.Key == "" {
			source.Key = "*"
		}
		log.Println("create cache", source.Key, "=>", source.IP, ":", source.PORT)
	}
}

func GetCache(key string) types.Cache {
	mux.Lock()
	defer mux.Unlock()
	return dbs[key]
}

func setCacheByKey(key string, cache types.Cache) {
	if key == "" {
		key = "*"
	}
	if GetCache(key) != nil {
		log.Fatalln("duplicate db key: ", key)
	}
	mux.Lock()
	defer mux.Unlock()
	dbs[key] = cache
}

func mustCreateCache(conf config.Cache) types.Cache {
	var creator = getCreatorByType(conf.Type)
	if creator == nil {
		log.Fatalln("fail to find creator for cache types:", conf.Type)
		return nil
	}
	cache, err := creator.Create(conf)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return cache
}
