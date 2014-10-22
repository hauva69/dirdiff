package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	//	"strings"
)

var mapDn1 = make(map[string]string)
var mapDn2 = make(map[string]string)
var mapDn map[string]string

func getMd5Sum(path string, info os.FileInfo) string {
	const fileChunk = 8192
	file, err := os.Open(path)
	if err != nil {
		log.Printf("%s", err.Error())
		return "ERROR"
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
	}
	//else if strings.Index(path, ".git") != -1 {
	//		return filepath.SkipDir
	//	}

	if !info.IsDir() {
		mapDn[path] = getMd5Sum(path, info)
	}

	return nil
}

func main() {
	dn1 := os.Args[1]
	mapDn = mapDn1
	err := filepath.Walk(dn1, walkFunc)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	dn2 := os.Args[2]
	mapDn = mapDn2
	err = filepath.Walk(dn2, walkFunc)
	if err != nil && err != filepath.SkipDir {
		log.Printf("%v\n", err)
	}

	for k, _ := range mapDn1 {
		fmt.Println(k)
	}

	for k, _ := range mapDn2 {
		fmt.Println(k)
	}

	for k, v := range mapDn1 {
		_, ok := mapDn2[dn2+k]
		if ok {
			fmt.Printf("%s=%s\n", k, v)
		} else {
			fmt.Printf("%s does not exist in %s\n", k, dn2)
		}
	}
}
