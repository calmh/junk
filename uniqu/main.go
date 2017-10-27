package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	seen := map[string]int{}
	br := bufio.NewScanner(os.Stdin)
	for br.Scan() {
		key := strings.TrimSpace(br.Text())
		seen[key]++
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%d\t%s\n", seen[key], key)
	}
}
