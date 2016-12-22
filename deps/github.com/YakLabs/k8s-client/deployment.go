package client

import "github.com/YakLabs/k8s-client/intstr"

const (
	// Kill all existing pods before creating new ones.
	RecreateDeploymentStrategyType DeploymentStrategyType = "Recreate"

	// Replace the old RCs by new one using rolling update i.e gradually scale down the old RCs and scale up the new one.
	RollingUpdateDeploymentStrategyType DeploymentStrategyType = "RollingUpdate"
)

type (
	// DeploymentInterface has methods to work with Deployment resources.
	DeploymentInterface interface {
		CreateDeployment(namespace string, item *Deployment) (*Deployment, error)
		GetDeployment(namespace, name string) (result *Deployment, err error)
		ListDeployments(namespace string, opts *ListOptions) (*DeploymentList, error)
		DeleteDeployment(namespace, name string) error
		UpdateDeployment(namespace string, item *Deployment) (*Deployment, error)
	}

	Deployment struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Specification of the desired behavior of the Deployment.
		Spec DeploymentSpec `json:"spec,omitempty"`

		// Most recently observed status of the Deployment.
		Status DeploymentStatus `json:"status,omitempty"`
	}

	DeploymentSpec struct {
		// Number of desired pods. This is a pointer to distinguish between explicit
		// zero and not specified. Defaults to 1.
		Replicas int `json:"replicas,omitempty"`

		// Label selector for pods. Existing ReplicaSets whose pods are
		// selected by this will be the ones affected by this deployment.
		Selector *LabelSelector `json:"selector,omitempty"`

		// Template describes the pods that will be created.
		Template PodTemplateSpec `json:"template"`

		// The deployment strategy to use to replace existing pods with new ones.
		Strategy DeploymentStrategy `json:"strategy,omitempty"`

		// Minimum number of seconds for which a newly created pod should be ready
		// without any of its container crashing, for it to be considered available.
		// Defaults to 0 (pod will be considered available as soon as it is ready)
		MinReadySeconds int `json:"minReadySeconds,omitempty"`

		// The number of old ReplicaSets to retain to allow rollback.
		// This is a pointer to distinguish between explicit zero and not specified.
		RevisionHistoryLimit *int `json:"revisionHistoryLimit,omitempty"`

		// Indicates that the deployment is paused and will not be processed by the
		// deployment controller.
		Paused bool `json:"paused,omitempty"`
		// The config this deployment is rolling back to. Will be cleared after rollback is done.
		RollbackTo *RollbackConfig `json:"rollbackTo,omitempty"`
	}

	DeploymentStrategy struct {
		// Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.
		Type DeploymentStrategyType `json:"type,omitempty"`

		// Rolling update config params. Present only if DeploymentStrategyType =
		// RollingUpdate.
		//---
		// TODO: Update this to follow our convention for oneOf, whatever we decide it
		// to be.
		RollingUpdate *RollingUpdateDeployment `json:"rollingUpdate,omitempty"`
	}

	DeploymentStrategyType string

	// Spec to control the desired behavior of rolling update.
	RollingUpdateDeployment struct {
		// The maximum number of pods that can be unavailable during the update.
		// Value can be an absolute number (ex: 5) or a percentage of total pods at the start of update (ex: 10%).
		// Absolute number is calculated from percentage by rounding up.
		// This can not be 0 if MaxSurge is 0.
		// By default, a fixed value of 1 is used.
		// Example: when this is set to 30%, the old RC can be scaled down by 30%
		// immediately when the rolling update starts. Once new pods are ready, old RC
		// can be scaled down further, followed by scaling up the new RC, ensuring
		// that at least 70% of original number of pods are available at all times
		// during the update.
		MaxUnavailable intstr.IntOrString `json:"maxUnavailable,omitempty"`

		// The maximum number of pods that can be scheduled above the original number of
		// pods.
		// Value can be an absolute number (ex: 5) or a percentage of total pods at
		// the start of the update (ex: 10%). This can not be 0 if MaxUnavailable is 0.
		// Absolute number is calculated from percentage by rounding up.
		// By default, a value of 1 is used.
		// Example: when this is set to 30%, the new RC can be scaled up by 30%
		// immediately when the rolling update starts. Once old pods have been killed,
		// new RC can be scaled up further, ensuring that total number of pods running
		// at any time during the update is atmost 130% of original pods.
		MaxSurge intstr.IntOrString `json:"maxSurge,omitempty"`
	}

	RollbackConfig struct {
		// The revision to rollback to. If set to 0, rollbck to the last revision.
		Revision int64 `json:"revision,omitempty"`
	}

	// PodTemplateSpec describes the data a pod should have when created from a template
	PodTemplateSpec struct {
		// Metadata of the pods created from this template.
		ObjectMeta `json:"metadata,omitempty"`

		// Spec defines the behavior of a pod.
		Spec PodSpec `json:"spec,omitempty"`
	}

	DeploymentStatus struct {
		// The generation observed by the deployment controller.
		ObservedGeneration int64 `json:"observedGeneration,omitempty"`

		// Total number of non-terminated pods targeted by this deployment (their labels match the selector).
		Replicas int `json:"replicas,omitempty"`

		// Total number of non-terminated pods targeted by this deployment that have the desired template spec.
		UpdatedReplicas int `json:"updatedReplicas,omitempty"`

		// Total number of available pods (ready for at least minReadySeconds) targeted by this deployment.
		AvailableReplicas int `json:"availableReplicas,omitempty"`

		// Total number of unavailable pods targeted by this deployment.
		UnavailableReplicas int `json:"unavailableReplicas,omitempty"`
	}

	DeploymentList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		// Items is the list of deployments.
		Items []Deployment `json:"items"`
	}
)
