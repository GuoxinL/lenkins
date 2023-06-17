package tools

import (
	"fmt"
	"io"
	"os"
)

func WriteFile(fileName string, content []byte) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, io.SeekEnd)
		_, err = f.WriteAt(content, n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}
