package util

import (
	"bytes"
	"fmt"
)

//ByteFormat 函数将二进制流转换成类似ue编辑器的格式显示
//in 表示输入流， line_width表示一行显示多少个字符
func ByteFormat(in []byte, line_width int) string {
	if line_width < 1 {
		line_width = 16
	}

	if in == nil {
		return ""
	}

	in_temp := in[:]
	line_num := len(in) / line_width

	//刚好是整数行则需处理多出一行的问题
	if len(in)%line_width == 0 {
		line_num = line_num - 1
	}

	//预分配内存
	//全部
	out := bytes.NewBuffer(make([]byte, 0))
	out.Grow((15 + line_width*5) * (line_num + 1))
	//16进制行
	xBuf := bytes.NewBuffer(make([]byte, 0))
	//字符形式行
	cBuf := bytes.NewBuffer(make([]byte, 0))

	for i := 0; i <= line_num; i++ {
		xBuf.Reset()
		xBuf.Grow(15 + line_width*3)
		xBuf.Reset()
		cBuf.Grow(5 + line_width*2)

		out.WriteString(fmt.Sprintf("%08d:  ", i))
		//最后一行特殊处理
		if i == line_num {
			for _, c := range in_temp {
				xBuf.WriteString(fmt.Sprintf("%02x ", c))
				if c > 31 && c < 127 {
					cBuf.WriteString(fmt.Sprintf("%c ", c))
				} else {
					cBuf.WriteString(". ")
				}
			}
			for i := 0; i < line_width-len(in_temp); i++ {
				xBuf.WriteString("   ")
			}
			xBuf.WriteString(";  ")
			out.WriteString(xBuf.String())
			out.WriteString(cBuf.String())
			return out.String()
		}
		for _, c := range in_temp[:line_width] {
			xBuf.WriteString(fmt.Sprintf("%02x ", c))
			if c > 31 && c < 127 {
				cBuf.WriteString(fmt.Sprintf("%c ", c))
			} else {
				cBuf.WriteString(". ")
			}
		}
		xBuf.WriteString(";  ")
		out.WriteString(xBuf.String())
		out.WriteString(cBuf.String())
		out.WriteString("\n")
		in_temp = in_temp[line_width:]
	}
	return out.String()
}

func BinHexOutput(in []byte) string {
	return ByteFormat(in, 16)
}
