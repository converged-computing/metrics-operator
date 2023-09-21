/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package io

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

// Ior means Flexible IO
// https://docs.gitlab.com/ee/administration/operations/filesystem_benchmarking.html

type Ior struct {
	metrics.StorageGeneric

	// Options
	workdir string
	command string
	pre     string
	post    string
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
	v, ok := metric.Options["pre"]
	if ok {
		m.pre = v.StrVal
	}
	v, ok = metric.Options["post"]
	if ok {
		m.post = v.StrVal
	}
}

func (m Ior) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)

	preBlock := `#!/bin/bash
echo "%s"
# Directory (and filename) for test assuming other storage mounts
cd %s
echo "%s"
echo "%s"
`

	postBlock := `
echo "%s"
%s
%s
`
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		m.workdir,
		metadata.CollectionStart,
		metadata.Separator,
	)

	postBlock = fmt.Sprintf(
		postBlock,
		metadata.CollectionEnd,
		m.post,
		interactive,
	)
	return m.StorageContainerSpec(preBlock, m.command, postBlock)
}

// Exported options and list options
func (m Ior) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"workdir": intstr.FromString(m.workdir),
		"command": intstr.FromString(m.command),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "io-ior",
		Summary:    "HPC IO Benchmark",
		Container:  "ghcr.io/converged-computing/metric-ior:latest",
	}
	storage := metrics.StorageGeneric{BaseMetric: base}
	Ior := Ior{StorageGeneric: storage}
	metrics.Register(&Ior)
}
