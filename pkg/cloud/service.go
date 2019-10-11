package cloud

import "github.com/xebia/aliyun-nuke/pkg/account"

// Service is a single service in a cloud provider
type Service interface {
	List(account.Region, account.Account) ([]Resource, error)
	IsGlobal() bool
}

// Resource is the single unit to be deleted within a single service
type Resource interface {
	Delete(account.Region, account.Account) error
	String() string
}
