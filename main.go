package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	api "github.com/ipfs/go-ipfs-api"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	sh := api.NewShell("localhost:5002")
	sh.PubSubPublish("starfire", randString()+"\n")

	blacklist, err := os.Open("blacklist")
	if nil != err {
		panic(err)
	}

	blacklistId, err := sh.Add(blacklist)
	if nil != err {
		panic(err)
	}
	fmt.Println("blacklist added [" + blacklistId + "]")

	moderateBlacklistCmd := map[string]interface{}{
		"type": "blacklist",
		"data": blacklistId,
	}
	moderateBlacklistCmdBytes, err := json.Marshal(moderateBlacklistCmd)
	if nil != err {
		panic(err)
	}
	moderateBlacklistCmdData := string(moderateBlacklistCmdBytes)
	sh.PubSubPublish("starfire", moderateBlacklistCmdData)
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
