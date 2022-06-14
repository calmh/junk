package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

func main() {
	bs, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	bw := bufio.NewWriter(os.Stdout)
	var conv []byte
	for _, rune := range string(bs) {
		if len(conv) > 0 {
			if rune <= 0x7f {
				log.Printf("%x %s", conv, conv)
				bw.Write(conv)
				conv = nil
				bw.WriteRune(rune)
				continue
			}
			conv = append(conv, byte(rune))
			continue
		}
		switch rune {
		case 0xc2, 0xc3, 0xc4, 0xc5, 0xe2:
			conv = append(conv, byte(rune))
		default:
			bw.WriteRune(rune)
		}
	}
	bw.Flush()
}
