package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"

	api "github.com/ipfs/go-ipfs-api"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	topic         = "starfire"
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

	userHome, err := homedir.Dir()
	if nil != err {
		panic(err)
	}
	configPath := filepath.Join(userHome, ".ipfs-starfire-moderator", "config")
	data, err := ioutil.ReadFile(configPath)
	if nil != err {
		panic(err)
	}
	config := map[string]interface{}{}
	if err := json.Unmarshal(data, &config); nil != err {
		panic(err)
	}
	identity := config["Identity"].(map[string]interface{})
	privKey := identity["PrivKey"].(string)
	privKeyBytes, err := base64.StdEncoding.DecodeString(privKey)
	if nil != err {
		panic(err)
	}
	key, err := crypto.UnmarshalPrivateKey(privKeyBytes)
	if nil != err {
		panic(err)
	}
	moderateBlacklistCmd := map[string]interface{}{
		"type": typeBlacklist,
		"data": blacklistId,
	}
	signBytes, err := key.Sign([]byte(blacklistId))
	if nil != err {
		panic(err)
	}
	sign := hex.EncodeToString(signBytes)
	moderateBlacklistCmd["sign"] = sign
	moderateBlacklistCmdBytes, err := json.Marshal(moderateBlacklistCmd)
	if nil != err {
		panic(err)
	}
	moderateBlacklistCmdData := string(moderateBlacklistCmdBytes)
	if err = sh.PubSubPublish(topic, moderateBlacklistCmdData); nil != err {
		panic(err)
	}

	fmt.Println("home publishing")
	start := time.Now()
	if err = sh.Publish("", homeId); nil != err {
		panic(err)
	}
	end := time.Now()
	publishElapsed := end.Sub(start).Seconds()
	fmt.Printf("home published [ipns/%s] in [%.2f]\n", id.ID, publishElapsed)

	homeResolved, err := sh.Resolve(id.ID)
	if nil != err {
		panic(err)
	}
	fmt.Println("home resolved [" + homeResolved + "]")
}
