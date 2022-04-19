package utils

func IgnoreErr[t any](e t, err error) t {
	return e
}
