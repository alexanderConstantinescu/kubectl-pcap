package v1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_PCAP = map[string]string{
	"":       "PCAP does layer 4 packet tracing across a Kubernetes cluster. The CR can be created by a user with required cluster permissions and allows for a programmable interface to debugging a Kubernetes networking problem. This CRD is mainly intended to be used by applications/user experiencing consistent connection failures to their endpoint. Applications can create this CR automatically whenever they experience a failed connection attempt. This CRD is not intended for sporadic connection failures since those require much more advanced techniques, and usually capturing more packets than what this CRD aims at doing. The goal of this CRD is to ease a network engineers life and reduce the amount of data that needs to analyzed when debugging a network problem. This CRD focuses on atomic issues and hence is very poor at analyzing broader scope problems. The CRD allows for a source resource and destination fields to be specified. Note that most fields in this CR require that you DO NOT specify IP. Kubernetes objects are not keyed by IP, and this CRD attempts at aligning it with the Kubernetes paradigm rather than the networking one.\n\nCompatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).",
	"spec":   "spec is the definition of the desired packet capture.",
	"status": "status is the observed trace of the desired packet capture. Read-only.",
}

func (PCAP) SwaggerDoc() map[string]string {
	return map_PCAP
}

var map_PCAPDestination = map[string]string{
	"":            "PCAPDestination defines a destination for the network connection",
	"destination": "destination is a freetext valued destination used for the connection. This can be a DNS name or an IP address. In the case that a DNS name is specified: packet traces including the DNS resolution will be included in the capture.",
	"port":        "port is the destination port specified for connection. It is specified as optional the purpose of ICMP packets",
	"protocol":    "protocol is the protocol specified for the connection. It is specified as optional the purpose of ICMP packets",
}

func (PCAPDestination) SwaggerDoc() map[string]string {
	return map_PCAPDestination
}

var map_PCAPList = map[string]string{
	"":      "Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer). PCAPList is the list of PCAP.",
	"items": "List of PCAP.",
}

func (PCAPList) SwaggerDoc() map[string]string {
	return map_PCAPList
}

var map_PCAPSource = map[string]string{
	"":          "PCAPSource defines a source for the network connection",
	"node":      "node is the node name, as specified by the Kubernetes field: node.metadata.name",
	"pod":       "pod is the pod name, as specified by the Kubernetes field: pod.metadata.name",
	"namespace": "namespace is the pod's namespace, as specified by the Kubernetes field: namespace.metadata.name",
}

func (PCAPSource) SwaggerDoc() map[string]string {
	return map_PCAPSource
}

var map_PCAPSpec = map[string]string{
	"":            "PCAPSpec defines a source and destination for the network connection",
	"id":          "id is defined by the initiator of the injected network connection and signals to all agents which ID they should filter on.",
	"source":      "source is the source Kubernetes resource initiating a network connection. Packets will be injected into this network namespace.",
	"destination": "destination is the destination endpoint targeted by the network connection",
}

func (PCAPSpec) SwaggerDoc() map[string]string {
	return map_PCAPSpec
}

var map_PCAPStatus = map[string]string{
	"":           "PCAPStatus specifies the packet traces found across the cluster for the connection, together with a global condition status indicating PCAP state to all agents",
	"traces":     "traces contains all traces across the cluster of the packets related to the connection",
	"conditions": "conditions coordinates state and action for all agents on the cluster.",
}

func (PCAPStatus) SwaggerDoc() map[string]string {
	return map_PCAPStatus
}

var map_Trace = map[string]string{
	"timestamp": "timestamp is the packet trace timestamp as returned by libpcap",
	"node":      "node is the node name the packet trace was captured on, as specified by the Kubernetes field: node.metadata.name",
	"pod":       "pod is the pod name associated with the source or destination network interface. Note: this is only set in case the L2 packet trace references an non-host networked interface.",
	"srciface":  "srciface is the source network interace associated with the L2 packet capture",
	"dstiface":  "dstiface is the destination network interace associated with the L2 packet capture",
	"srcip":     "srcip is the source IP of the layer 3 packet trace",
	"dstip":     "dstip is the destination IP of the layer 3 packet trace",
	"srcport":   "srcport is the source port of the layer 4 packet trace",
	"dstport":   "dstport is the destination port of the layer 4 packet trace",
}

func (Trace) SwaggerDoc() map[string]string {
	return map_Trace
}

// AUTO-GENERATED FUNCTIONS END HERE
