package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	// Metrics are registered here! Importing registers once
	"github.com/converged-computing/metrics-operator/pkg/metrics"
	_ "github.com/converged-computing/metrics-operator/pkg/metrics/app"
	//	_ "github.com/converged-computing/metrics-operator/pkg/metrics/io"
	//	_ "github.com/converged-computing/metrics-operator/pkg/metrics/network"
	//	_ "github.com/converged-computing/metrics-operator/pkg/metrics/perf"
	//
	// +kubebuilder:scaffold:imports
)

var (
	baseurl = "https://converged-computing.github.io/metrics-operator/getting_started/metrics.html"
)

type MetricOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Family      string `json:"family"`
	Type        string `json:"type"`
	Image       string `json:"image"`
	Url         string `json:"url"`
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Please provide a filename to write to")
	}
	filename := os.Args[1]
	records := []MetricOutput{}
	for _, metric := range metrics.Registry {
		newRecord := MetricOutput{
			Name:        metric.Name(),
			Description: metric.Description(),
			Family:      metric.Family(),
			Image:       metric.Image(),
			Url:         metric.Url(),
		}
		records = append(records, newRecord)
	}

	// Ensure we are consistent in ordering
	sort.Slice(records, func(i, j int) bool {
		return records[i].Name < records[j].Name
	})

	file, err := json.MarshalIndent(records, "", " ")
	if err != nil {
		log.Fatalf("Could not marshall records %s\n", err.Error())
	}
	err = os.WriteFile(filename, file, 0644)
	if err != nil {
		log.Fatalf("Could not write to file %s: %s\n", filename, err.Error())
	}
}
