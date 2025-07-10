package types

type CahceKey string

func (ck CahceKey) String() string {
	return string(ck)
}

const (
	AccountCacheKey CahceKey = "accountCache"
)
