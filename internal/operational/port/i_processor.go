package port

type IProcessor interface {
	Process(data ...[]byte) error
}
