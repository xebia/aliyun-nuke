package cloud

import "github.com/xebia/aliyun-nuke/pkg/account"

// Service is a single service in a cloud provider
type Service interface {
	List(account.Account) ([]Resource, error)
}

// Resource is the single unit to be deleted with a single service
type Resource interface {
	Delete() (bool, error)
	String() string
}
