package encoding

import "fmt"

/*
utf-8是一种针对Unicode的可变长度字符编码，也是一种前缀码。
它可以用一至四个字节对Unicode字符集中的所有有效编码点进行编码，
属于Unicode标准的一部分

对于某一个字符的UTF-8编码，如果只有一个字节则其最高二进制位为0,与ASCII兼容；
如果是多字节，其第一个字节从最高位开始，连续的二进制位值为1的个数决定了其编码的位数，其余各字节均以10开头。
具体表示如下：
0XXX_XXXX
110X_XXXX 10XX_XXXX
1110_XXXX 10XX_XXXX 10XX_XXXX
1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
*/

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}

func isUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}
