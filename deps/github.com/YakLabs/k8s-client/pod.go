package client

import "github.com/YakLabs/k8s-client/intstr"

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"
	DNSClusterFirst        DNSPolicy     = "ClusterFirst"
	DNSDefault             DNSPolicy     = "Default"
	// URISchemeHTTP means that the scheme used will be http://
	URISchemeHTTP URIScheme = "HTTP"
	// URISchemeHTTPS means that the scheme used will be https://
	URISchemeHTTPS URIScheme = "HTTPS"
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// PullAlways means that kubelet always attempts to pull the latest image.  Container will fail If the pull fails.
	PullAlways PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image.  Container will fail if the image isn't present
	PullNever PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	PullIfNotPresent PullPolicy = "IfNotPresent"
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	PodUnknown PodPhase = "Unknown"
)

type (
	PodInterface interface {
		CreatePod(namespace string, item *Pod) (*Pod, error)
		GetPod(namespace, name string) (result *Pod, err error)
		ListPods(namespace string, opts *ListOptions) (*PodList, error)
		DeletePod(namespace, name string) error
		UpdatePod(namespace string, item *Pod) (*Pod, error)
	}

	Pod struct {
		TypeMeta   `json:",inline"`
		ObjectMeta `json:"metadata,omitempty"`
		Spec       PodSpec   `json:"spec,omitempty"`
		Status     PodStatus `json:"status,omitempty"`
	}

	PodList struct {
		TypeMeta `json:",inline"`
		ListMeta `json:"metadata,omitempty"`
		Items    []Pod `json:"items"`
	}

	PodSpec struct {
		Volumes                       []Volume               `json:"volumes"`
		Containers                    []Container            `json:"containers"`
		RestartPolicy                 RestartPolicy          `json:"restartPolicy,omitempty"`
		TerminationGracePeriodSeconds *int64                 `json:"terminationGracePeriodSeconds,omitempty"`
		ActiveDeadlineSeconds         *int64                 `json:"activeDeadlineSeconds,omitempty"`
		DNSPolicy                     DNSPolicy              `json:"dnsPolicy,omitempty"`
		NodeSelector                  map[string]string      `json:"nodeSelector,omitempty"`
		ServiceAccountName            string                 `json:"serviceAccountName"`
		NodeName                      string                 `json:"nodeName,omitempty"`
		ImagePullSecrets              []LocalObjectReference `json:"imagePullSecrets,omitempty"`
	}

	PodStatus struct {
		Phase             PodPhase          `json:"phase,omitempty"`
		Conditions        []PodCondition    `json:"conditions,omitempty"`
		Message           string            `json:"message,omitempty"`
		Reason            string            `json:"reason,omitempty"`
		HostIP            string            `json:"hostIP,omitempty"`
		PodIP             string            `json:"podIP,omitempty"`
		StartTime         *Time             `json:"startTime,omitempty"`
		ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
	}

	PodCondition struct {
		Type               PodConditionType `json:"type"`
		Status             ConditionStatus  `json:"status"`
		LastProbeTime      Time             `json:"lastProbeTime,omitempty"`
		LastTransitionTime Time             `json:"lastTransitionTime,omitempty"`
		Reason             string           `json:"reason,omitempty"`
		Message            string           `json:"message,omitempty"`
	}

	PodConditionType string
	StorageMedium    string
	RestartPolicy    string
	DNSPolicy        string
	Protocol         string
	URIScheme        string

	Volume struct {
		Name         string `json:"name"`
		VolumeSource `json:",inline,omitempty"`
	}

	VolumeSource struct {
		EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty"`
		Secret   *SecretVolumeSource   `json:"secret,omitempty"`
	}

	EmptyDirVolumeSource struct {
		Medium StorageMedium `json:"medium,omitempty"`
	}

	SecretVolumeSource struct {
		SecretName string `json:"secretName,omitempty"`
	}

	Container struct {
		Name                   string               `json:"name"`
		Image                  string               `json:"image"`
		Command                []string             `json:"command,omitempty"`
		Args                   []string             `json:"args,omitempty"`
		WorkingDir             string               `json:"workingDir,omitempty"`
		Ports                  []ContainerPort      `json:"ports,omitempty"`
		Env                    []EnvVar             `json:"env,omitempty"`
		Resources              ResourceRequirements `json:"resources,omitempty"`
		VolumeMounts           []VolumeMount        `json:"volumeMounts,omitempty"`
		LivenessProbe          *Probe               `json:"livenessProbe,omitempty"`
		ReadinessProbe         *Probe               `json:"readinessProbe,omitempty"`
		Lifecycle              *Lifecycle           `json:"lifecycle,omitempty"`
		TerminationMessagePath string               `json:"terminationMessagePath,omitempty"`
		ImagePullPolicy        PullPolicy           `json:"imagePullPolicy"`
		Stdin                  bool                 `json:"stdin,omitempty"`
		StdinOnce              bool                 `json:"stdinOnce,omitempty"`
		TTY                    bool                 `json:"tty,omitempty"`
	}

	ContainerPort struct {
		Name          string   `json:"name,omitempty"`
		HostPort      int      `json:"hostPort,omitempty"`
		ContainerPort int      `json:"containerPort"`
		Protocol      Protocol `json:"protocol,omitempty"`
		HostIP        string   `json:"hostIP,omitempty"`
	}

	EnvVar struct {
		Name      string        `json:"name"`
		Value     string        `json:"value,omitempty"`
		ValueFrom *EnvVarSource `json:"valueFrom,omitempty"`
	}

	EnvVarSource struct {
		FieldRef        *ObjectFieldSelector  `json:"fieldRef,omitempty"`
		ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
		SecretKeyRef    *SecretKeySelector    `json:"secretKeyRef,omitempty"`
	}

	// VolumeMount describes a mounting of a Volume within a container.
	VolumeMount struct {
		// Required: This must match the Name of a Volume [above].
		Name string `json:"name"`
		// Optional: Defaults to false (read-write).
		ReadOnly bool `json:"readOnly,omitempty"`
		// Required. Must not contain ':'.
		MountPath string `json:"mountPath"`
	}

	// Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
	Probe struct {
		// The action taken to determine the health of a container
		Handler `json:",inline"`
		// Length of time before health checking is activated.  In seconds.
		InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`
		// Length of time before health checking times out.  In seconds.
		TimeoutSeconds int `json:"timeoutSeconds,omitempty"`
		// How often (in seconds) to perform the probe.
		PeriodSeconds int `json:"periodSeconds,omitempty"`
		// Minimum consecutive successes for the probe to be considered successful after having failed.
		// Must be 1 for liveness.
		SuccessThreshold int `json:"successThreshold,omitempty"`
		// Minimum consecutive failures for the probe to be considered failed after having succeeded.
		FailureThreshold int `json:"failureThreshold,omitempty"`
	}

	// Handler defines a specific action that should be taken TODO: pass structured data to these actions, and document that data here.
	Handler struct {
		// One and only one of the following should be specified.
		// Exec specifies the action to take.
		Exec *ExecAction `json:"exec,omitempty"`
		// HTTPGet specifies the http request to perform.
		HTTPGet *HTTPGetAction `json:"httpGet,omitempty"`
		// TCPSocket specifies an action involving a TCP port.
		// TODO: implement a realistic TCP lifecycle hook
		TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
	}

	// ExecAction describes a "run in container" action.
	ExecAction struct {
		// Command is the command line to execute inside the container, the working directory for the
		// command  is root ('/') in the container's filesystem.  The command is simply exec'd, it is
		// not run inside a shell, so traditional shell instructions ('|', etc) won't work.  To use
		// a shell, you need to explicitly call out to that shell.
		Command []string `json:"command,omitempty"`
	}

	// HTTPGetAction describes an action based on HTTP Get requests.
	HTTPGetAction struct {
		// Optional: Path to access on the HTTP server.
		Path string `json:"path,omitempty"`
		// Required: Name or number of the port to access on the container.
		Port intstr.IntOrString `json:"port,omitempty"`
		// Optional: Host name to connect to, defaults to the pod IP. You
		// probably want to set "Host" in httpHeaders instead.
		Host string `json:"host,omitempty"`
		// Optional: Scheme to use for connecting to the host, defaults to HTTP.
		Scheme URIScheme `json:"scheme,omitempty"`
		// Optional: Custom headers to set in the request. HTTP allows repeated headers.
		HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty"`
	}

	// HTTPHeader describes a custom header to be used in HTTP probes
	HTTPHeader struct {
		// The header field name
		Name string `json:"name"`
		// The header field value
		Value string `json:"value"`
	}

	// ResourceRequirements describes the compute resource requirements.
	ResourceRequirements struct {
		// Limits describes the maximum amount of compute resources allowed.
		Limits ResourceList `json:"limits,omitempty"`
		// Requests describes the minimum amount of compute resources required.
		// If Request is omitted for a container, it defaults to Limits if that is explicitly specified,
		// otherwise to an implementation-defined value
		Requests ResourceList `json:"requests,omitempty"`
	}

	// ResourceList is a set of (resource name, quantity) pairs.
	ResourceList map[ResourceName]string

	// ResourceName is the name identifying various resources in a ResourceList.
	ResourceName string

	// TCPSocketAction describes an action based on opening a socket
	TCPSocketAction struct {
		// Required: Port to connect to.
		Port intstr.IntOrString `json:"port,omitempty"`
	}

	// Lifecycle describes actions that the management system should take in response to container lifecycle events. For the PostStart and PreStop lifecycle handlers, management of the container blocks until the action is complete, unless the container process fails, in which case the handler is aborted.
	Lifecycle struct {
		// PostStart is called immediately after a container is created.  If the handler fails, the container
		// is terminated and restarted.
		PostStart *Handler `json:"postStart,omitempty"`
		// PreStop is called immediately before a container is terminated.  The reason for termination is
		// passed to the handler.  Regardless of the outcome of the handler, the container is eventually terminated.
		PreStop *Handler `json:"preStop,omitempty"`
	}

	// PullPolicy describes a policy for if/when to pull a container image
	PullPolicy string

	// PodPhase is a label for the condition of a pod at the current time.
	PodPhase string

	// ContainerStatus is the status of a single container within a pod.
	ContainerStatus struct {
		// Each container in a pod must have a unique name.
		Name                 string         `json:"name"`
		State                ContainerState `json:"state,omitempty"`
		LastTerminationState ContainerState `json:"lastState,omitempty"`
		// Ready specifies whether the container has passed its readiness check.
		Ready bool `json:"ready"`
		// Note that this is calculated from dead containers.  But those containers are subject to
		// garbage collection.  This value will get capped at 5 by GC.
		RestartCount int    `json:"restartCount"`
		Image        string `json:"image"`
		ImageID      string `json:"imageID"`
		ContainerID  string `json:"containerID,omitempty"`
	}

	// ContainerState holds a possible state of container. Only one of its members may be specified. If none of them is specified, the default one is ContainerStateWaiting.
	ContainerState struct {
		Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
		Running    *ContainerStateRunning    `json:"running,omitempty"`
		Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
	}

	ContainerStateWaiting struct {
		// A brief CamelCase string indicating details about why the container is in waiting state.
		Reason string `json:"reason,omitempty"`
		// A human-readable message indicating details about why the container is in waiting state.
		Message string `json:"message,omitempty"`
	}

	ContainerStateRunning struct {
		StartedAt Time `json:"startedAt,omitempty"`
	}

	ContainerStateTerminated struct {
		ExitCode    int    `json:"exitCode"`
		Signal      int    `json:"signal,omitempty"`
		Reason      string `json:"reason,omitempty"`
		Message     string `json:"message,omitempty"`
		StartedAt   Time   `json:"startedAt,omitempty"`
		FinishedAt  Time   `json:"finishedAt,omitempty"`
		ContainerID string `json:"containerID,omitempty"`
	}
)
