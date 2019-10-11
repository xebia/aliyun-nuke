package nuker

import (
	"reflect"
	"testing"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type dummyService struct{}

type dummyResource struct {
	Name  string
	Index int
}

var dummyResources []cloud.Resource

func TestNuke(t *testing.T) {
	dummyResources = []cloud.Resource{
		dummyResource{Name: "Resource 1"},
		dummyResource{Name: "Resource 2"},
		dummyResource{Name: "Resource 3"},
	}

	type args struct {
		currentAccount account.Account
		services       []cloud.Service
		regions        []account.Region
	}
	tests := []struct {
		name string
		args args
		want []cloud.Resource
	}{
		{
			name: "No-op returns empty list of deleted resources",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{}, regions: []account.Region{"nl-hilversum"}},
			want: []cloud.Resource{},
		},
		{
			name: "Returns deleted resources of a single region",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{dummyService{}}, regions: []account.Region{"nl-hilversum"}},
			want: []cloud.Resource{
				dummyResource{Name: "Resource 1"},
				dummyResource{Name: "Resource 2"},
				dummyResource{Name: "Resource 3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Nuke(tt.args.currentAccount, tt.args.services, tt.args.regions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Nuke() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (d dummyService) IsGlobal() bool {
	return false
}

func (d dummyService) List(account.Region, account.Account) ([]cloud.Resource, error) {
	return dummyResources, nil
}

func (d dummyService) String() string {
	return "dummyService"
}

func (d dummyResource) Delete(account.Region, account.Account) error {
	// Remove resource from dummy resource list
	dummyResources = dummyResources[1:]
	return nil
}

func (d dummyResource) String() string {
	return d.Name
}
