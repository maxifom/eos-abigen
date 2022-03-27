package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/maxifom/eos-abigen/examples/contracts/eosio/go/generated/eosio"
	"github.com/maxifom/eos-abigen/pkg/client"
)

func main() {
	rpcClient, err := client.NewRPCClient("https://eos.greymass.com", &http.Client{
		Timeout: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	eosioClient := eosio.NewClient(rpcClient)
	producers, err := eosioClient.Producers(context.TODO())
	if err != nil {
		panic(err)
	}

	for _, row := range producers.Rows {
		log.Println(row.Owner, time.Time(row.LastClaimTime))
	}

	eosioRes, err := eosioClient.Userres(context.TODO())
	if err != nil {
		panic(err)
	}

	for _, row := range eosioRes.Rows {
		log.Println(row.Owner, row.NetWeight, row.CpuWeight, row.RamBytes)
	}

	eosioTokenRes, err := eosioClient.Userres(context.TODO(), client.Scope("eosio.token"))
	if err != nil {
		panic(err)
	}

	for _, row := range eosioTokenRes.Rows {
		log.Println(row.Owner, row.NetWeight, row.CpuWeight, row.RamBytes)
	}
}
