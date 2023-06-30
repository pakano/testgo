package main

import "fmt"

func IsSeller(name string) bool {
	return name == "jonny1"
}

func main() {
	queue := make([]string, 0)
	graph := make(map[string][]string)
	searched := make(map[string]bool)
	graph["you"] = []string{"bob", "alice", "claire"}
	graph["bob"] = []string{"you", "anuj", "peggy"}
	graph["alice"] = []string{"peggy"}
	graph["claire"] = []string{"thom", "jonny"}
	graph["anuj"] = nil
	graph["peggy"] = nil
	graph["thom"] = nil
	graph["jonny"] = nil
	queue = append(queue, "you")

	for len(queue) != 0 {
		fmt.Println(queue)
		name := queue[0]
		queue = queue[1:]

		if !searched[name] {
			continue
		}

		if IsSeller(name) {
			fmt.Println(name, "is seller")
			return
		}
		searched[name] = true
		neighbors := graph[name]
		queue = append(queue, neighbors...)
	}
}
