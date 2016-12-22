package client

type (
	HorizontalPodAutoscalerInterface interface {
		CreateHorizontalPodAutoscaler(namespace string, item *HorizontalPodAutoscaler) (*HorizontalPodAutoscaler, error)
		GetHorizontalPodAutoscaler(namespace, name string) (result *HorizontalPodAutoscaler, err error)
		ListHorizontalPodAutoscalers(namespace string, opts *ListOptions) (*HorizontalPodAutoscalerList, error)
		DeleteHorizontalPodAutoscaler(namespace, name string) error
		UpdateHorizontalPodAutoscaler(namespace string, item *HorizontalPodAutoscaler) (*HorizontalPodAutoscaler, error)
	}

	// list of horizontal pod autoscaler objects.
	HorizontalPodAutoscalerList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`

		// list of horizontal pod autoscaler objects.
		Items []HorizontalPodAutoscaler `json:"items"`
	}

	// configuration of a horizontal pod autoscaler.
	HorizontalPodAutoscaler struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// behaviour of autoscaler. More info: http://releases.k8s.io/release-1.3/docs/devel/api-conventions.md#spec-and-status.
		Spec HorizontalPodAutoscalerSpec `json:"spec,omitempty"`

		// current information about the autoscaler.
		Status HorizontalPodAutoscalerStatus `json:"status,omitempty"`
	}

	// current status of a horizontal pod autoscaler
	HorizontalPodAutoscalerStatus struct {
		// most recent generation observed by this autoscaler.
		ObservedGeneration *int64 `json:"observedGeneration,omitempty"`

		// last time the HorizontalPodAutoscaler scaled the number of pods;
		// used by the autoscaler to control how often the number of pods is changed.
		LastScaleTime *Time `json:"lastScaleTime,omitempty"`

		// current number of replicas of pods managed by this autoscaler.
		CurrentReplicas int32 `json:"currentReplicas"`

		// desired number of replicas of pods managed by this autoscaler.
		DesiredReplicas int32 `json:"desiredReplicas"`

		// current average CPU utilization over all pods, represented as a percentage of requested CPU,
		// e.g. 70 means that an average pod is using now 70% of its requested CPU.
		CurrentCPUUtilizationPercentage *int32 `json:"currentCPUUtilizationPercentage,omitempty"`
	}

	// specification of a horizontal pod autoscaler.
	HorizontalPodAutoscalerSpec struct {
		// reference to scaled resource; horizontal pod autoscaler will learn the current resource consumption
		// and will set the desired number of pods by using its Scale subresource.
		ScaleTargetRef CrossVersionObjectReference `json:"scaleTargetRef"`
		// lower limit for the number of pods that can be set by the autoscaler, default 1.
		MinReplicas *int32 `json:"minReplicas,omitempty"`
		// upper limit for the number of pods that can be set by the autoscaler. It cannot be smaller than MinReplicas.
		MaxReplicas int32 `json:"maxReplicas"`
		// target average CPU utilization (represented as a percentage of requested CPU) over all the pods;
		// if not specified the default autoscaling policy will be used.
		TargetCPUUtilizationPercentage *int32 `json:"targetCPUUtilizationPercentage,omitempty"`
	}

	// CrossVersionObjectReference contains enough information to let you identify the referred resource.
	CrossVersionObjectReference struct {
		// Kind of the referent; More info: http://releases.k8s.io/release-1.3/docs/devel/api-conventions.md#types-kinds"
		Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
		// Name of the referent; More info: http://releases.k8s.io/release-1.3/docs/user-guide/identifiers.md#names
		Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
		// API version of the referent
		APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,3,opt,name=apiVersion"`
	}
)

// NewHorizontalPodAutoscalere creates a new HorizontalPodAutoscaler struct
func NewHorizontalPodAutoscaler(namespace, name string) *HorizontalPodAutoscaler {
	return &HorizontalPodAutoscaler{
		TypeMeta: TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling",
		},
		ObjectMeta: ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
	}
}
