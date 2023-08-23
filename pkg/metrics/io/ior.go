/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package io

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/jobs"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
)

// Ior means Flexible IO
// https://docs.gitlab.com/ee/administration/operations/filesystem_benchmarking.html

type Ior struct {
	jobs.StorageGeneric

	// Options
	workdir string
	command string
}

func (m Ior) Url() string {
	return "https://github.com/hpc/ior"
}

// Set custom options / attributes for the metric
func (m *Ior) SetOptions(metric *api.Metric) {
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set defaults for options
	m.command = "ior -w -r -o testfile"
	m.workdir = "/opt/ior"

	// https://ior.readthedocs.io/en/latest/
	// https://ior.readthedocs.io/en/latest/userDoc/tutorial.html
	// with mpirun mpirun -n 64 ./ior -t 1m -b 16m -s 16
	// For example commands
	command, ok := metric.Options["command"]
	if ok {
		m.command = command.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		m.workdir = workdir.StrVal
	}
}

// Generate the entrypoint for measuring the storage
func (m Ior) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Prepare metadata for set and separator
	metadata := metrics.Metadata(spec, metric)
	template := `#!/bin/bash
echo "%s"
# Directory (and filename) for test assuming other storage mounts
cd %s
echo "%s"
echo "%s"
%s
echo "%s"
%s
%s
%s
`
	script := fmt.Sprintf(
		template,
		metadata,
		m.workdir,
		metrics.CollectionStart,
		metrics.Separator,
		m.command,
		metrics.CollectionEnd,
		spec.Spec.Storage.Commands.Post,
		spec.Spec.Storage.Commands.Prefix,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)
	return []metrics.EntrypointScript{
		{Script: script},
	}

}

// Exported options and list options
func (m Ior) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"workdir": intstr.FromString(m.workdir),
		"command": intstr.FromString(m.command),
	}
}

func init() {
	storage := jobs.StorageGeneric{
		Identifier: "io-ior",
		Summary:    "HPC IO Benchmark",
		Container:  "ghcr.io/converged-computing/metric-ior:latest",
	}
	Ior := Ior{StorageGeneric: storage}
	metrics.Register(&Ior)
}
