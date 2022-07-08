package finder

import "fmt"

// FlightGraphItem is one item in the graph.
type FlightGraphItem struct {
	airport    string
	departures []*FlightGraphItem
	arrivals   []*FlightGraphItem
}

// AddDeparture adds the item to departures.
func (g *FlightGraphItem) AddDeparture(n *FlightGraphItem) {
	g.departures = append(g.departures, n)
}

// AddArrival adds the item to arrival.
func (g *FlightGraphItem) AddArrival(n *FlightGraphItem) {
	g.arrivals = append(g.arrivals, n)
}

// FlightGraph is the dest graph.
type FlightGraph struct {
	items map[string]*FlightGraphItem
}

// NewFlightGraph creates a new graph.
func NewFlightGraph(points [][]string) *FlightGraph {
	items := make(map[string]*FlightGraphItem)
	// fill empty points
	for i := range points {
		items[points[i][0]] = &FlightGraphItem{
			airport: points[i][0],
		}
		items[points[i][1]] = &FlightGraphItem{
			airport: points[i][1],
		}
	}
	// build graph relationships
	for i := range points {
		fromItem := items[points[i][0]]
		toItem := items[points[i][1]]
		fromItem.AddDeparture(toItem)
		toItem.AddArrival(fromItem)
	}
	return &FlightGraph{
		items: items,
	}
}

// FindPath returns the arrival and departure from all points in the graph or error.
func (g *FlightGraph) FindPath() ([]string, error) {
	departure, arrival, err := g.findPotentialDepartureAndArrival()
	if err != nil {
		return nil, err
	}

	if departure == "" && arrival == "" {
		return nil, fmt.Errorf("the path is cycled")
	}

	// in this block we might have either departure or arrival or both,
	// but we need at least something to find/validate the missing.
	if arrival != "" {
		fullPath, err := g.findFullPath(func(item *FlightGraphItem) []*FlightGraphItem {
			return item.arrivals
		}, arrival)
		if err != nil {
			return nil, err
		}

		departure = fullPath[len(fullPath)-1]

		return []string{departure, arrival}, nil
	}

	// or find arrival by departure
	fullPath, err := g.findFullPath(func(item *FlightGraphItem) []*FlightGraphItem {
		return item.departures
	}, departure)
	if err != nil {
		return nil, err
	}

	arrival = fullPath[len(fullPath)-1]

	return []string{departure, arrival}, nil
}

// findPotentialDepartureAndArrival searches from potential departure and arrival.
func (g *FlightGraph) findPotentialDepartureAndArrival() (string, string, error) {
	var departure, arrival string
	for airport := range g.items {
		item := g.items[airport]
		if len(item.arrivals) == 0 {
			if departure != "" {
				return "", "", fmt.Errorf("found two potential departures")
			}
			departure = airport
		}
		if len(item.departures) == 0 {
			if arrival != "" {
				return "", "", fmt.Errorf("found two potential arrivals")
			}
			arrival = airport
		}
	}
	return departure, arrival, nil
}

// findFullPath searches for the longest path in the graph, validate whether it covers all points
// and returns the full path, or error.
func (g *FlightGraph) findFullPath(direction func(*FlightGraphItem) []*FlightGraphItem, from string) ([]string, error) {
	roteLimits := make(map[string]int)
	expectedPathLength := 1
	for i := range g.items {
		items := direction(g.items[i])
		for j := range items {
			rotePoint := items[j].airport
			count := roteLimits[rotePoint]
			roteLimits[rotePoint] = count + 1
			expectedPathLength++
		}
	}

	path := g.searchLongestPath(direction, g.items[from], roteLimits, expectedPathLength, make([]string, 0))
	if len(path) == 0 || expectedPathLength != len(path) {
		return nil, fmt.Errorf("can't find the path from %s airport", from)
	}

	return path, nil
}

// searchLongestPath creates all possible paths from the graph and returns the longest.
func (g *FlightGraph) searchLongestPath(
	direction func(*FlightGraphItem) []*FlightGraphItem,
	item *FlightGraphItem,
	roteLimit map[string]int,
	expectedPathLength int,
	path []string,
) []string {
	path = append(path, item.airport)
	// optimisation to rerun the longest pass faster
	if len(path) == expectedPathLength {
		return path
	}
	maxPath := path
	items := direction(item)
	for i := range items {
		if roteLimit[items[i].airport] == 0 {
			continue
		}
		// copy to prevent mutations
		childRoteLimits := make(map[string]int)
		for j := range roteLimit {
			childRoteLimits[j] = roteLimit[j]
		}
		count := childRoteLimits[items[i].airport]
		childRoteLimits[items[i].airport] = count - 1
		// copy to prevent mutations
		childPath := make([]string, 0)
		for i := range path {
			childPath = append(childPath, path[i])
		}

		newPath := g.searchLongestPath(direction, items[i], childRoteLimits, expectedPathLength, childPath)
		if len(maxPath) < len(newPath) {
			maxPath = newPath
		}
	}

	return maxPath
}
