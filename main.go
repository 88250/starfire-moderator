package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.SetLogLevel("*", "warn")

	ds, err := BadgerDatastore("ds")
	if err != nil {
		panic(err)
	}
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}

	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")

	h, dht, err := SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
	)

	if err != nil {
		panic(err)
	}

	lite, err := New(ctx, ds, h, dht, nil)
	if err != nil {
		panic(err)
	}

	lite.Bootstrap(DefaultBootstrapPeers())

	c, _ := cid.Decode("QmQKMPwjMBHqV5G3rus4caXSteGgefFsihCwXNAB7fUCQo")
	rsc, err := lite.GetFile(ctx, c)

	if err != nil {
		panic(err)
	}
	defer rsc.Close()
	content, err := ioutil.ReadAll(rsc)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}
