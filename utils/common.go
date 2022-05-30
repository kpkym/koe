package utils

func IgnoreErr[t any](e t, _ error) t {
	return e
}
