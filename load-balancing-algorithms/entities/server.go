package entities

import (
	"math/rand"
	"time"
)

type Server struct {
	ID              int
	Load            int
	LastContactTime time.Time
	Weight          int
	Requests        int
}

// createServers initializes a slice of Server structs with the specified count.
// Each server is assigned a unique ID, a load of 0, a random last contact time within the last 1000 minutes, and a random weight between 1 and 100.
func createServers(count int) []Server {
	servers := make([]Server, count) // Initialize a slice of Server structs with the given count
	for i := range servers {
		servers[i] = Server{
			ID:              i + 1,                                                         // Assign a unique ID to each server, starting from 1
			Load:            0,                                                             // Initialize the load of each server is random , to demonstrate how algorithms behave with respect to the existing load
			LastContactTime: time.Now().Add(-time.Duration(rand.Intn(1000)) * time.Minute), // Assign a random last contact time within the last 1000 minutes to demonstrate how algorithms behave with respect to varying Last Contact Time
			Weight:          rand.Intn(100) + 1,                                            // Assign a random weight between 1 and 100 to demonstrate behavious or Weighted Random anf Weighted Round Robin with respect to the Unweighted ones
		}
	}
	return servers // Return the initialized slice of servers
}
