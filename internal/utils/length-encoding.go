package utils

import "encoding/binary"

func ToByteArr32(len int32, base int) (arr []byte) {
	//stored inside the arr
	binary.BigEndian.PutUint32(arr[0:base], uint32(len))
	return
}

func ToByteArr64(len int64, base int) (arr []byte) {
	//stored inside the arr
	binary.BigEndian.PutUint64(arr[0:base], uint64(len))
	return
}

func FromByteArr32(arr []byte, base int) (len uint32) {
	return binary.BigEndian.Uint32(arr[0:base])
}

func FromByteArr64(arr []byte, base int) (len uint64) {
	return binary.BigEndian.Uint64(arr[0:base])
}
