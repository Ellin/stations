package unitTest

import (
	"pathinder/internal/parsers"
	"pathinder/model"
	"reflect"
	"testing"
)

func TestParseNetworkMap(t *testing.T) {
	londonMapText := `stations:
# south stations
waterloo  , 3 , 1
victoria,6,7

# north stations
euston,11,23
st_pancras,5,15 # international


connections:
waterloo-victoria
waterloo-euston
st_pancras-euston
victoria-st_pancras`

	for _, tests := range []struct {
		name           string
		mapText        string
		wantError      bool
		wantStationMap map[model.StationName]model.Station
		wantNetworkMap map[model.StationName]map[model.StationName]struct{}
	}{
		{
			name:      "Valid network map",
			mapText:   londonMapText,
			wantError: false,
			wantStationMap: map[model.StationName]model.Station{
				"waterloo": model.Station{
					Name: "waterloo",
					X:    3,
					Y:    1,
				},
				"victoria": model.Station{
					Name: "victoria",
					X:    6,
					Y:    7,
				},
				"euston": model.Station{
					Name: "euston",
					X:    11,
					Y:    23,
				},
				"st_pancras": model.Station{
					Name: "st_pancras",
					X:    5,
					Y:    15,
				},
			},
			wantNetworkMap: map[model.StationName]map[model.StationName]struct{}{
				"waterloo": {
					"victoria": {},
					"euston":   {},
				},
				"victoria": {
					"waterloo":   {},
					"st_pancras": {},
				},
				"euston": {
					"waterloo":   {},
					"st_pancras": {},
				},
				"st_pancras": {
					"euston":   {},
					"victoria": {},
				},
			},
		},
		{
			name: "No stations section header",
			mapText: `# south stations
waterloo  , 3 , 1
victoria,6,7

connections:
waterloo-victoria`,
			wantError:      true,
			wantStationMap: nil,
			wantNetworkMap: nil,
		},
		{
			name: "No connections section header",
			mapText: `stations:
# south stations
waterloo  , 3 , 1
victoria,6,7

victoria-waterloo`,
			wantError:      true,
			wantStationMap: nil,
			wantNetworkMap: nil,
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			got, err := parsers.ParseNetworkMap(tests.mapText)

			if tests.wantError {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Test %s > unexpected error: %v", tests.name, err)
			}

			if !reflect.DeepEqual(got.StationMap, tests.wantStationMap) {
				t.Errorf("Test %v: got %v want %v", tests.name, got, tests.wantStationMap)
			}

			if !reflect.DeepEqual(got.NetworkMap, tests.wantNetworkMap) {
				t.Errorf("Test %v: got %v want %v", tests.name, got, tests.wantNetworkMap)
			}
		})
	}
}
