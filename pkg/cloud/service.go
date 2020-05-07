package cloud

import (
	"github.com/xebia/aliyun-nuke/pkg/account"
)

// Service is a single service in a cloud provider
type Service interface {
	List(account.Region, account.Account, bool) ([]Resource, error)
	IsGlobal() bool
}

// Resource is the single unit to be deleted within a single service
type Resource interface {
	Delete(account.Region, account.Account) error
	Id() string
	Type() string
}

var Services = make([]Service, 0)

// RegisterService is called by init() functions of cloud services to register themselves
func RegisterService(service Service) {
	Services = append(Services, service)
}
