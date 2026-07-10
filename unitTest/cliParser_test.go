package unitTest

import (
	"flag"
	"fmt"
	"os"
	"pathinder/internal/parsers"
	"testing"
)

func TestParseArgs(t *testing.T) {
	for _, tests := range []struct {
		name    string
		args    []string
		wantErr bool
		want    string
	}{
		{
			name: "valid",
			args: []string{
				"network/m5.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: false,
			want:    "waterloo,st_pancras,3",
		},
		{
			name: "Insufficient argument",
			args: []string{
				"network/m5.map",
				"waterloo",
			},
			wantErr: true,
			want:    "Error: Insufficient number of arguments.\n",
		},
		{
			name: "notfound map",
			args: []string{
				"network/notfound.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: true,
			want:    "Error: File network/notfound.map Does Not Exist.\n",
		},
		{
			name: "invalid file ext",
			args: []string{
				"network/map2.txt",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: true,
			want:    "Error: Invalid File Extension '.txt', Allowed extentions = [.map].",
		},
	} {
		t.Run("", func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flag  with each test
			os.Args = append([]string{"program"}, tests.args...)

			got, err := parsers.ParseArgs()

			if tests.wantErr {
				if err == nil {
					t.Errorf("want error, got nil")
				}
				if err.Error() != tests.want {
					t.Errorf("error = %v, want %v", err, tests.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("test %s unexpected error: %v", tests.name, err)
			}

			gotresult := fmt.Sprintf("%s%s%d", got.StartStation, got.EndStation, got.TrainCount)
			if gotresult != tests.want {
				t.Errorf("got %v, want %v", gotresult, tests.want)
			}
		})
	}

}
