package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
)

var mapDn1 = make(map[string]string)
var mapDn2 = make(map[string]string)

func getMd5Sum(path string, info os.FileInfo) string {
	const fileChunk = 8192
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
	defer file.Close()
	blocks := uint64(math.Ceil(float64(info.Size()) / float64(fileChunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(fileChunk,
			float64(info.Size()-int64(i*fileChunk))))
		buf := make([]byte, blockSize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	} else if path == ".git" {
		return filepath.SkipDir
	}
	fmt.Println(path)
	mapDn1[path] = getMd5Sum(path, info)
	fmt.Println(mapDn1[path])

	return nil
}

func main() {
	dn1 := os.Args[1]
	err := filepath.Walk(dn1, walkFunc)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
