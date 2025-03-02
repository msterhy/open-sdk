package captcha

import (
	"time"

	"pinnacle-primary-be/core/store/rds"
)

type (
	Redis struct {
		rds       *rds.Redis
		expireSec int
		baseName  string
	}
)

func NewRedisStore(rds *rds.Redis, expire time.Duration, baseName string) *Redis {
	return &Redis{
		rds:       rds,
		expireSec: int(expire.Seconds()),
		baseName:  baseName,
	}
}

func (r *Redis) Set(id string, value string) error {
	return r.rds.Setex(r.baseName+"/"+id, value, r.expireSec)
}

func (r *Redis) Get(id string, clear bool) string {
	value, _ := r.rds.Get(r.baseName + "/" + id)
	if clear {
		r.rds.Del(r.baseName + "/" + id)
	}
	return value
}

func (r *Redis) Verify(id, answer string, clear bool) bool {
	v := r.Get(id, clear)
	return v == answer
}
