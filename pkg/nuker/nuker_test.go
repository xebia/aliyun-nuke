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
		want []NukeResult
	}{
		{
			name: "No-op returns empty list of deleted resources",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{}, regions: []account.Region{"region-0"}},
			want: []NukeResult{},
		},
		{
			name: "Returns deleted resources of a single region",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{dummyService{}}, regions: []account.Region{"region-1"}},
			want: []NukeResult{
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 1"},
					Error:    nil,
				},
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 2"},
					Error:    nil,
				},
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 3"},
					Error:    nil,
				},
			},
		},
		{
			name: "Returns deleted resources of multiple regions",
			args: args{currentAccount: account.Account{}, services: []cloud.Service{dummyService{}}, regions: []account.Region{"region-2", "region-3"}},
			want: []NukeResult{
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 1"},
					Error:    nil,
				},
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 2"},
					Error:    nil,
				},
				{
					Success:  true,
					Resource: dummyResource{Name: "Resource 3"},
					Error:    nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := Nuke(tt.args.currentAccount, tt.args.services, tt.args.regions)
			expected := tt.want
			got := make([]NukeResult, 0)
			for result := range results {
				got = append(got, result)
			}
			for _, gotResource := range got {
				if !contains(expected, gotResource) {
					t.Errorf("Nuke() returned unexpected result %v, wanted: %v", gotResource, tt.want)
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

func (d dummyResource) Id() string {
	return d.Name
}

func (d dummyResource) Type() string {
	return "dummyResource"
}

func contains(results []NukeResult, result NukeResult) bool {
	for _, item := range results {
		if item == result {
			return true
		}
	}
	return false
}
