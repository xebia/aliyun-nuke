package aliyun

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type OssBuckets struct{}

const START string = ""

type OssBucket struct {
	Name     string
	Location string

	items []item
}

type item struct {
	Key string
}

func init() {
	cloud.RegisterService(OssBuckets{})
}

func (s OssBuckets) IsGlobal() bool {
	return true
}

func (s OssBuckets) List(region account.Region, account account.Account, force bool) ([]cloud.Resource, error) {
	client, err := getOSSClient(account, "oss")

	if err != nil {
		return nil, err
	}

	bucketResult, err := client.ListBuckets()
	if err != nil {
		return nil, err
	}

	buckets := make([]cloud.Resource, len(bucketResult.Buckets))
	for i, bucket := range bucketResult.Buckets {
		b := OssBucket{Name: bucket.Name, Location: bucket.Location}
		items, err := listItemsInBucket(account, b)
		if err != nil {
			return nil, err
		}

		b.items = items
		buckets[i] = b
	}

	return buckets, nil
}

func (r OssBucket) Id() string {
	return fmt.Sprintf("%s (%d items)", r.Name, len(r.items))
}

func (r OssBucket) Type() string {
	return "OSS bucket"
}

func (r OssBucket) Delete(region account.Region, account account.Account) error {
	client, err := getOSSClient(account, r.Location)
	if err != nil {
		return err
	}

	clientBucket, err := client.Bucket(r.Name)
	if err != nil {
		return err
	}

	for _, item := range r.items {
		err = clientBucket.DeleteObject(item.Key)
		if err != nil {
			return nil
		}
	}

	err = client.DeleteBucket(r.Name)
	if err != nil {
		return nil
	}

	return nil
}

func listItemsInBucket(account account.Account, r OssBucket) ([]item, error) {
	client, err := getOSSClient(account, r.Location)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(r.Name)
	if err != nil {
		return nil, err
	}

	itemResults, err := getItemsFromMarker(bucket, START)
	items := make([]item, len(itemResults))
	for i, object := range itemResults {
		items[i] = item{Key: object.Key}
	}

	return items, nil
}

func getItemsFromMarker(bucket *oss.Bucket, start string) ([]oss.ObjectProperties, error) {
	results := make([]oss.ObjectProperties, 0)
	options := make([]oss.Option, 0)
	if len(start) > 0 {
		options = append(options, oss.Marker(start))
	}

	itemResult, err := bucket.ListObjects(options...)
	if err != nil {
		return results, err
	}

	for _, item := range itemResult.Objects {
		results = append(results, item)
	}

	if len(itemResult.NextMarker) > 0 {
		moreResults, err := getItemsFromMarker(bucket, itemResult.NextMarker)
		if err != nil {
			return results, err
		}
		results = append(results, moreResults...)
	}

	return results, nil
}

func getOSSClient(account account.Account, endpoint string) (*oss.Client, error) {
	return oss.New(fmt.Sprintf("%s.aliyuncs.com", endpoint), account.AccessKeyID, account.AccessKeySecret)
}
