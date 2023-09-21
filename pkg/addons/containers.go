/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Container addons are typically for applications
type ApplicationAddon struct {
	AddonBase

	// Container image
	image string

	// Name for container
	name string

	// command to execute
	command string

	// Working Directory
	workdir string

	// Entrypoint of container, if different from command
	entrypoint string

	// A pull secret for the application container
	pullSecret string

	// Resources include limits and requests for the application
	resources map[string]map[string]intstr.IntOrString

	// Container Spec has attributes for the container
	// Do we run this in privileged mode?
	privileged bool
}

// Validate we have an executable provided, and args and optional
func (a *ApplicationAddon) Validate() bool {
	if a.image == "" {
		logger.Error("The application addon requires a container 'image'.")
		return false
	}
	if a.name == "" {
		a.name = "app-addon"
	}
	if a.command == "" {
		logger.Error("The application addon requires a container 'command'.")
		return false
	}
	return true
}

// AssembleContainers adds the addon application container
func (a ApplicationAddon) AssembleContainers() []specs.ContainerSpec {
	return []specs.ContainerSpec{{
		Image:      a.image,
		Name:       a.name,
		WorkingDir: a.workdir,
		Command:    strings.Split(a.command, " "),
		// TODO these need to be mapped from m.resources
		Resources: &api.ContainerResources{},
		Attributes: &api.ContainerSpec{
			SecurityContext: api.SecurityContext{
				Privileged: a.privileged,
				// TODO add the caps here ptrace admin
			},
		},
	}}
}

// Set custom options / attributes for the metric
func (a *ApplicationAddon) SetDefaultOptions(metric *api.MetricAddon) {
	a.resources = map[string]map[string]intstr.IntOrString{}

	image, ok := metric.Options["image"]
	if ok {
		a.image = image.StrVal
	}
	command, ok := metric.Options["command"]
	if ok {
		a.command = command.StrVal
	}
	entrypoint, ok := metric.Options["entrypoint"]
	if ok {
		a.entrypoint = entrypoint.StrVal
	}
	pullSecret, ok := metric.Options["pullSecret"]
	if ok {
		a.pullSecret = pullSecret.StrVal
	}
	workdir, ok := metric.Options["workdir"]
	if ok {
		a.workdir = workdir.StrVal
	}
	priv, ok := metric.Options["privileged"]
	if ok {
		if priv.StrVal == "true" || priv.StrVal == "yes" {
			a.privileged = true
		}
	}
	resources, ok := metric.MapOptions["resourceLimits"]
	if ok {
		a.resources["limits"] = map[string]intstr.IntOrString{}
		for key, value := range resources {
			a.resources["limits"][key] = value
		}
	}
	resources, ok = metric.MapOptions["resourceRequests"]
	if ok {
		a.resources["requests"] = map[string]intstr.IntOrString{}
		for key, value := range resources {
			a.resources["requests"][key] = value
		}
	}
	if a.entrypoint == "" {
		a.setDefaultEntrypoint()
	}
}

// Set the default entrypoint
func (a *ApplicationAddon) setDefaultEntrypoint() {
	a.entrypoint = fmt.Sprintf("/metrics_operator/%s-entrypoint.sh", a.Identifier)
}

// Calling the default allows a custom application that uses this to do the same
func (a *ApplicationAddon) SetOptions(metric *api.MetricAddon) {
	a.SetDefaultOptions(metric)
}

// Underlying function that can be shared
func (a *ApplicationAddon) DefaultOptions() map[string]intstr.IntOrString {
	values := map[string]intstr.IntOrString{
		"image":      intstr.FromString(a.image),
		"workdir":    intstr.FromString(a.workdir),
		"entrypoint": intstr.FromString(a.entrypoint),
		"command":    intstr.FromString(a.command),
	}
	if a.privileged {
		values["privileged"] = intstr.FromString("true")
	} else {
		values["privileged"] = intstr.FromString("false")
	}
	return values
}

// Exported options and list options
func (a *ApplicationAddon) Options() map[string]intstr.IntOrString {
	return a.DefaultOptions()
}

// Return formatted map options
func (a *ApplicationAddon) MapOptions() map[string]map[string]intstr.IntOrString {
	requests := map[string]intstr.IntOrString{}
	limits := map[string]intstr.IntOrString{}
	for k, value := range a.resources["limits"] {
		limits[k] = value
	}
	for k, value := range a.resources["requests"] {
		requests[k] = value
	}
	return map[string]map[string]intstr.IntOrString{
		"resourceLimits":   limits,
		"resourceRequests": requests,
	}
}

func init() {

	// Config map volume type
	base := AddonBase{
		Identifier: "application",
		Summary:    "basic application (container) type",
	}
	app := ApplicationAddon{AddonBase: base}
	Register(&app)
}
