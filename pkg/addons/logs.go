/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/utils"
	"go.uber.org/zap"
)

// Consistent logging identifiers that should be echoed to have newline after
var (
	logger *zap.SugaredLogger
)

// Default metadata (in JSON) to also put at the top of addons
// That append to an entrypoint with their metadata
func Metadata(a Addon) string {

	export := metadata.MetricExport{
		MetricName:        a.Name(),
		MetricDescription: a.Description(),
		MetricOptions:     a.Options(),
		MetricListOptions: a.ListOptions(),
	}
	meta, err := json.Marshal(export)
	if err != nil {
		logger.Errorf("Warning, error serializing spec metadata: %s", err.Error())
	}
	// We need to escape the quotes for printing in bash
	metadataEscaped := utils.EscapeCharacters(string(meta))
	return fmt.Sprintf("ADDON METADATA START %s\nADDON METADATA END", metadataEscaped)
}

func init() {
	handle, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	logger = handle.Sugar()
	defer handle.Sync()
}
