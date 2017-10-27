package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	fmt.Println(tokens([]byte("this is some data")))
}

func tokens(l []byte) []string {
	now := time.Now().UTC()
	return []string{
		token(l, now),
		token(l, now.Add(-time.Hour)),
	}
}

const timeFormat = "2006-01-02T15"

func token(l []byte, t time.Time) string {
	hours := uint32(time.Now().Unix() / 3600)
	token := make([]byte, 4+len(l))
	binary.BigEndian.PutUint32(token, hours)
	copy(token[4:], l)
	fmt.Printf("%q\n", token)

	hash := sha256.Sum256(token)
	return hex.EncodeToString(hash[:])
}
