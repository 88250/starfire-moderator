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

const (
	topic = "starfire"
	typeBlacklist = "blacklist"
)

func main() {
	sh := api.NewShell("localhost:5002")

	id, err := sh.ID()
	if nil != err {
		panic(err)
	}
	fmt.Println("moderator id [" + id.ID + "]")

	homeId, err := sh.AddDir("home")
	if nil != err {
		panic(err)
	}
	fmt.Println("home [" + homeId + "]")

	blacklist, err := os.Open("home/blacklist")
	if nil != err {
		panic(err)
	}
	blacklistId, err := sh.Add(blacklist)
	if nil != err {
		panic(err)
	}
	fmt.Println("blacklist [" + blacklistId + "]")

	moderateBlacklistCmd := map[string]interface{}{
		"type": typeBlacklist,
		"data": blacklistId,
	}
	moderateBlacklistCmdBytes, err := json.Marshal(moderateBlacklistCmd)
	if nil != err {
		panic(err)
	}
	moderateBlacklistCmdData := string(moderateBlacklistCmdBytes)
	sh.PubSubPublish(topic, moderateBlacklistCmdData)

	err = sh.Publish(id.ID, homeId)
	if nil != err {
		panic(err)
	}
}

