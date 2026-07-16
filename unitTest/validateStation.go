package unitTest

import (
	"pathinder/internal/parsers"
	"testing"
)

func TestparseStations(t *testing.T) {
	for _, tests := range []struct {
		name      string
		lines     []string
		linenum   int
		wantError bool
		want      string
	}{
		{
			name:      "valid",
			lines:     []string{"# south stations", "waterloo  , 3 , 1", "victoria,6,7", "# north stations", "euston,11,23", "st_pancras,5,15", "# international"},
			linenum:   1,
			wantError: false,
			want:      "map[euston:{euston 11 23} st_pancras:{st_pancras 5 15} victoria:{victoria 6 7} waterloo:{waterloo 3 1}]",
		},
	} {
		t.Run("", func(t *testing.T) {
			_, err := parsers.ParseStations(tests.lines, tests.linenum)

			if tests.wantError {
				if err == nil {
					t.Errorf("want error, got nil")
				}
				if tests.want != err.Error() {
					t.Errorf("Test %v: got %v, want %v", tests.name, err, tests.want)
				}
			}

			if err != nil {
				t.Fatalf("Test %s > unexpected error: %v", tests.name, err)
			}

			// if got != tests.want {
			// 	t.Fatalf("Test %v: got %v want %v", tests.name, got, tests.want)
			// }
		})
	}
}
