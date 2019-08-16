package main

import (
	"log"
	"os"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/nuker"
)

func main() {
	accessKeyId, ok := os.LookupEnv("ALIYUN_NUKE_ACCESS_KEY_ID")
	accessKeySecret, ok := os.LookupEnv("ALIYUN_NUKE_ACCESS_KEY_SECRET")

	if !ok {
		log.Fatal("credential error: ALIYUN_NUKE_ACCESS_KEY_ID and ALIYUN_NUKE_ACCESS_KEY_SECRET undefined")
	}

	currentAccount := account.Account{
		Credentials: account.Credentials{
			AccessKeyID:     accessKeyId,
			AccessKeySecret: accessKeySecret,
		},
	}

	deleted := nuker.NukeItAll(currentAccount)
	for _, resource := range deleted {
		log.Println(resource)
	}
}
