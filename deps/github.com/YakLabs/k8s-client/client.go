package client

type (
	// Client is an interface that represents a kubernetes client
	Client interface {
		ConfigMapInterface
		DeploymentInterface
		HorizontalPodAutoscalerInterface
		NamespaceInterface
		NodeInterface
		PodInterface
		ReplicaSetInterface
		SecretInterface
		ServiceAccountInterface
		ServiceInterface
	}

	ListOptions struct {
		LabelSelector LabelSelector
		FieldSelector FieldSelector
	}
)
