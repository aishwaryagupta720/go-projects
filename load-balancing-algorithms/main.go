package main

import (
	"fmt"
	"os"

	"load-balancing-algorithms/entities"

	"github.com/go-echarts/go-echarts/v2/components"
)

func main() {
	page := components.NewPage()
	page.AddCharts(
		entities.CreateComparisonChart("Low", 1000),
		entities.CreateComparisonChart("Medium", 5000),
		entities.CreateComparisonChart("High", 100000),
	)

	// Render the page
	f, err := os.Create("load_balancing_comparison.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	if err := page.Render(f); err != nil {
		fmt.Println("Error rendering page:", err)
		return
	}
}
