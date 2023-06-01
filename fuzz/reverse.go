package reverse


func Reverse(str string) string {
	data := []byte(str)
	for i, j := 0, len(data)-1; i < len(data)/2; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return string(data)
}
