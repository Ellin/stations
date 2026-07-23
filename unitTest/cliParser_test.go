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
			name: "Valid arguments",
			args: []string{
				"../network/london_network.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: false,
			want:    "waterloo st_pancras 3",
		},
		{
			name: "Insufficient arguments",
			args: []string{
				"../network/london_network.map",
				"waterloo",
			},
			wantErr: true,
			want:    "Error: Incorrect Length Of Arguments. To See The Usage run: \n\n\t\t'go run . -h'\n",
		}, {
			name: "Too many arguments",
			args: []string{
				"../network/london_network.map",
				"waterloo",
				"st_pancras",
				"3",
				"4",
			},
			wantErr: true,
			want:    "Error: Incorrect Length Of Arguments. To See The Usage run: \n\n\t\t'go run . -h'\n",
		},
		{
			name: "Non-existent map",
			args: []string{
				"../network/notfound.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: true,
			want:    "Error: File ../network/notfound.map Does Not Exist.\n",
		},
		{
			name: "Invalid file ext",
			args: []string{
				"../network/map2.txt",
				"waterloo",
				"st_pancras",
				"3",
			},
			wantErr: true,
			want:    "Error: Invalid File Extension '.txt', Allowed extentions = [.map].",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flag  with each test
			os.Args = append([]string{"program"}, tests.args...)

			got, err := parsers.ParseArgs()

			if tests.wantErr {
				if err == nil {
					t.Fatalf("want error, got nil")
				}
				if err.Error() != tests.want {
					t.Fatalf("Test %v: got %v, want %v", tests.name, err, tests.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("Test %s > unexpected error: %v", tests.name, err)
			}

			gotresult := fmt.Sprintf("%s %s %d", got.StartStation, got.EndStation, got.TrainCount)
			if gotresult != tests.want {
				t.Errorf("Test %v: got %v, want %v", tests.name, gotresult, tests.want)
			}
		})
	}

}
