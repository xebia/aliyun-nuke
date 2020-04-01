# aliyun-nuke

Clears out all resources in a given Alibaba Cloud account. **Use with extreme caution!**

## Deletion process

Because Alibaba Cloud has many regions and every region can contain different resources, aliyun-nuke loops over all services and all regions to bruteforce the deletion of the resources. If a resource cannot be deleted because a dependent resource still exists, the dependent resource is deleted first and the process is started again. This repeats until no more resources were deleted for a single run.

## Supported services

Currently supported services are:

| Service | Elements                       |
| ------- | ------------------------------ |
| ECS     | Instances, security groups     |
| OSS     | Buckets, objects               |
| VPC     | VPCs, VSwitches, NAT gateways  |
| RAM     | Users, groups, roles, policies |

Any other resources will be kept as-is. If any unsupported resources block the deletion of the above resource types, aliyun-nuke will stop the deletion process and quit.

## Getting started as CLI tool

Build aliyun-cli from the source code by running `go build .` in the root directory of the repository. This assumes your Go environment is set up correctly.

When you have built an executable, define the following environment variables:

```bash
export ALIYUN_NUKE_ACCESS_KEY_ID=--YOUR ACCESS KEY ID--
export ALIYUN_NUKE_ACCESS_KEY_SECRET=--YOUR ACCESS KEY SECRET--
```

You can find these on the [AccessKey](https://ak-console.aliyun.com/) page in the Alibaba Cloud console. If you haven't created access keys before you might have to create them.

Then run aliyun-nuke to clear the account:

```bash
./aliyun-nuke
```

## Library usage

aliyun-nuke can also be used as a library:

```go
package main

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/nuker"
)

func main() {
	currentAccount := account.Account{
		Credentials: account.Credentials{
			AccessKeyID:     "--YOUR ACCESS KEY ID--",
			AccessKeySecret: "--YOUR ACCESS KEY SECRET--",
		},
	}

	done, _, _ := nuker.NukeItAll(currentAccount)
    <-done
}
```
