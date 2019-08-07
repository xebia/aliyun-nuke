package main

import (
	"fmt"
	"os"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/nuker"
)

func main() {
	currentAccount := account.Account{
		Credentials: account.Credentials{
			AccessKeyID:     os.Getenv("ALIYUN_NUKE_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("ALIYUN_NUKE_ACCESS_KEY_SECRET"),
		},
	}

	deleted := nuker.NukeItAll(currentAccount)
	for _, resource := range deleted {
		fmt.Println(resource)
	}
}
