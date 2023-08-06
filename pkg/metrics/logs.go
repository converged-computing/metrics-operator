/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"encoding/json"
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/utils"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Consistent logging identifiers that should be echoed to have newline after
var (
	Separator       = "METRICS OPERATOR TIMEPOINT"
	CollectionStart = "METRICS OPERATOR COLLECTION START"
	CollectionEnd   = "METRICS OPERATOR COLLECTION END"
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

// Default metadata (in JSON) to also put at the top of logs for parsing
// I'd like to improve upon this manual approach, it's a bit messy.
func Metadata(set *api.MetricSet, metric *Metric) string {

	m := (*metric)
	export := MetricExport{

		// Global
		Pods:        set.Spec.Pods,
		Completions: set.Spec.Completions,

		// Application
		ApplicationImage:   set.Spec.Application.Image,
		ApplicationCommand: set.Spec.Application.Command,

		// Storage
		StorageVolumePath:          set.Spec.Storage.Volume.Path,
		StorageVolumeHostPath:      set.Spec.Storage.Volume.HostPath,
		StorageVolumeSecretName:    set.Spec.Storage.Volume.SecretName,
		StorageVolumeClaimName:     set.Spec.Storage.Volume.ClaimName,
		StorageVolumeConfigMapName: set.Spec.Storage.Volume.ConfigMapName,

		// Metric
		MetricName:        m.Name(),
		MetricDescription: m.Description(),
		MetricType:        m.Type(),
		MetricOptions:     m.Options(),
		MetricListOptions: m.ListOptions(),
	}
	metadata, err := json.Marshal(export)
	if err != nil {
		fmt.Printf("Warning, error serializing spec metadata: %s", err.Error())
	}
	// We need to escape the quotes for printing in bash
	metadataEscaped := utils.EscapeCharacters(string(metadata))
	return fmt.Sprintf("METADATA START %s\nMETADATA END", metadataEscaped)
}
