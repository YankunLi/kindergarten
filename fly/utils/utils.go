package utils

import (
	//	"errors"
	"fmt"
	//	"io"
)

//const MAX_SLICE_SIZE uint32 = 4092
//
//func read(file *File, buf []byte, count uint32) (uint32, error) {
//	if count < 0 {
//		return 0, errors.New("parameter error")
//	}
//
//	temp_buf := make([]byte, MAX_SLICE_SIZE, MAX_SLICE_SIZE)
//	var to_read = count
//
//	n, err := file.Read(buf)
//	for {
//		if err != nil && err != io.EOF {
//			return 0, errors.New("read fail")
//		}
//	}
//	n, err := file.Read(temp_buf)
//
//	return num, err
//
//}

func Show() {
	fmt.Println("import utils test")
}

func StringToArr(str string, s []byte) {
	temp := []byte(str)
	for i := 0; i < len(temp); i++ {
		s[i] = temp[i]
	}
}
