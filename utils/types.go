package utils

func ToStrings(t []interface{}) []string {
	st := make([]string, len(t))
	for i, e := range t {
		st[i] = e.(string)
	}

	return st
}
