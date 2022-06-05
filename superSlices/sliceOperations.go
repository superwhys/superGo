package superSlices

func Reverse(slice []interface{}) []interface{} {
	ret := make([]interface{}, len(slice))
	ret = append(slice[:0:0], slice...)
	for i := len(ret)/2 - 1; i >= 0; i-- {
		opp := len(ret) - 1 - i
		ret[i], ret[opp] = ret[opp], ret[i]
	}
	return ret
}

func SliceContainString(s []string, s1 string) bool {
	for _, item := range s {
		if item == s1 {
			return true
		}
	}
	return false
}
