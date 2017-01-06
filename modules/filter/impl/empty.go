package impl

type EmptyFilter struct {
}

func (this *EmptyFilter) Open(fileName string) error {

	return nil
}

func (this *EmptyFilter) Close() error {

	return nil

}

func (filter *EmptyFilter) Exists(key []byte) bool {
	return false
}

func (filter *EmptyFilter) Add(key []byte) error {

	return nil
}

func (filter *EmptyFilter) Delete(key []byte) error {

	return nil
}
