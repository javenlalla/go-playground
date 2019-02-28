package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"log"
)

func main() {
	s := "test-string-to-hash"

	shaHash := shaByString(s)
	log.Println(shaHash)

	mdfiveHash := mdfiveByString(s)
	log.Println(mdfiveHash)
}

func shaByString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func mdfiveByString(s string) string {
	h := md5.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil))
}
