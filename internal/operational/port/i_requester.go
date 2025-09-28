package port

type IRequester interface {
	Post(request []byte) ([]byte, error)
}
