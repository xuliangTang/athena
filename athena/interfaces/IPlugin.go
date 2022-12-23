package interfaces

type IPlugin interface {
	Enabler() bool
	InitModule()
}
