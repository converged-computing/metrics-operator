/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package sys

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	hwlocIdentifier = "sys-hwloc"
	hwlocSummary    = "install hwloc for inspecting hardware locality"
	hwlocContainer  = "ghcr.io/converged-computing/metric-hwloc:latest"
)

type Hwloc struct {
	metrics.SingleApplication

	// Custom Options
	commands []string
}

func (m Hwloc) Url() string {
	return "https://www.open-mpi.org/projects/hwloc/tutorials/20120702-POA-hwloc-tutorial.html"
}

func (m *Hwloc) Famliy() string {
	return metrics.SystemFamily
}

// Set custom options / attributes for the metric
func (m *Hwloc) SetOptions(metric *api.Metric) {

	m.Identifier = hwlocIdentifier
	m.Summary = hwlocSummary
	m.Container = hwlocContainer

	// Defaults for lstopo command
	m.ResourceSpec = &metric.Resources
	m.AttributeSpec = &metric.Attributes
	m.commands = []string{"lstopo architecture.png", "hwloc-ls machine.xml"}

	cmd, ok := metric.ListOptions["command"]
	if ok {
		m.commands = []string{}
		for _, val := range cmd {
			m.commands = append(m.commands, val.StrVal)
		}
	}
}

func (m Hwloc) ListOptions() map[string][]intstr.IntOrString {
	opts := map[string][]intstr.IntOrString{}
	for _, val := range m.commands {
		opts["commands"] = append(opts["commands"], intstr.FromString(val))
	}
	return opts
}

func (m Hwloc) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)

	// Assemble commands into separate things
	commands := ""
	for _, cmd := range m.commands {
		commands += fmt.Sprintf("\necho %s\n%s\n echo '%s'", cmd, cmd, metadata.Separator)
	}
	preBlock := `#!/bin/bash
echo "%s"	
. /root/.profile
export PATH=/opt/view/bin:$PATH
echo "%s"
%s
echo "%s"
ls
`

	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		metadata.CollectionStart,
		commands,
		metadata.CollectionEnd,
	)
	postBlock := fmt.Sprintf("\n%s\n", interactive)
	return m.ApplicationContainerSpec(preBlock, "", postBlock)
}

func init() {
	base := metrics.BaseMetric{
		Identifier: hwlocIdentifier,
		Summary:    hwlocSummary,
		Container:  hwlocContainer,
	}
	app := metrics.SingleApplication{BaseMetric: base}
	Hwloc := Hwloc{SingleApplication: app}
	metrics.Register(&Hwloc)
}
