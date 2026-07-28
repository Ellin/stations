package unitTest

import (
	"maps"
	"pathfinder/internal/parsers"
	"pathfinder/model"
	"testing"
)

func TestParseStations(t *testing.T) {

	for _, tests := range []struct {
		name        string
		lines       []string
		linenum     int
		wantError   bool
		errorString string
		want        map[string]model.Station
	}{
		{
			name:      "Valid",
			lines:     []string{"waterloo  , 3 , 1", "victoria,6,7", "euston,11,23", "st_pancras,5,15"},
			linenum:   1,
			wantError: false,
			want:      map[string]model.Station{"euston": {Name: "euston", X: 11, Y: 23}, "st_pancras": {Name: "st_pancras", X: 5, Y: 15}, "victoria": {Name: "victoria", X: 6, Y: 7}, "waterloo": {Name: "waterloo", X: 3, Y: 1}},
		},
		{

			name:      "Valid",
			lines:     []string{"jungle,1,3", "green_belt,3,4", "village,4,3", "mountain,5,2", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:   1,
			wantError: false,
			want:      map[string]model.Station{"clouds": {Name: "clouds", X: 6, Y: 5}, "desert": {Name: "desert", X: 9, Y: 9}, "downtown": {Name: "downtown", X: 5, Y: 4}, "farms": {Name: "farms", X: 3, Y: 5}, "grasslands": {Name: "grasslands", X: 3, Y: 6}, "green_belt": {Name: "green_belt", X: 3, Y: 4}, "industrial": {Name: "industrial", X: 8, Y: 6}, "jungle": {Name: "jungle", X: 1, Y: 3}, "metropolis": {Name: "metropolis", X: 7, Y: 8}, "mountain": {Name: "mountain", X: 5, Y: 2}, "suburbs": {Name: "suburbs", X: 4, Y: 6}, "treetop": {Name: "treetop", X: 6, Y: 4}, "village": {Name: "village", X: 4, Y: 3}, "wetlands": {Name: "wetlands", X: 7, Y: 4}},
		},
		{

			name:        "Malformed Station Data",
			lines:       []string{"jungle,1,3", "green_belt,3,4", "village,4,3", "mountain,5", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Malformed station format on line #5: mountain,5\n",
		},
		{
			name:        "Malformed Station Name",
			lines:       []string{"jungle,1,3", "&green_belt,3,4", "village,4,3", "mountain,5,6", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Invalid station name in line #3: &green_belt\n",
		},
		{
			name:        "Malformed Station Coordinates: Not an integer",
			lines:       []string{"jungle,1,3", "green_belt,2.5,4", "village,4,3", "mountain,5,6", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Invalid coordinates at station green_belt in line #3: 2.5, 4\n",
		},
		{
			name:        "Malformed Station Coordinates: Negative integer",
			lines:       []string{"jungle,1,3", "green_belt,-1,4", "village,4,3", "mountain,5,6", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Invalid coordinates at station green_belt in line #3: -1, 4\n",
		},
		{
			name:        "Duplicate Station Coordinates",
			lines:       []string{"jungle,1,3", "green_belt,3,4", "village,3,4", "mountain,5,6", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Duplicate coordinates in line #4: 3,4\n",
		},
		{
			name:        "Duplicate Station Name",
			lines:       []string{"jungle,1,3", "green_belt,4,3", "jungle,3,4", "mountain,5,6", "treetop,6,4", "grasslands,3,6", "suburbs,4,6", "clouds,6,5", "wetlands,7,4", "farms,3,5", "downtown,5,4", "metropolis,7,8", "industrial,8,6", "desert,9,9"},
			linenum:     1,
			wantError:   true,
			errorString: "Duplicate station name jungle in line #4\n",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			got, err := parsers.ParseStations(tests.lines, tests.linenum)

			if tests.wantError {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				if tests.errorString != err.Error() {
					t.Fatalf("got %v, want %v", err, tests.errorString)
				}
				return
			}

			if err != nil {
				t.Fatalf("Test %s > unexpected error: %v", tests.name, err)
			}

			if !maps.Equal(got, tests.want) {
				t.Fatalf("Test %v: got %v want %v", tests.name, got, tests.want)
			}
		})
	}
}
