package port

type IComparer interface {
	Compare([]byte, []byte) (bool, string)
}
