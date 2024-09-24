/*
Copyright 2024 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package application

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/converged-computing/metrics-operator/pkg/metadata"
	metrics "github.com/converged-computing/metrics-operator/pkg/metrics"
	"github.com/converged-computing/metrics-operator/pkg/specs"
)

const (
	customIdentifier = "app-custom"
	customSummary    = "Provide a custom application for MPI trace"
)

type CustomApp struct {
	metrics.LauncherWorker
}

func (m CustomApp) Url() string {
	return "https://converged-computing.github.io/metrics-operator"
}

func (m CustomApp) Family() string {
	return metrics.ProxyAppFamily
}

// Set custom options / attributes for the metric
func (m *CustomApp) SetOptions(metric *api.Metric) {

	m.Identifier = customIdentifier
	m.Summary = customSummary

	// Ensure we set sole tenancy if desired
	st, ok := metric.Options["soleTenancy"]
	if ok && st.StrVal == "true" || st.StrVal == "yes" {
		m.SoleTenancy = true
	}

	// We require both a command and workdir
	m.SetDefaultOptions(metric)
	if m.Command == "" || m.Container == "" {
		fmt.Printf("Either \"command\" or \"container\" is not defined - this will not work as expected")
	}
}

// We don't know if the app can run on one node or not
func (m CustomApp) Validate(spec *api.MetricSet) bool {
	return true
}

// Exported options and list options
func (m CustomApp) Options() map[string]intstr.IntOrString {
	values := map[string]intstr.IntOrString{
		"command":     intstr.FromString(m.Command),
		"workdir":     intstr.FromString(m.Workdir),
		"soleTenancy": intstr.FromString("false"),
	}
	if m.SoleTenancy {
		values["soleTenancy"] = intstr.FromString("true")
	}
	return values
}

// Prepare containers with jobs and entrypoint scripts
func (m CustomApp) PrepareContainers(
	spec *api.MetricSet,
	metric *metrics.Metric,
) []*specs.ContainerSpec {

	// Metadata to add to beginning of run
	meta := metrics.Metadata(spec, metric)
	hosts := m.GetHostlist(spec)
	prefix := m.GetCommonPrefix(meta, m.Command, hosts)

	// Template blocks for launcher script
	preBlock := `
echo "%s"
`

	postBlock := `
echo "%s"
%s
`
	interactive := metadata.Interactive(spec.Spec.Logging.Interactive)
	preBlock = prefix + fmt.Sprintf(preBlock, metadata.Separator)
	postBlock = fmt.Sprintf(postBlock, metadata.CollectionEnd, interactive)

	// Entrypoint for the launcher
	launcherEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.LauncherScript),
		Path:    m.LauncherScript,
		Pre:     preBlock,
		Command: m.Command,
		Post:    postBlock,
	}

	// Entrypoint for the worker
	// Just has a sleep infinity added to the prefix
	workerEntrypoint := specs.EntrypointScript{
		Name:    specs.DeriveScriptKey(m.WorkerScript),
		Path:    m.WorkerScript,
		Pre:     prefix,
		Command: "sleep infinity",
	}

	// These are associated with replicated jobs via JobName
	launcherContainer := m.GetLauncherContainerSpec(launcherEntrypoint)
	workerContainer := m.GetWorkerContainerSpec(workerEntrypoint)

	// Return the script templates for each of launcher and worker
	return []*specs.ContainerSpec{&launcherContainer, &workerContainer}
}

func init() {
	base := metrics.BaseMetric{
		Identifier: customIdentifier,
		Summary:    customSummary,
	}
	launcher := metrics.LauncherWorker{BaseMetric: base}
	custom := CustomApp{LauncherWorker: launcher}
	metrics.Register(&custom)
}
