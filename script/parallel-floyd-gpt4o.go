package main

import (
	"fmt"
	"math"
	"slices"
	"sync"
)

// INF represents infinity
const INF = math.MaxInt

// Initialize the graph with distances
func initializeGraph(V int) [][]int {
	dist := make([][]int, V)
	for i := 0; i < V; i++ {
		dist[i] = make([]int, V)
		for j := 0; j < V; j++ {
			if i == j {
				dist[i][j] = 0
			} else {
				dist[i][j] = INF
			}
		}
	}
	return dist
}

// ParallelFloydWarshall implements the Floyd-Warshall algorithm in parallel
func ParallelFloydWarshall(graph [][]int, V, numWorkers int) [][]int {
	dist := make([][]int, V)
	for i := range graph {
		dist[i] = make([]int, V)
		copy(dist[i], graph[i])
	}

	var wg sync.WaitGroup
	chunks := V / numWorkers

	for k := 0; k < V; k++ {
		wg.Add(numWorkers)
		for worker := 0; worker < numWorkers; worker++ {
			go func(worker int) {
				defer wg.Done()
				start := worker * chunks
				end := start + chunks
				if worker == numWorkers-1 {
					end = V
				}
				for i := start; i < end; i++ {
					for j := 0; j < V; j++ {
						if dist[i][k] != INF && dist[k][j] != INF && dist[i][j] > dist[i][k]+dist[k][j] {
							dist[i][j] = dist[i][k] + dist[k][j]
						}
					}
				}
			}(worker)
		}
		wg.Wait()
	}

	return dist
}

func main() {
	V := 4
	numWorkers := 4

	graph := initializeGraph(V)
	graph[0][1] = 3
	graph[0][2] = 1
	graph[1][2] = 7
	graph[1][3] = 5
	graph[2][3] = 2

	fmt.Println("Initial graph:")
	for _, row := range graph {
		fmt.Println(row)
	}

	dist := ParallelFloydWarshall(graph, V, numWorkers)

	fmt.Println("Distance matrix after executing Floyd-Warshall:")
	for _, row := range dist {
		fmt.Println(row)
	}

	// Below is written by me, not AI ( ﾟ∀。)

	dist1 := ParallelFloydWarshall(graph, V, 1)
	fmt.Println("Distance matrix after executing Floyd-Warshall (serial):")
	for _, row := range dist1 {
		fmt.Println(row)
	}

	eq := slices.EqualFunc(dist, dist1, func(row1 []int, row2 []int) bool {
		return slices.Equal(row1, row2)
	})
	fmt.Println(eq)
}
