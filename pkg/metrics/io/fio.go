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

// FIO means Flexible IO
// https://docs.gitlab.com/ee/administration/operations/filesystem_benchmarking.html

type Fio struct {
	jobs.StorageGeneric

	// Options
	testname  string
	blocksize string
	iodepth   int
	size      string
	directory string
}

func (m Fio) Url() string {
	return "https://fio.readthedocs.io/en/latest/fio_doc.html"
}

// Set custom options / attributes for the metric
func (m *Fio) SetOptions(metric *api.Metric) {
	m.Rate = metric.Rate
	m.Completions = metric.Completions
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes

	// Set defaults for options
	m.testname = "test"
	m.blocksize = "4k"
	m.iodepth = 64
	m.size = "4G"
	m.directory = "/tmp"

	v, ok := metric.Options["testname"]
	if ok {
		m.testname = v.StrVal
	}
	v, ok = metric.Options["blocksize"]
	if ok {
		m.blocksize = v.StrVal
	}
	v, ok = metric.Options["size"]
	if ok {
		m.size = v.StrVal
	}
	v, ok = metric.Options["directory"]
	if ok {
		m.directory = v.StrVal
	}
	v, ok = metric.Options["iodepth"]
	if ok {
		m.iodepth = int(v.IntVal)
	}
}

// Generate the entrypoint for measuring the storage
func (m Fio) EntrypointScripts(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []metrics.EntrypointScript {

	// Prepare metadata for set and separator
	metadata := metrics.Metadata(spec, metric)
	template := `#!/bin/bash

echo "%s"
# Directory (and filename) for test assuming other storage mounts
filename=%s/test-$(cat /dev/urandom | tr -cd 'a-f0-9' | head -c 32)
# Run the pre-command here so it has access to the filename.
%s
command="%s fio --randrepeat=1 --ioengine=libaio --direct=1 --gtod_reduce=1 --name=%s --bs=%s --iodepth=%d --readwrite=randrw --rwmixread=75 --size=%s --filename=$filename --output-format=json"
echo "FIO COMMAND START"
echo $command
echo "FIO COMMAND END"
# FIO just has one command, we don't need to think about completions / etc!
echo "%s"
echo "%s"
$command
echo "%s"
# Run command here so it's after collection finish, but before removing the filename
%s 
%s rm -rf $filename
%s
`
	script := fmt.Sprintf(
		template,
		metadata,
		m.directory,
		spec.Spec.Storage.Commands.Pre,
		spec.Spec.Storage.Commands.Prefix,
		m.testname,
		m.blocksize,
		m.iodepth,
		m.size,
		metrics.CollectionStart,
		metrics.Separator,
		metrics.CollectionEnd,
		spec.Spec.Storage.Commands.Post,
		spec.Spec.Storage.Commands.Prefix,
		metrics.Interactive(spec.Spec.Logging.Interactive),
	)
	// The entrypoint is the entrypoint for the container, while
	// the command is expected to be what we are monitoring. Often
	// they are the same thing. We return an empty Name so it's automatically
	// assigned
	return []metrics.EntrypointScript{
		{Script: script},
	}

}

// Exported options and list options
func (m Fio) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"rate":        intstr.FromInt(int(m.Rate)),
		"completions": intstr.FromInt(int(m.Completions)),
		"testname":    intstr.FromString(m.testname),
		"blocksize":   intstr.FromString(m.blocksize),
		"iodepth":     intstr.FromInt(m.iodepth),
		"size":        intstr.FromString(m.size),
		"directory":   intstr.FromString(m.directory),
	}
}

func init() {
	storage := jobs.StorageGeneric{
		Identifier: "io-fio",
		Summary:    "Flexible IO Tester (FIO)",
		Container:  "ghcr.io/converged-computing/metric-fio:latest",
	}
	fio := Fio{StorageGeneric: storage}
	metrics.Register(&fio)
}
