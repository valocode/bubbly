package resource

type provider interface {
	Query(string) (string, error)
	Save(string, string) error
}

// ProviderType is a data access layer provider
type ProviderType string

const (
	// Etcd is an Etcd provider
	Etcd ProviderType = "etcd"
	// Buntdb is a Buntdb provider
	Buntdb ProviderType = "buntdb"
)
