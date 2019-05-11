package main

import (
	"math/rand"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	sh := api.NewShell("localhost:5001")
	data := randString() + "\n"
	sh.PubSubPublish("test", data)
}

func randString() string {
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	l := rand.Intn(10) + 2

	var s string
	for i := 0; i < l; i++ {
		s += string([]byte{alpha[rand.Intn(len(alpha))]})
	}
	return s
}
