package v1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PCAP does layer 4 packet tracing across a Kubernetes cluster. The CR can be
// created by a user with required cluster permissions and allows for a
// programmable interface to debugging a Kubernetes networking problem. This CRD
// is mainly intended to be used by applications/user experiencing consistent
// connection failures to their endpoint. Applications can create this CR
// automatically whenever they experience a failed connection attempt. This CRD
// is not intended for sporadic connection failures since those require much
// more advanced techniques, and usually capturing more packets than what this
// CRD aims at doing. The goal of this CRD is to ease a network engineers life
// and reduce the amount of data that needs to analyzed when debugging a network
// problem. This CRD focuses on atomic issues and hence is very poor at
// analyzing broader scope problems. The CRD allows for a source resource and
// destination fields to be specified. Note that most fields in this CR require
// that you DO NOT specify IP. Kubernetes objects are not keyed by IP, and this
// CRD attempts at aligning it with the Kubernetes paradigm rather than the
// networking one.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=pcaps,scope=Cluster
type PCAP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// spec is the definition of the desired packet capture.
	// +kubebuilder:validation:Required
	// +required
	Spec PCAPSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// status is the observed trace of the desired packet capture. Read-only.
	// +kubebuilder:validation:Optional
	// +optional
	Status PCAPStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PCAPSpec defines a source and destination for the network connection
// +k8s:openapi-gen=true
type PCAPSpec struct {
	// id is defined by the initiator of the injected network connection and
	// signals to all agents which ID they should filter on.
	// +kubebuilder:validation:Required
	// +required
	ID string `json:"id" protobuf:"bytes,1,opt,name=id"`
	// source is the source Kubernetes resource initiating a network connection.
	// Packets will be injected into this network namespace.
	// +kubebuilder:validation:Required
	// +required
	Source PCAPSource `json:"source" protobuf:"bytes,2,opt,name=source"`
	// destination is the destination endpoint targeted by the network connection
	// +kubebuilder:validation:Required
	// +required
	Destination PCAPDestination `json:"destination" protobuf:"bytes,3,opt,name=destination"`
}

