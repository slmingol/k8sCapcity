package main

import (
	"encoding/json"
	"fmt"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

func runDaemonMode(clusterInfo ClusterInfo) {
	daemonLog := DaemonLog{}
	daemonLog.UtilizationFactorPods = make(map[string]float64)
	daemonLog.UtilizationFactorMemoryRequests = make(map[string]float64)
	daemonLog.UtilizationFactorCPURequests = make(map[string]float64)

	for name, node := range clusterInfo.NodeInfo {
		if node.PrintOutput {
			clusterInfo.ClusterUsedCPURequests.Add(node.UsedCPURequests)
			clusterInfo.ClusterUsedCPU.Add(node.UsedCPU)
			clusterInfo.ClusterUsedMemoryRequests.Add(node.UsedMemoryRequests)
			clusterInfo.ClusterUsedMemory.Add(node.UsedMemory)
			clusterInfo.ClusterUsedPods = clusterInfo.ClusterUsedPods + node.UsedPods
			clusterInfo.ClusterUsedMemoryLimits.Add(node.UsedMemoryLimits)
			daemonLog.UtilizationFactorPods[name] = float64(node.UsedPods) / float64(node.AllocatablePods.Value())
			daemonLog.UtilizationFactorMemoryRequests[name] = float64(node.UsedMemoryRequests.Value()) / float64(node.AllocatableMemory.Value())
			daemonLog.UtilizationFactorCPURequests[name] = float64(node.UsedCPURequests.Value()) / float64(node.AllocatableCPU.Value())

		}
	}

	daemonLog.EventKind = "metric"
	daemonLog.EventModule = "k8s_quota"
	daemonLog.EventProvider = "k8sCapcity"
	daemonLog.EventType = "info"
	daemonLog.EventVersion = "03/06/2020-01"
	daemonLog.NodeLabel = clusterInfo.NodeLabel
	daemonLog.ResourceQuotaCPURequestCores = clusterInfo.RqclusterAllocatedRequestsCPU.Value()
	daemonLog.ResourceQuotaCPURequestMilliCores = clusterInfo.RqclusterAllocatedRequestsCPU.ScaledValue(resource.Milli)
	daemonLog.ResourceQuotaMemoryRequest = clusterInfo.RqclusterAllocatedRequestsMemory.Value()
	daemonLog.ResourceQuotaMemoryLimit = clusterInfo.RqclusterAllocatedLimitsMemory.Value()
	daemonLog.ResourceQuotaPods = clusterInfo.RqclusterAllocatedPods.Value()
	daemonLog.ContainerResourceCPURequestCores = clusterInfo.ClusterUsedCPURequests.Value()
	daemonLog.ContainerResourceCPURequestMilliCores = clusterInfo.ClusterUsedCPURequests.ScaledValue(resource.Milli)
	daemonLog.ContainerResourceMemoryRequest = clusterInfo.ClusterUsedMemoryRequests.Value()
	daemonLog.ContainerResourceMemoryLimit = clusterInfo.ClusterUsedMemoryLimits.Value()
	daemonLog.ContainerResourcePods = clusterInfo.ClusterUsedPods
	daemonLog.AllocatableMemoryTotal = clusterInfo.ClusterAllocatableMemory.Value()
	daemonLog.AllocatableMemoryNminusone = clusterInfo.ClusterAllocatableMemory.Value() - clusterInfo.NminusMemory.Value()
	daemonLog.AllocatableCPUTotal = clusterInfo.ClusterAllocatableCPU.Value()
	daemonLog.AllocatableCPUNminusone = clusterInfo.ClusterAllocatableCPU.Value() - clusterInfo.NminusCPU.Value()
	daemonLog.AllocatablePodsTotal = clusterInfo.ClusterAllocatablePods.Value()
	daemonLog.AllocatablePodsNminusone = clusterInfo.ClusterAllocatablePods.Value() - clusterInfo.NminusPods.Value()
	if float64(daemonLog.AllocatableMemoryTotal) == 0 {
		daemonLog.SubscriptionFactorMemoryRequestTotal = 0
	} else {
		daemonLog.SubscriptionFactorMemoryRequestTotal = float64(daemonLog.ResourceQuotaMemoryRequest) / float64(daemonLog.AllocatableMemoryTotal)
	}
	if daemonLog.AllocatableMemoryNminusone == 0 {
		daemonLog.SubscriptionFactorMemoryRequestNminusone = 0
	} else {
		daemonLog.SubscriptionFactorMemoryRequestNminusone = float64(daemonLog.ResourceQuotaMemoryRequest) / float64(daemonLog.AllocatableMemoryNminusone)
	}
	if float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli)) == 0 {
		daemonLog.SubscriptionFactorCPURequestTotal = 0
	} else {
		daemonLog.SubscriptionFactorCPURequestTotal = float64(daemonLog.ResourceQuotaCPURequestMilliCores) / float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli))
	}
	if float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli)-clusterInfo.NminusCPU.ScaledValue(resource.Milli)) == 0 {
		daemonLog.SubscriptionFactorCPURequestNminusone = 0
	} else {
		daemonLog.SubscriptionFactorCPURequestNminusone = float64(daemonLog.ResourceQuotaCPURequestMilliCores) / float64(clusterInfo.ClusterAllocatableCPU.ScaledValue(resource.Milli)-clusterInfo.NminusCPU.ScaledValue(resource.Milli))
	}
	if float64(daemonLog.AllocatablePodsTotal) == 0 {
		daemonLog.SubscriptionFactorPodsTotal = 0
	} else {
		daemonLog.SubscriptionFactorPodsTotal = float64(daemonLog.ResourceQuotaPods) / float64(daemonLog.AllocatablePodsTotal)
	}
	if float64(daemonLog.AllocatablePodsNminusone) == 0 {
		daemonLog.SubscriptionFactorPodsNminusone = 0
	} else {
		daemonLog.SubscriptionFactorPodsNminusone = float64(daemonLog.ResourceQuotaPods) / float64(daemonLog.AllocatablePodsNminusone)
	}
	if float64(daemonLog.AllocatablePodsTotal) == 0 {
		daemonLog.UtilizationFactorPodsTotal = 0
	} else {
		daemonLog.UtilizationFactorPodsTotal = float64(clusterInfo.ClusterUsedPods) / float64(daemonLog.AllocatablePodsTotal)
	}
	if float64(daemonLog.AllocatablePodsNminusone) == 0 {
		daemonLog.UtilizationFactorPodsNminusone = 0
	} else {
		daemonLog.UtilizationFactorPodsNminusone = float64(clusterInfo.ClusterUsedPods) / float64(daemonLog.AllocatablePodsNminusone)
	}
	if float64(daemonLog.AllocatableMemoryTotal) == 0 {
		daemonLog.UtilizationFactorMemoryRequestsTotal = 0
	} else {
		daemonLog.UtilizationFactorMemoryRequestsTotal = float64(daemonLog.ContainerResourceMemoryRequest) / float64(daemonLog.AllocatableMemoryTotal)
	}
	if float64(daemonLog.AllocatableMemoryNminusone) == 0 {
		daemonLog.UtilizationFactorMemoryRequestsNminusone = 0
	} else {
		daemonLog.UtilizationFactorMemoryRequestsNminusone = float64(daemonLog.ContainerResourceMemoryRequest) / float64(daemonLog.AllocatableMemoryNminusone)
	}
	if float64(daemonLog.AllocatableCPUTotal) == 0 {
		daemonLog.UtilizationFactorCPURequestsTotal = 0
	} else {
		daemonLog.UtilizationFactorCPURequestsTotal = float64(clusterInfo.ClusterUsedCPURequests.Value()) / float64(daemonLog.AllocatableCPUTotal)
	}
	if float64(daemonLog.AllocatableCPUNminusone) == 0 {
		daemonLog.UtilizationFactorCPURequestsNminusone = 0
	} else {
		daemonLog.UtilizationFactorCPURequestsNminusone = float64(clusterInfo.ClusterUsedCPURequests.Value()) / float64(daemonLog.AllocatableCPUNminusone)
	}
	daemonLog.AvailableMemoryRequestTotal = daemonLog.AllocatableMemoryTotal - daemonLog.ContainerResourceMemoryRequest
	daemonLog.AvailableMemoryRequestNminusone = daemonLog.AllocatableMemoryNminusone - daemonLog.ContainerResourceMemoryRequest
	daemonLog.AvailableCPURequestTotal = daemonLog.AllocatableCPUTotal - daemonLog.ContainerResourceCPURequestCores
	daemonLog.AvailableCPURequestNminusone = daemonLog.AllocatableCPUNminusone - daemonLog.ContainerResourceCPURequestCores
	daemonLog.AvailablePodsTotal = daemonLog.AllocatablePodsTotal - daemonLog.ContainerResourcePods
	daemonLog.AvailablePodsNminusone = daemonLog.AllocatablePodsNminusone - daemonLog.ContainerResourcePods
	result, err := json.Marshal(daemonLog)
	if err != nil {
		fmt.Printf("There was an error during json.Marshal, Error: %s\n", err)
		panic(err)
	}
	fmt.Println(string(result))
}
