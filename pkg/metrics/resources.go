/*
Copyright 2022-2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

import (
	"fmt"

	api "github.com/converged-computing/metrics-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// TODO add this to spec
// also find logging library
// getResourceGroup can return a ResourceList for either requests or limits
func getResourceGroup(items api.ContainerResource) (corev1.ResourceList, error) {

	logger.Info("🍅️ Resource", "items", items)
	list := corev1.ResourceList{}
	for key, unknownValue := range items {
		if unknownValue.Type == intstr.Int {

			value := unknownValue.IntVal
			logger.Info("🍅️ ResourceKey", "Key", key, "Value", value)
			limit, err := resource.ParseQuantity(fmt.Sprintf("%d", value))
			if err != nil {
				return list, err
			}

			if key == "memory" {
				list[corev1.ResourceMemory] = limit
			} else if key == "cpu" {
				list[corev1.ResourceCPU] = limit
			} else {
				list[corev1.ResourceName(key)] = limit
			}
		} else if unknownValue.Type == intstr.String {

			value := unknownValue.StrVal
			logger.Info("🍅️ ResourceKey", "Key", key, "Value", value)
			if key == "memory" {
				list[corev1.ResourceMemory] = resource.MustParse(value)
			} else if key == "cpu" {
				list[corev1.ResourceCPU] = resource.MustParse(value)
			} else {
				list[corev1.ResourceName(key)] = resource.MustParse(value)
			}
		}
	}
	return list, nil
}

// getContainerResources determines if any resources are requested via the spec
func getContainerResources(spec *api.ContainerResources) (corev1.ResourceRequirements, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources := corev1.ResourceRequirements{}

	// Limits
	limits, err := getResourceGroup(spec.Limits)
	if err != nil {
		logger.Error("🍅️ Resources for Container.Limits", err.Error())
		return resources, err
	}
	resources.Limits = limits

	// Requests
	requests, err := getResourceGroup(spec.Requests)
	if err != nil {
		logger.Error("🍅️ Resources for Container.Requests", err.Error())
		return resources, err
	}
	resources.Requests = requests
	return resources, nil

}

// getPodResources determines if any resources are requested via the spec
func getPodResources(set *api.MetricSet) (corev1.ResourceList, error) {

	// memory int, setCPURequest, setCPULimit, setGPULimit int64
	resources, err := getResourceGroup(set.Spec.Resources)
	if err != nil {
		logger.Error("🍅️ Resources for Pod.Resources", err.Error())
		return resources, err
	}
	return resources, nil
}
