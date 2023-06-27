package tools

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	index := strings.Index(dir, "lenkins")

	fmt.Println(index)
	return dir[:index+7]
}

func SetCurrentProjectWorkingDir() {
	path := GetCurrentPath()
	fmt.Println("working dir: ", path)
	_ = os.Chdir(path)
}
