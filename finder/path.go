// Package finder contains finders implementation.
package finder

import (
	"fmt"
	"strings"
)

// FlightPath is the finder which searches for the path based on the input flights.
type FlightPath struct{}

// NewFlightPath create a new instance of FlightPath.
func NewFlightPath() *FlightPath {
	return &FlightPath{}
}

// Find searches for the path from flights in FlightPath and returns it.
// Or error in case the path not found.
func (f *FlightPath) Find(flights [][]string) ([]string, error) {
	if len(flights) == 0 {
		return nil, fmt.Errorf("the 'flights' can't be empty")
	}
	for i := range flights {
		if len(flights[i]) != 2 { // nolint:gomnd // 2 is the length of a pair
			return nil, fmt.Errorf("the 'flight' %v is incorrect, expect 2 dest in the array", flights[i])
		}
		for j := range flights[i] {
			if strings.TrimSpace(flights[i][j]) == "" {
				return nil, fmt.Errorf("the 'flight' %v is incorrect, one item is empty", flights[i])
			}
		}
		if flights[i][0] == flights[i][1] {
			return nil, fmt.Errorf("the 'flight' %v is incorrect, departure equals arrival", flights[i])
		}
	}

	return NewFlightGraph(flights).FindPath()
}
