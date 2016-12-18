package client

// Common object elements

type (
	UID             string
	FinalizerName   string
	ConditionStatus string

	TypeMeta struct {
		Kind       string `json:"kind,omitempty"`
		APIVersion string `json:"apiVersion,omitempty"`
	}

	ObjectMeta struct {
		Name              string            `json:"name,omitempty"`
		Namespace         string            `json:"namespace,omitempty"`
		SelfLink          string            `json:"selfLink,omitempty"`
		UID               UID               `json:"uid,omitempty"`
		ResourceVersion   string            `json:"resourceVersion,omitempty"`
		CreationTimestamp *Time             `json:"creationTimestamp,omitempty"`
		DeletionTimestamp *Time             `json:"deletionTimestamp,omitempty"`
		Generation        int64             `json:"generation,omitempty"`
		Labels            map[string]string `json:"labels,omitempty"`
		Annotations       map[string]string `json:"annotations,omitempty"`
	}

	ListMeta struct {
		SelfLink        string `json:"selfLink,omitempty"`
		ResourceVersion string `json:"resourceVersion,omitempty"`
	}

	ObjectReference struct {
		Kind            string `json:"kind,omitempty"`
		Namespace       string `json:"namespace,omitempty"`
		Name            string `json:"name,omitempty"`
		UID             UID    `json:"uid,omitempty"`
		APIVersion      string `json:"apiVersion,omitempty"`
		ResourceVersion string `json:"resourceVersion,omitempty"`
		FieldPath       string `json:"fieldPath,omitempty"`
	}

	LocalObjectReference struct {
		Name string `json:"name"`
	}

	ObjectFieldSelector struct {
		APIVersion string `json:"apiVersion"`
		FieldPath  string `json:"fieldPath"`
	}

	ConfigMapKeySelector struct {
		LocalObjectReference `json:",inline"`
		Key                  string `json:"key"`
	}

	SecretKeySelector struct {
		LocalObjectReference `json:",inline"`
		Key                  string `json:"key"`
	}

	LabelSelector struct {
		// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
		// map is equivalent to an element of matchExpressions, whose key field is "key", the
		// operator is "In", and the values array contains only "value". The requirements are ANDed.
		MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	}

	FieldSelector map[string]string

	Object interface {
		GetKind() string
		GetAnnotations() map[string]string
		GetLabels() map[string]string
		SetLabels(labels map[string]string)
	}

	NamespacedObject interface {
		Object
		GetNamespace() string
	}
	ListObject interface {
		Object
		GetItems()
	}
)

func (t *TypeMeta) GetKind() string {
	return t.Kind
}

func (o *ObjectMeta) GetNamespace() string {
	return o.Namespace
}

func (o *ObjectMeta) GetAnnotations() map[string]string {
	return o.Annotations
}

func (o *ObjectMeta) GetLabels() map[string]string {
	return o.Labels
}

func (o *ObjectMeta) SetLabels(labels map[string]string) {
	o.Labels = labels
}
