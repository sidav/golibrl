package string_operations

func ReverseString(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}

func GetRotatedStringArray(arr *[]string) *[]string { // rotates 90 degrees clockwise
	newArr := make([]string, 0)
	for i := 0; i < len((*arr)[0]); i++ {
		str := ""
		for j := 0; j < len(*arr); j++ {
			str += string((*arr)[j][i])
		}
		newArr = append(newArr, str)
	}
	return &newArr
}

func GetMirroredStringArray(arr *[]string, v, h bool) *[]string {
	newArr := make([]string, 0)
	if v && h {
		for i := len(*arr) - 1; i >= 0; i-- {
			newArr = append(newArr, ReverseString((*arr)[i]))
		}
		return &newArr
	}
	if v {
		for i := len(*arr) - 1; i >= 0; i-- {
			newArr = append(newArr, (*arr)[i])
		}
		return &newArr
	}
	if h {
		for i := 0; i < len(*arr); i++ {
			newArr = append(newArr, ReverseString((*arr)[i]))
		}
		return &newArr
	}
	return arr // nil
}