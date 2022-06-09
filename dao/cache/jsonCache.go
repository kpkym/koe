package cache

type jsonCache[T any] struct{}

func (j jsonCache[T]) Get(key string) (value T, ok bool) {
	// TODO implement me
	panic("implement me")
}

func (j jsonCache[T]) Set(key string, value T) {
	// TODO implement me
	panic("implement me")
}

func (j jsonCache[T]) GetOrSet(key string, fn func() (value T)) (value T) {
	// TODO implement me
	panic("implement me")
}

func NewJsonCache[T any]() jsonCache[T] {
	return jsonCache[T]{}
}
