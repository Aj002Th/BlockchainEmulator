package stateStorage

// StateStorage 针对状态树的持久化存储
type StateStorage interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Has([]byte) bool
	Delete([]byte) error
	BatchPut([][2][]byte) error
	Close() error
}
