package util

type DeduplicatePlugin interface {
	Init(fileName string) error
	Persist() error
	Lookup(key []byte) bool
	Add(key []byte) error
}
