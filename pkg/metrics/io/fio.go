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

// FIO means Flexible IO
// https://docs.gitlab.com/ee/administration/operations/filesystem_benchmarking.html

type Fio struct {
	metrics.StorageGeneric

	// Options
	testname  string
	blocksize string
	iodepth   int
	size      string
	directory string

	// extra commands for pre, post, etc.
	pre    string
	post   string
	prefix string
}

func (m Fio) Url() string {
	return "https://fio.readthedocs.io/en/latest/fio_doc.html"
}

// Set custom options / attributes for the metric
func (m *Fio) SetOptions(metric *api.Metric) {
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
	v, ok = metric.Options["prefix"]
	if ok {
		m.prefix = v.StrVal
	}
	v, ok = metric.Options["pre"]
	if ok {
		m.pre = v.StrVal
	}
	v, ok = metric.Options["post"]
	if ok {
		m.post = v.StrVal
	}
}

func (m Fio) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)

	preBlock := `#!/bin/bash
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
`
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		m.directory,
		m.pre,
		m.prefix,
		m.testname,
		m.blocksize,
		m.iodepth,
		m.size,
		metadata.CollectionStart,
		metadata.Separator,
	)

	postBlock := `
echo "%s"
# Run command here so it's after collection finish, but before removing the filename
%s 
%s rm -rf $filename
%s	
`

	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	postBlock = fmt.Sprintf(
		postBlock,
		metadata.CollectionEnd,
		m.post,
		m.prefix,
		interactive,
	)
	return m.StorageContainerSpec(preBlock, "$command", postBlock)
}

// Exported options and list options
func (m Fio) Options() map[string]intstr.IntOrString {
	return map[string]intstr.IntOrString{
		"testname":  intstr.FromString(m.testname),
		"blocksize": intstr.FromString(m.blocksize),
		"iodepth":   intstr.FromInt(m.iodepth),
		"size":      intstr.FromString(m.size),
		"directory": intstr.FromString(m.directory),
	}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: "io-fio",
		Summary:    "Flexible IO Tester (FIO)",
		Container:  "ghcr.io/converged-computing/metric-fio:latest",
	}
	storage := metrics.StorageGeneric{BaseMetric: base}
	fio := Fio{StorageGeneric: storage}
	metrics.Register(&fio)
}
