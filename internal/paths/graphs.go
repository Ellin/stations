package paths

import (
	"fmt"
	"pathfinder/model"
)

type StationName = string

type Graph struct {
	EKGraph       map[int][]Edge
	VertexIDMap   map[int]StationName
	VertexNameMap map[StationName]int
	LeveL         []int // use this array to make the level graph is seprates the nodes into tree like form
	// it stores the distance bw current node and source
	EdgeBook []int // keep record of the current neighbor edge so we dont keep checking the same edge each loop
	DeadEnd  []bool
}

type Edge struct {
	To   int // node index
	Cap  int
	Rev  int
	Real bool // indicates whether an edge is real (forward) edge or false (reverse) edge
}

func (g *Graph) CreateVertexMaps(nd model.NetworkData) {
	// Split every vertex except start and sink into V_in, V_out.
	// While splitting vertexes, create a vertexID map mapping numerical ids with station names
	g.VertexIDMap = make(map[int]StationName)
	g.VertexNameMap = make(map[StationName]int) // maps station names with V_in ID. V_out ID == V_in ID + 1.

	var i int
	for station := range nd.NetworkMap {
		g.VertexIDMap[i] = station
		g.VertexNameMap[station] = i

		// split all vertexes that are not start and end
		if station != nd.Start && station != nd.End {
			g.VertexIDMap[i+1] = station // V_out
			i += 2
		} else {
			i++
		}
	}
	g.LeveL = make([]int, len(g.VertexIDMap))
	g.DeadEnd = make([]bool, len(g.VertexIDMap))
	g.EdgeBook = make([]int, len(g.VertexIDMap))
}

// Given edge A->B, create an edge from A->B and also a corresponding reverse edge B->A with cap: 0
func (g *Graph) AddEdge(fromID, toID int) {
	forward := Edge{
		To:   toID,
		Cap:  1,                    // Since all tracks and stations have max cap 1
		Rev:  len(g.EKGraph[toID]), // the array index of the reverse edge within a node's adjacency list (NOT a vertex ID)
		Real: true,
	}

	// Since the reverse edge is always added to the adjacency list of the opposite node and added at the same time as creation,
	// we know the index of the newly added edge will be the length of the list (i.e. the index of the last element + 1)

	reverse := Edge{
		To:   fromID,
		Cap:  0, // Since the reverse edges are not real edges, capacity is 0
		Rev:  len(g.EKGraph[fromID]),
		Real: false,
	}

	// Add the forward and reverse edges to the corresponding adjacency lists
	g.EKGraph[fromID] = append(g.EKGraph[fromID], forward)
	g.EKGraph[toID] = append(g.EKGraph[toID], reverse)
}

func (g *Graph) CreateEKGraph(nd model.NetworkData) {
	g.EKGraph = make(map[int][]Edge)

	for station, adjNodes := range nd.NetworkMap {
		if station != nd.Start && station != nd.End {
			inID := g.VertexNameMap[station]
			outID := inID + 1

			g.AddEdge(inID, outID) // Create edge between V_in and V_out halves of split vertex

			for adjNode := range adjNodes {
				g.AddEdge(outID, g.VertexNameMap[adjNode])
			}
		} else {
			for adjNode := range adjNodes {
				g.AddEdge(g.VertexNameMap[station], g.VertexNameMap[adjNode])
			}
		}
	}
}

func (g *Graph) PrintGraph(nd model.NetworkData) {
	fmt.Printf("\nTRANSFORMED EKGRAPH \n")

	for key, list := range g.EKGraph {

		station := g.VertexIDMap[key]
		if station != nd.Start && station != nd.End {
			if nextStation, ok := g.VertexIDMap[key+1]; ok && nextStation == station {
				fmt.Printf("%v_in --- ", station)
			} else {
				fmt.Printf("%v_out --- ", station)
			}

		} else {
			fmt.Printf("%v --- Edges: ", station)
		}

		for _, edge := range list {
			fmt.Printf("%v (cap: %v), ", g.VertexIDMap[edge.To], edge.Cap)
		}
		fmt.Println()
	}
}
