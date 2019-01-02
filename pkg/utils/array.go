package utils

func DeleteFromArray(e string, arr []string) []string {
	res := []string{}
	for _, a := range arr {
		if a != e {
			res = append(res, a)
		}
	}

	return res
}

func UniqueArray(arr []string) []string {
	uniq := map[string]struct{}{}
	res := []string{}
	for _, a := range arr {
		if _, ok := uniq[a]; !ok {
			res = append(res, a)
			uniq[a] = struct{}{}
		}
	}

	return res
}

func ArrayInclude(arr []string, e string) bool {
	for _, a := range arr {
		if a == e {
			return true
		}
	}
	return false
}
