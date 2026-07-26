package unitTest

import (
	"flag"
	"os"
	"pathinder/internal/parsers"
	"pathinder/internal/paths"
	"testing"
)

type schedulerTest struct {
	name        string
	args        []string
	wantErr     bool
	errorString string
	want        int
}

func TestOutputCheck(t *testing.T) {
	for _, tests := range []schedulerTest{
		{
			name: "London network, 3 trains, EdmondsKarp",
			args: []string{
				"../network/london_network.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			want: 3,
		},
		{
			name: "London network, 1 train, EdmondsKarp",
			args: []string{
				"../network/london_network.map",
				"waterloo",
				"st_pancras",
				"1",
			},
			want: 2,
		},
		{
			name: "London network, 3 trains, Dinic",
			args: []string{
				"-alg", "Dinic",
				"../network/london_network.map",
				"waterloo",
				"st_pancras",
				"3",
			},
			want: 3,
		},
		{
			name: "Single long path, 1 train",
			args: []string{
				"../network/onepath.map",
				"a",
				"z",
				"1",
			},
			want: 4,
		},
		{
			name: "Single long path, 2 trains",
			args: []string{
				"../network/onepath.map",
				"a",
				"z",
				"2",
			},
			want: 5,
		},
		{
			name: "onepath2, 1 train, EdmondsKarp",
			args: []string{
				"../network/onepath2.map",
				"b",
				"z",
				"1",
			},
			want: 3,
		},
		{
			name: "onepath2, 3 trains, Dinic",
			args: []string{
				"-alg", "Dinic",
				"../network/onepath2.map",
				"b",
				"z",
				"3",
			},
			want: 5,
		},
		{
			name: "map1, 2 trains, EdmondsKarp",
			args: []string{
				"../network/map1.map",
				"waterloo",
				"st_pancras",
				"2",
			},
			want: 3,
		},
		{
			name: "map2, 10 trains, EdmondsKarp",
			args: []string{
				"../network/map2.map",
				"tampere",
				"helsinki",
				"10",
			},
			want: 9,
		},
		{
			name: "map2, 10 trains, Dinic",
			args: []string{
				"-alg", "Dinic",
				"../network/map2.map",
				"tampere",
				"helsinki",
				"10",
			},
			want: 9,
		},
		{
			name: "map2, single track, 20 trains, EdmondsKarp",
			args: []string{
				"../network/map2.map",
				"hanko",
				"kokkola",
				"20",
			},
			want: 25,
		},
		{
			name: "No path between start and end",
			args: []string{
				"../network/nopath.map",
				"a",
				"z",
				"1",
			},
			wantErr:     true,
			errorString: "Error: No Path From Start To End Stations",
		},
	} {
		t.Run(tests.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // reset flag  with each test
			os.Args = append([]string{"program"}, tests.args...)

			got, schedulerErr := runSchulerTest(t, tests)
			if tests.wantErr {
				if schedulerErr == nil {
					t.Fatalf("want error, got nil")
				}
				if schedulerErr.Error() != tests.errorString {
					t.Fatalf("Test %v: got %v, want %v", tests.name, schedulerErr, tests.want)
				}
				return
			}

			if schedulerErr != nil {
				t.Fatalf("Test %s > unexpected ReadFile error: %v", tests.name, schedulerErr)

			}
			if got != tests.want {
				t.Errorf("Test %v: got %v, want %v", tests.name, got, tests.want)
			}
		})
	}

}

func runSchulerTest(t *testing.T, tests schedulerTest) (int, error) {

	arg, argErr := parsers.ParseArgs()
	if argErr != nil {
		t.Fatalf("Test %s > unexpected error: %v", tests.name, argErr)
	}

	fileText, err := os.ReadFile(arg.MapFile)
	if err != nil {
		t.Fatalf("Test %s > unexpected ReadFile error: %v", tests.name, err)
	}

	networkData, err := parsers.ParseNetworkMap(string(fileText))
	if err != nil {
		t.Fatalf("Test %s > unexpected ParseNetworkMap error: %v", tests.name, err)
	}

	if err := parsers.ValidateStartAndEndStation(&networkData, arg.StartStation, arg.EndStation); err != nil {
		t.Fatalf("Test %s > unexpected ValidateStartAndEndStation error: %v", tests.name, err)
	}

	networkData.Start, networkData.End = arg.StartStation, arg.EndStation

	graph := paths.Graph{}
	graph.CreateVertexMaps(networkData)
	graph.CreateEKGraph(networkData)

	got, schedulerErr := graph.RunScheduler(arg.StartStation, arg.EndStation, arg.Algo, arg.TrainCount)

	return got, schedulerErr
}
