package pcap

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	pcpav1 "github.com/alexanderConstantinescu/kubectl-pcap/api/v1"
)

const (
	GroupName = "pcap.k8s.io"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(pcpav1.Install)
	// Install is a function which adds every version of this group to a scheme
	Install = schemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

func Kind(kind string) schema.GroupKind {
	return schema.GroupKind{Group: GroupName, Kind: kind}
}
