package unitTest

import (
	"pathinder/internal/parsers"
	"reflect"
	"testing"
)

func TestParseConnections(t *testing.T) {
	for _, tests := range []struct {
		name           string
		lines          []string
		ConnectionLine []string
		linenum        int
		wantError      bool
		errorString    string
		want           map[parsers.StationName]map[parsers.StationName]struct{}
	}{
		{
			name:           "Valid Connection",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"waterloo-victoria", "waterloo-euston", "st_pancras-euston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      false,
			want:           map[string]map[string]struct{}{"euston": {"st_pancras": {}, "waterloo": {}}, "st_pancras": {"euston": {}, "victoria": {}}, "victoria": {"st_pancras": {}, "waterloo": {}}, "waterloo": {"euston": {}, "victoria": {}}},
		},
		{
			name:           "Malformed Connection Format",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"waterloo-victoria", "waterloo-euston", "st_pancraseuston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      true,
			errorString:    "Malformed connection format in line #14: st_pancraseuston\n",
		},
		{
			name:           "Malformed Station Names",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"waterloo-victoria", "&waterloo-euston", "st_pancras-euston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      true,
			errorString:    "Malformed station names in line #13: &waterloo-euston\n",
		},
		{
			name:           "Non Existent Start Station",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"moon-victoria", "waterloo-euston", "st_pancras-euston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      true,
			errorString:    "Non-existent start station in line #12: moon\n",
		},
		{
			name:           "Non Existent End Station",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"waterloo-moon", "waterloo-euston", "st_pancras-euston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      true,
			errorString:    "Non-existent end station in line #12: moon\n",
		},
		{
			name:           "Same Start End Station Connection",
			lines:          []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			ConnectionLine: []string{"waterloo-waterloo", "waterloo-euston", "st_pancras-euston", "victoria-st_pancras"},
			linenum:        11,
			wantError:      true,
			errorString:    "Error on line #12. Start and end connections are the same: waterloo-waterloo",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			stationsMap, _ := parsers.ParseStations(tests.lines, tests.linenum)

			got, err := parsers.ParseConnections(tests.ConnectionLine, tests.linenum, stationsMap)
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

			if !reflect.DeepEqual(got, tests.want) {
				t.Fatalf("Test %v: got %v want %v", tests.name, got, tests.want)
			}
		})
	}
}
