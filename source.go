package badman

// Source is interface of BlackList.
type Source interface {
	Load() chan *BadEntityMessage
}

// DefaultSources is default set of blacklist source that is maintained by badman.
var DefaultSources = []Source{}
