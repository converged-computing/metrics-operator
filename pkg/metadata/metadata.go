/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metadata

import (
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Consistent logging identifiers that should be echoed to have newline after
var (
	Separator       = "METRICS OPERATOR TIMEPOINT"
	CollectionStart = "METRICS OPERATOR COLLECTION START"
	CollectionEnd   = "METRICS OPERATOR COLLECTION END"
	handle          *zap.Logger
	logger          *zap.SugaredLogger
)

// Metric Export is a flattened structure with minimal required metadata for now
// It would be nice if we could just dump everything.
type MetricExport struct {

	// Global
	Pods        int32 `json:"pods"`
	Completions int32 `json:"completions"`

	// Application
	ApplicationImage   string `json:"applicationImage,omitempty"`
	ApplicationCommand string `json:"applicationCommand,omitempty"`

	// Storage
	StorageVolumePath          string `json:"storageVolumePath,omitempty"`
	StorageVolumeHostPath      string `json:"storageVolumeHostPath,omitempty"`
	StorageVolumeSecretName    string `json:"storageVolumeSecretName,omitempty"`
	StorageVolumeClaimName     string `json:"storageVolumeClaimName,omitempty"`
	StorageVolumeConfigMapName string `json:"storageVolumeConfigMapName,omitempty"`

	// Metric
	MetricName        string                          `json:"metricName,omitempty"`
	MetricDescription string                          `json:"metricDescription,omitempty"`
	MetricType        string                          `json:"metricType,omitempty"`
	MetricOptions     map[string]intstr.IntOrString   `json:"metricOptions,omitempty"`
	MetricListOptions map[string][]intstr.IntOrString `json:"metricListOptions,omitempty"`
}

// Interactive returns a sleep infinity if interactive is true
func Interactive(interactive bool) string {
	if interactive {
		return "sleep infinity"
	}
	return ""
}
