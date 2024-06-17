package entities

import (
	"math/rand"
	"time"
)

// roundRobin distributes requests sequentially to each server in a circular manner.
func roundRobin(servers []Server, requests int) {
	serverCount := len(servers)
	for i := 0; i < requests; i++ {
		// Increment the load of the current server
		servers[i%serverCount].Load++
		// Update the last contact time of the current server
		servers[i%serverCount].LastContactTime = time.Now()
	}
}

// powerOfTwoChoices distributes requests to the server with the least load among two randomly chosen servers.
func powerOfTwoChoices(servers []Server, requests int) {
	serverCount := len(servers)
	for i := 0; i < requests; i++ {
		// Randomly select two servers
		choice1 := rand.Intn(serverCount)
		choice2 := rand.Intn(serverCount)
		// Assign the request to the server with the least load
		if servers[choice1].Load < servers[choice2].Load {
			servers[choice1].Load++
			servers[choice1].LastContactTime = time.Now()
		} else {
			servers[choice2].Load++
			servers[choice2].LastContactTime = time.Now()
		}
	}
}

// randomAllocation distributes requests to servers chosen at random.
func randomAllocation(servers []Server, requests int) {
	serverCount := len(servers)
	for i := 0; i < requests; i++ {
		// Randomly select a server
		randomServer := rand.Intn(serverCount)
		// Increment the load of the selected server
		servers[randomServer].Load++
		// Update the last contact time of the selected server
		servers[randomServer].LastContactTime = time.Now()
	}
}

// leastRecentlyContacted distributes requests to the server that has not been contacted for the longest time.
func leastRecentlyContacted(servers []Server, requests int) {
	for i := 0; i < requests; i++ {
		// Find the least recently contacted server
		lrcIndex := 0
		for j, server := range servers {
			if server.LastContactTime.Before(servers[lrcIndex].LastContactTime) {
				lrcIndex = j
			}
		}
		// Increment the load of the least recently contacted server
		servers[lrcIndex].Load++
		// Update the last contact time of the least recently contacted server
		servers[lrcIndex].LastContactTime = time.Now()
	}
}

// weightedRoundRobin distributes requests to servers based on their weight in a round-robin manner.
func weightedRoundRobin(servers []Server, requests int) {
	// Calculate the total weight of all servers
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	requestCount := 0
	for requestCount < requests {
		for i := 0; i < len(servers) && requestCount < requests; i++ {
			// Distribute requests according to the server's weight
			for j := 0; j < servers[i].Weight && requestCount < requests; j++ {
				servers[i].Load++
				servers[i].LastContactTime = time.Now()
				requestCount++
			}
		}
	}
}

// weightedRandom distributes requests to servers randomly, with the probability proportional to their weight.
func weightedRandom(servers []Server, requests int) {
	// Calculate the total weight of all servers
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	for i := 0; i < requests; i++ {
		// Select a server based on its weight
		r := rand.Intn(totalWeight)
		sum := 0
		for j, server := range servers {
			sum += server.Weight
			if r < sum {
				servers[j].Load++
				servers[j].LastContactTime = time.Now()
				break
			}
		}
	}
}
