package dao

type Dao[T any] interface {
	GetData(model *T, code string, fn func() (T, error)) error
}
