package hashcache

type HashCache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Remove(key string) error
	Clear() error
	Size() uint64
}
