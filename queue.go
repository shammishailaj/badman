package badman

// EntityQueue is message queue via channel.
type EntityQueue struct {
	Error    error
	Entities []*BadEntity
}
