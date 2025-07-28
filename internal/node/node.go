package node

type Node interface {
	Type() string
	Name() string
	Start() error
	Stop() error
	Varz() map[string]interface{}
}
