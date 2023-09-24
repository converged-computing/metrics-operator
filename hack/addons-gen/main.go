package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	"github.com/converged-computing/metrics-operator/pkg/addons"
	// Metrics are registered here! Importing registers once
	//
	// +kubebuilder:scaffold:imports
)

var (
	baseurl = "https://converged-computing.github.io/metrics-operator/getting_started/addons.html"
)

type AddonOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Family      string `json:"family"`
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Please provide a filename to write to")
	}
	filename := os.Args[1]
	records := []AddonOutput{}
	for _, addon := range addons.Registry {
		newRecord := AddonOutput{
			Name:        addon.Name(),
			Description: addon.Description(),
			Family:      addon.Family(),
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
