/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"encoding/json"
	"fmt"
	"log"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/utils"
	"go.uber.org/zap"
)

// Consistent logging identifiers that should be echoed to have newline after
var (
	logger *zap.SugaredLogger
)

// Default metadata (in JSON) to also put at the top of logs for parsing
// I'd like to improve upon this manual approach, it's a bit messy.
func Metadata(set *api.MetricSet, metric *Metric) string {

	m := (*metric)
	export := metadata.MetricExport{

		// Global
		Pods: set.Spec.Pods,

		// Metric
		MetricName:        m.Name(),
		MetricDescription: m.Description(),
		MetricOptions:     m.Options(),
		MetricListOptions: m.ListOptions(),
	}
	metadata, err := json.Marshal(export)
	if err != nil {
		logger.Errorf("Warning, error serializing spec metadata: %s", err.Error())
	}
	// We need to escape the quotes for printing in bash
	metadataEscaped := utils.EscapeCharacters(string(metadata))
	return fmt.Sprintf("METADATA START %s\nMETADATA END", metadataEscaped)
}

func init() {
	handle, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = handle.Sugar()
	defer handle.Sync()
}
