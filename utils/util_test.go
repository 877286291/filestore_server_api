package utils

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestFileMD5(t *testing.T) {
	file, err := os.Open("../operation.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(FileMD5(file))
}
func TestFileSha1(t *testing.T) {
	file, err := os.Open("../ITM_Situation.lookup")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(FileSha1(file))
}