// PCAPSource defines a source for the network connection
// +k8s:openapi-gen=true
type PCAPSource struct {
	// node is the node name, as specified by the Kubernetes field: node.metadata.name
	// +kubebuilder:validation:Optional
	// +optional
	Node string `json:"node,omitempty" protobuf:"bytes,1,opt,name=node"`
	// pod is the pod name, as specified by the Kubernetes field: pod.metadata.name
	// +kubebuilder:validation:Optional
	// +optional
	Pod string `json:"pod,omitempty" protobuf:"bytes,2,opt,name=pod"`
	// namespace is the pod's namespace, as specified by the Kubernetes field: namespace.metadata.name
	// +kubebuilder:validation:Optional
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

// PCAPDestination defines a destination for the network connection
// +k8s:openapi-gen=true
type PCAPDestination struct {
	// destination is a freetext valued destination used for the
	// connection. This can be a DNS name or an IP address. In the
	// case that a DNS name is specified: packet traces including the DNS
	// resolution will be included in the capture.
	// +kubebuilder:validation:Required
	// +required
	Destination string `json:"destination" protobuf:"bytes,1,opt,name=destination"`
	// port is the destination port specified for connection. It is specified as
	// optional the purpose of ICMP packets
	// +kubebuilder:validation:Optional
	// +optional
	Port int32 `json:"port" protobuf:"varint,2,opt,name=port"`
	// protocol is the protocol specified for the connection. It is specified as
	// optional the purpose of ICMP packets
	// +kubebuilder:validation:Optional
	// +optional
	Protocol corev1.Protocol `json:"protocol" protobuf:"bytes,3,opt,name=protocol,casttype=k8s.io/api/core/v1.Protocol"`
}

type PCAPTimeout time.Duration

const (
	// PCAPTerminateTimeout is the termination timeout for all node agents'
	// packet capture. It is specified so that all agents do not continuously
	// keep packet tracing in case an error occurs with the coordination of
	// conditions.
	PCAPTerminateTimeout PCAPTimeout = PCAPTimeout(time.Second * 30)
)

type PCAPCondition string

const (
	// AgentsReady is the condition type indicating that all node agents have
	// started their packet capture. Each agent will increment
	// ObservedGeneration by one. It is only paired with the
	// ConditionStatus: True, and is only set to True once ObservedGeneration
	// equals the amount of cluster nodes. Only the source sets
	// ConditionStatus: True
	PCAPAgentsReadyCondition PCAPCondition = "AgentsReady"
	// Terminate is the condition type indicating to all agents that they should
	// terminate their packet capture. All agents will terminate their packet
	// capture PCAPTerminateTimeout seconds after starting it in any case.
	// It is paired with the following ConditionStatus:
	// - True, indicating a finished or error execution from the source agent
	// - False, indicating on-going execution
	PCAPTerminateCondition PCAPCondition = "Terminate"
)

// PCAPStatus specifies the packet traces found across the cluster for the
// connection, together with a global condition status indicating PCAP state to
// all agents
// +k8s:openapi-gen=true
type PCAPStatus struct {
	// traces contains all traces across the cluster of the packets related to
	// the connection
	// +kubebuilder:validation:Optional
	// +optional
	Traces []Trace `json:"traces,omitempty" protobuf:"bytes,1,rep,name=traces"`
	// conditions coordinates state and action for all agents on the cluster.
	// +kubebuilder:validation:Required
	// +required
	Conditions []metav1.Condition `json:"conditions" protobuf:"bytes,2,rep,name=conditions"`
}

// +k8s:openapi-gen=true
type Trace struct {
	// timestamp is the packet trace timestamp as returned by libpcap
	// +kubebuilder:validation:Required
	// +required
	Timestamp metav1.Time `json:"timestamp" protobuf:"bytes,1,opt,name=timestamp"`
	// node is the node name the packet trace was captured on, as specified by
	// the Kubernetes field: node.metadata.name
	// +kubebuilder:validation:Required
	// +required
	Node string `json:"node" protobuf:"bytes,2,opt,name=node"`
	// pod is the pod name associated with the source or destination network
	// interface. Note: this is only set in case the L2 packet trace references
	// an non-host networked interface.
	// +kubebuilder:validation:Optional +optional
	Pod string `json:"pod,omitempty" protobuf:"bytes,3,opt,name=pod"`
	// srciface is the source network interace associated with the L2 packet capture
	// +kubebuilder:validation:Required
	// +required
	SourceInterface string `json:"srciface" protobuf:"bytes,4,opt,name=srciface"`
	// dstiface is the destination network interace associated with the L2 packet capture
	// +kubebuilder:validation:Required
	// +required
	DestinationInterface string `json:"dstiface" protobuf:"bytes,5,opt,name=dstiface"`
	// srcip is the source IP of the layer 3 packet trace
	// +kubebuilder:validation:Required
	// +required
	SourceIP string `json:"srcip" protobuf:"bytes,6,opt,name=srcip"`
	// dstip is the destination IP of the layer 3 packet trace
	// +kubebuilder:validation:Required
	// +required
	DestinationIP string `json:"dstip" protobuf:"bytes,7,opt,name=dstip"`
	// srcport is the source port of the layer 4 packet trace
	// +kubebuilder:validation:Required
	// +required
	SourcePort int32 `json:"srcport" protobuf:"varint,8,opt,name=srcport"`
	// dstport is the destination port of the layer 4 packet trace
	// +kubebuilder:validation:Required
	// +required
	DestinationPort int32 `json:"dstport" protobuf:"varint,9,opt,name=dstport"`
}

// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=pcap
// PCAPList is the list of PCAP.
type PCAPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of PCAP.
	Items []PCAP `json:"items" protobuf:"bytes,2,rep,name=items"`
}
