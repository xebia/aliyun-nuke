# aliyun-nuke

Clears out all resources in a given Alibaba Cloud account. **Use with extreme caution!**

## Deletion process

Because Alibaba Cloud has many regions and every region can contain different resources, `aliyun-nuke` loops over all services and all regions to bruteforce the deletion of the resources. 
If a resource cannot be deleted because a dependent resource still exists, the dependent resource is deleted first and the process is started again. This repeats until no more resources 
were deleted for a single run.

## Supported services

Currently supported services are:

| Service            | Elements                                                |
| ------------------ | ------------------------------------------------------- |
| ECS                | Disks, instances, key pairs, security groups, snapshots |
| ESS (Auto Scaling) | Scaling groups, scheduled tasks                         |
| OSS                | Buckets, objects                                        |
| RAM                | Users, groups, roles, policies                          |
| RDS                | Instances                                               |
| SLB                | Load balancers                                          |
| VPC                | VPCs, VSwitches, NAT gateways, EIPs                     |

Any other resources will be kept as-is. If any unsupported resources block the deletion of the above resource types, `aliyun-nuke` will stop the deletion process and quit
after 60 seconds of retrying.

## Getting started as CLI tool

Build aliyun-cli from the source code by running `go build .` in the root directory of the repository. This assumes your Go environment is set up correctly.

When you have built an executable, define the following environment variables:

```bash
export ALIYUN_NUKE_ACCESS_KEY_ID=--YOUR ACCESS KEY ID--
export ALIYUN_NUKE_ACCESS_KEY_SECRET=--YOUR ACCESS KEY SECRET--
```

You can find these on the [AccessKey](https://ak-console.aliyun.com/) page in the Alibaba Cloud console. If you haven't created access keys before you might have to create them.

Then run `aliyun-nuke` to clear the account:

```bash
./aliyun-nuke destroy [--regions <region1,region2,...>]
```

If you want to see the help for the CLI or a subcommand, use the `help` subcommand:

```bash
./aliyun-nuke help
./aliyun-nuke help destroy
```

## Library usage

`aliyun-nuke` can also be used as a library:

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
    
    excludedIds := make([]string, 0) // fill this with resource IDs that you don't want to have deleted
    force := false // set to true to also delete service resources, like RAM roles
	results := nuker.NukeItAll(currentAccount, account.AllRegions, excludedIds, force)
    for range results {
        // process results (always consume this channel!)	
    }
}
```
