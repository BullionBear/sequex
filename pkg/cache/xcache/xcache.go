package xcache

type XCache interface {
	Set(key int64, data interface{}) error
	GetLatest(size int64) (interface{}, error)
	Size() uint64
	Clear() error
	RemoveOldest(size int64) error
	Close() error
}
