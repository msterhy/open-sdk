package rdsx

import (
	"github.com/trancecho/open-sdk/cache"
	"github.com/trancecho/open-sdk/cache/types"
	"log"
)

var Cache types.Cache

func InitCache() {
	Cache = cache.GetCache("MainRedis")
	if Cache == nil {
		log.Fatalln("fail to get cache")
	}
}
