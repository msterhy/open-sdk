package cache

import (
	"github.com/trancecho/open-sdk/cache/driver"
	"github.com/trancecho/open-sdk/cache/types"
	"github.com/trancecho/open-sdk/config"
)

type Creator interface {
	Create(conf config.Cache) (types.Cache, error)
}

func init() {
	typeMap["redis"] = driver.RedisCreator{}
}

var typeMap = make(map[string]Creator)

func getCreatorByType(cacheType string) Creator {
	return typeMap[cacheType]
}
