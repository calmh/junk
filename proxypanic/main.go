package main

import (
	"fmt"
	"net/http"
)

func main() {
	cli := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	_, err := cli.Get("https://google.com/")
	if err != nil {
		fmt.Println(err)
	}
}
