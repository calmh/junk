package main

import "golang.org/x/time/rate"

func main() {
	r := rate.NewLimiter(0, 100)
	if r.Allow() {
		panic("should not allow")
	}
}
