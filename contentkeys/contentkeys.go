package contentkeys

type contextKey string

const (
	ServiceTypeKey contextKey = "serviceType"
	DBKey          contextKey = "db"
	CacheKey       contextKey = "cache"
)
