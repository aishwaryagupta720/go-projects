package entities

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func CreateComparisonChart(loadDescription string, requests int) *charts.Bar {
	serverCount := 1000
	originalServers := createServers(serverCount)

	serversRR := make([]Server, serverCount)
	copy(serversRR, originalServers)
	roundRobin(serversRR, requests)

	serversP2C := make([]Server, serverCount)
	copy(serversP2C, originalServers)
	powerOfTwoChoices(serversP2C, requests)

	serversRandom := make([]Server, serverCount)
	copy(serversRandom, originalServers)
	randomAllocation(serversRandom, requests)

	serversLRC := make([]Server, serverCount)
	copy(serversLRC, originalServers)
	leastRecentlyContacted(serversLRC, requests)

	serversWeightedR := make([]Server, serverCount)
	copy(serversWeightedR, originalServers)
	weightedRandom(serversWeightedR, requests)

	serversWeightedRR := make([]Server, serverCount)
	copy(serversWeightedRR, originalServers)
	weightedRoundRobin(serversWeightedRR, requests)

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    fmt.Sprintf("Load Balancing Performance: %s Load", loadDescription),
			Subtitle: "Comparison of Round Robin, Power of Two Choices, Random, Least Recently Contacted, Weighted Random, and Weighted Round Robin",
			Left:     "center",
			Top:      "5%",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:       true,
			Right:      "10%",
			Top:        "top",
			Orient:     "horizontal",
			ItemWidth:  25,
			ItemHeight: 10,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "2200px",
			Height: "600px",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Servers",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Number of Requests",
		}),
	)

	serverIDs := make([]string, serverCount)
	for i := range serverIDs {
		serverIDs[i] = fmt.Sprintf("Server %d", i+1)
	}
	bar.SetXAxis(serverIDs).
		AddSeries("Round Robin", generateBarItems(serversRR)).
		AddSeries("Power of Two Choices", generateBarItems(serversP2C)).
		AddSeries("Random", generateBarItems(serversRandom)).
		AddSeries("Least Recently Contacted", generateBarItems(serversLRC)).
		AddSeries("Weighted Random", generateBarItems(serversWeightedR)).
		AddSeries("Weighted Round Robin", generateBarItems(serversWeightedRR))

	return bar
}

func generateBarItems(servers []Server) []opts.BarData {
	items := make([]opts.BarData, len(servers))
	for i, server := range servers {
		items[i] = opts.BarData{Value: server.Load}
	}
	return items
}
