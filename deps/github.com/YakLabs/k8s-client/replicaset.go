package client

type (
	// ReplicaSetInterface has methods to work with ReplicaSet resources.
	ReplicaSetInterface interface {
		CreateReplicaSet(namespace string, item *ReplicaSet) (*ReplicaSet, error)
		GetReplicaSet(namespace, name string) (result *ReplicaSet, err error)
		ListReplicaSets(namespace string, opts *ListOptions) (*ReplicaSetList, error)
		DeleteReplicaSet(namespace, name string) error
		UpdateReplicaSet(namespace string, item *ReplicaSet) (*ReplicaSet, error)
	}

	// ReplicaSetList is a collection of ReplicaSets.
	ReplicaSetList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []ReplicaSet `json:"items"`
	}

	// ReplicaSet represents the configuration of a replica set.
	ReplicaSet struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`

		// Spec defines the desired behavior of this ReplicaSet.
		Spec ReplicaSetSpec `json:"spec,omitempty"`

		// Status is the current status of this ReplicaSet. This data may be
		// out of date by some window of time.
		Status ReplicaSetStatus `json:"status,omitempty"`
	}

	// ReplicaSetStatus represents the current status of a ReplicaSet.
	ReplicaSetStatus struct {
		// Replicas is the number of actual replicas.
		Replicas int32 `json:"replicas"`

		// The number of pods that have labels matching the labels of the pod template of the replicaset.
		FullyLabeledReplicas int32 `json:"fullyLabeledReplicas,omitempty"`

		// ObservedGeneration is the most recent generation observed by the controller.
		ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	}

	// ReplicaSetSpec is the specification of a replication set.
	ReplicaSetSpec struct {
		// Replicas is the number of desired replicas.
		Replicas int32 `json:"replicas"`

		// Selector is a label query over pods that should match the replica count.
		// Must match in order to be controlled.
		// If empty, defaulted to labels on pod template.
		// More info: http://releases.k8s.io/release-1.3/docs/user-guide/labels.md#label-selectors
		Selector *LabelSelector `json:"selector,omitempty"`

		// Template is the object that describes the pod that will be created if
		// insufficient replicas are detected.
		Template PodTemplateSpec `json:"template,omitempty"`
	}
)

// NewService creates a new ReplicaSet struct
func NewReplicaSet(namespace, name string) *ReplicaSet {
	return &ReplicaSet{
		TypeMeta: TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: ObjectMeta{
			Namespace:   namespace,
			Name:        name,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Spec: ReplicaSetSpec{
			Template: PodTemplateSpec{
				ObjectMeta: ObjectMeta{
					Labels: make(map[string]string),
				},
			},
		},
	}
}
