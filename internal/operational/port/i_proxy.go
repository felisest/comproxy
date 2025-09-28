package port

type IProxy interface {
	Run() error
	SetProcessor(func([]byte, []byte))
}
