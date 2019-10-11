package nuker

import (
	"testing"

	"github.com/xebia/aliyun-nuke/pkg/account"
	"github.com/xebia/aliyun-nuke/pkg/cloud"
)

type dummyService struct{}

type dummyResource struct {
	Name  string
	Index int
}

var dummyResources map[string][]cloud.Resource

func TestNuke(t *testing.T) {
	dummyResources = map[string][]cloud.Resource{
		"region-1": {
			dummyResource{Name: "Resource 1"},
			dummyResource{Name: "Resource 2"},
			dummyResource{Name: "Resource 3"},
		},
		"region-2": {
			dummyResource{Name: "Resource 1"},
		},
		"region-3": {
			dummyResource{Name: "Resource 2"},
			dummyResource{Name: "Resource 3"},
		},
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
			args: args{currentAccount: account.Account{}, services: []cloud.Service{}, regions: []account.Region{"region-0"}},
			want: []cloud.Resource{},
		},
		{
			name: "Returns deleted resources of a single region",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{dummyService{}}, regions: []account.Region{"region-1"}},
			want: []cloud.Resource{
				dummyResource{Name: "Resource 1"},
				dummyResource{Name: "Resource 2"},
				dummyResource{Name: "Resource 3"},
			},
		},
		{
			name: "Returns deleted resources of multiple regions",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{dummyService{}}, regions: []account.Region{"region-2", "region-3"}},
			want: []cloud.Resource{
				dummyResource{Name: "Resource 1"},
				dummyResource{Name: "Resource 2"},
				dummyResource{Name: "Resource 3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Nuke(tt.args.currentAccount, tt.args.services, tt.args.regions)
			expected := tt.want
			if len(got) != len(expected) {
				t.Errorf("Nuke() results were of incorrect length (got %v, wanted %v)", got, tt.want)
			}

			for _, gotResource := range got {
				if !contains(expected, gotResource) {
					t.Errorf("Nuke() results did not contain %v, wanted: %v", gotResource, tt.want)
				}
			}
		})
	}
}

func (d dummyService) IsGlobal() bool {
	return false
}

func (d dummyService) List(region account.Region, account account.Account) ([]cloud.Resource, error) {
	return dummyResources[string(region)], nil
}

func (d dummyService) String() string {
	return "dummyService"
}

func (d dummyResource) Delete(region account.Region, account account.Account) error {
	// Remove resource from dummy resource list
	dummyResources[string(region)] = dummyResources[string(region)][1:]
	return nil
}

func (d dummyResource) String() string {
	return d.Name
}

func contains(resources []cloud.Resource, resource cloud.Resource) bool {
	for _, item := range resources {
		if item == resource {
			return true
		}
	}
	return false
}
