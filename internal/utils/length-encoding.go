package utils

import "encoding/binary"

func ToByteArr(len int) (arr []byte) {
	//stored inside the arr
	binary.BigEndian.PutUint32(arr[0:4], uint32(len))
	return
}

func FromByteArr(arr []byte) (len int) {
	return int(binary.BigEndian.Uint32(arr[0:4]))
}
