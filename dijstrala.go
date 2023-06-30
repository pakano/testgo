package main

import (
	"errors"
	"fmt"
	"math"
)

func findLowestCostNode(costs map[string]int, processed map[string]bool) (lowestCostNode string, lowestCost int, err error) {
	lowestCost = math.MaxInt
	for k := range costs {
		if costs[k] < lowestCost && !processed[k] {
			lowestCost = costs[k]
			lowestCostNode = k
		}
	}
	if lowestCostNode == "" {
		err = errors.New("not found")
	}
	return
}

func main() {
	graph := make(map[string]map[string]int)
	costs := make(map[string]int)
	parents := make(map[string]string)
	processed := make(map[string]bool)

	graph["start"] = make(map[string]int)
	graph["start"]["A"] = 6
	graph["start"]["B"] = 2
	graph["A"] = make(map[string]int)
	graph["A"]["end"] = 1
	graph["B"] = make(map[string]int)
	graph["B"]["A"] = 3
	graph["B"]["end"] = 5

	costs["A"] = 6
	costs["B"] = 2
	costs["end"] = math.MaxInt

	parents["A"] = "start"
	parents["B"] = "start"
	parents["end"] = ""

	for {
		lowestCostNode, lowestCost, err := findLowestCostNode(costs, processed)
		if err != nil {
			fmt.Println(err)
			break
		}

		if lowestCostNode == "end" {
			break
		}

		neighbors := graph[lowestCostNode]
		for k := range neighbors {
			cost := lowestCost + neighbors[k]
			costold, ok := costs[k]
			if !ok || costold > cost {
				costs[k] = cost
				parents[k] = lowestCostNode
			}
		}
		processed[lowestCostNode] = true

		fmt.Println(costs)
		fmt.Println(parents)

	}
	fmt.Println(parents)

	nodes := []string{}
	node := "end"
	for {
		nodes = append(nodes, node)
		node = parents[node]
		if node == "" {
			break
		}
	}
	for i := len(nodes) - 1; i >= 0; i-- {
		fmt.Print(nodes[i], "->")
	}
	fmt.Println()
}
