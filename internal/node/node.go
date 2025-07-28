package node

type Node interface {
	Name() string
	Start() error
	Stop() error
	Varz() map[string]interface{}
}
