package badman

// Source is interface of BlackList.
type Source interface {
	Download() chan *EntityQueue
}
