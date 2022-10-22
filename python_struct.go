package IntBinary

import (
	"bytes"
	"math"
	"strings"
)

// StructBinary 整数转二进制 python struct?
// pack unpack
// b,h,l
type StructBinary struct{}

func (receiver *StructBinary) BytesToInt(body []byte, symbolic bool) int64 {
	// 二进制转整数
	//symbolic == true: 负数
	var Number, s int64
	var length = len(body)
	for i := range body {
		var y = float64(length - i - 1)
		var num = int64(math.Pow(256, y))
		Number += num * int64(body[i])
	}
	if symbolic {
		s = int64(math.Pow(256, float64(length)))
	}
	return Number - s
}

func (receiver *StructBinary) IntToBytes(n int, size int) []byte {
	// 整数转二进制
	// size长度
	var PowNumber = int64(math.Pow(256, float64(size)))
	var Number = PowNumber + int64(n)
	var body = make([]byte, size)
	var remainder = Number % 256
	for i := range body {
		body[size-i-1] = uint8(remainder)
		Number = Number / 256
		remainder = Number % 256
	}
	return body
}

func (receiver *StructBinary) PackValue(Format string, value []int) []byte {
	// 封包
	var body = make([][]byte, len(value))
	for i := 0; i < len(Format); i++ {
		var Fmt = string(Format[i])
		switch Fmt {
		case "B", "b":
			body[i] = receiver.IntToBytes(value[i], 1)
		case "H", "h":
			body[i] = receiver.IntToBytes(value[i], 2)
		case "L", "l":
			body[i] = receiver.IntToBytes(value[i], 4)
		}
	}
	return bytes.Join(body, []byte(""))
}

func (receiver *StructBinary) UnPackValue(Format string, value []byte) []int64 {
	// 解包
	var now int64 = 0
	var body []int64
	for i := 0; i < len(Format); i++ {
		var Fmt = string(Format[i])
		switch Fmt {
		case "B", "b":
			var SliBody = value[now : now+1]
			var jut = strings.ToUpper(Fmt) != Fmt
			var num = receiver.BytesToInt(SliBody, jut)
			body = append(body, num)
			now += 1
		case "H", "h":
			var SliBody = value[now : now+2]
			var jut = strings.ToUpper(Fmt) != Fmt
			var num = receiver.BytesToInt(SliBody, jut)
			body = append(body, num)
			now += 2
		case "L", "l":
			var SliBody = value[now : now+4]
			var jut = strings.ToUpper(Fmt) != Fmt
			var num = receiver.BytesToInt(SliBody, jut)
			body = append(body, num)
			now += 4
		}
	}
	return body
}

func Pack(Format string, value []int) []byte {
	// StructBinary pack
	return new(StructBinary).PackValue(Format, value)
}

func UnPack(Format string, value []byte) []int64 {
	// StructBinary unpack
	return new(StructBinary).UnPackValue(Format, value)
}
