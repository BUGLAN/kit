package listutil

func GroupList(datas []string, size int) [][]string {
	if size == 0 {
		return [][]string{}
	}
	length := len(datas)
	if length < size {
		return [][]string{datas}
	}

	arrs := make([][]string, 0)

	lastIndex := 0
	for i := 0; i < length; i++ {
		if i-lastIndex == size {
			arrs = append(arrs, datas[lastIndex:i])
			lastIndex = i
		}
	}
	arrs = append(arrs, datas[lastIndex:])
	return arrs
}
