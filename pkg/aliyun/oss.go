package aliyun

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

// OssService represents the OSS service
type OssService struct{}

// OssResource wraps OSS buckets
type OssResource struct {
	Account  account.Account
	Name     string
	Location string

	items []item
}

// Item is a single object in a bucket
type item struct {
	Account  account.Account
	Location string
	Key      string
}

// List returns a list of all buckets in an account
func (s OssService) List(account account.Account) ([]cloud.Resource, error) {
	client, err := oss.New("oss.aliyuncs.com", account.AccessKeyID, account.AccessKeySecret)

	if err != nil {
		return nil, err
	}

	bucketResult, err := client.ListBuckets()
	if err != nil {
		return nil, err
	}

	buckets := make([]cloud.Resource, len(bucketResult.Buckets))
	for i, bucket := range bucketResult.Buckets {
		b := OssResource{Name: bucket.Name, Location: bucket.Location, Account: account}
		items, err := listItemsInBucket(account, b)
		if err != nil {
			return nil, err
		}

		b.items = items
		buckets[i] = b
	}

	return buckets, nil
}

// String returns the name of the resource
func (r OssResource) String() string {
	return fmt.Sprintf("Bucket: %s (%d items deleted)", r.Name, len(r.items))
}

// Delete removes a bucket
func (r OssResource) Delete() (bool, error) {
	client, err := getOSSClient(r.Account, r.Location)
	if err != nil {
		return false, err
	}

	clientBucket, err := client.Bucket(r.Name)
	if err != nil {
		return false, err
	}

	for _, item := range r.items {
		err = clientBucket.DeleteObject(item.Key)
		if err != nil {
			return false, nil
		}
	}

	err = client.DeleteBucket(r.Name)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func listItemsInBucket(account account.Account, r OssResource) ([]item, error) {
	client, err := getOSSClient(account, r.Location)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(r.Name)
	if err != nil {
		return nil, err
	}

	itemResult, err := bucket.ListObjects()
	if err != nil {
		return nil, err
	}

	items := make([]item, len(itemResult.Objects))
	for i, object := range itemResult.Objects {
		items[i] = item{Account: account, Location: r.Location, Key: object.Key}
	}

	return items, nil
}

func getOSSClient(account account.Account, endpoint string) (*oss.Client, error) {
	return oss.New(fmt.Sprintf("%s.aliyuncs.com", endpoint), account.AccessKeyID, account.AccessKeySecret)
}
