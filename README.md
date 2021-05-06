# kubectl-pcap

This project provides a cloud-native way of performing "distributed network tracing", I would like to emphasize the quotation marks since doing distributed tracing of network packets has a lot of subtle intricacies which, to the trained eye could seem like a bold statement.

This project uses [google/gopacket](https://github.com/google/gopacket) (which uses libpcap under the hood - which in turn uses BPF under its hood) to both inject and trace network packets across a Kubernetes cluster.

## Motivation

Network failures in a highly dynamical system such as Kubernetes is very common, most often difficult to debug and very state dependent. This project tries to solve both problems in a programmatic way:

- Difficult to debug: this project allows a cluster admin the possibility to inject tailored packets into a network namespace and trace its path to the packet destination (or as far as the packet can make it). "Tailored packets" here means that it sets specific bits in the TCP/IP and ethernet headers to be make it uniquely identifiable on a cluster. The ID is announced to all agents across the cluster which start their watchers looking for the packet in question, every agent monitors layer 2 - 4 packet data and reports the result into the custom CR result which is the uniquely identifiable and Kubernetes native construct used to report the global packet tracing result
- State dependent: the fact that this is using a Kubernetes CRD allows for an extensible and programmatic way for applications to create distributed packet traces automatically whenever they experience a failed network connection. The CRD can later be saved, replayed and sent to other - hopefully network teams responsible, for investigation.

As mentioned in the description, doing this has a lot of subtle difficulties. Both src/dst port can change depending on the kernel the packet passes through, most networking solutions use some form of encapsulation protocol for its overlay, and once you finish jumping over those hurdles you might run into conntrack.

The motivation behind this project is to allow networking engineers the ability to easily trace packets following the format:

```
+-----------+------+-------------------------------+-----+--------+--------+----------+----------+
| TIMESTAMP | NODE | ASSOCIATED NIC RESOURCE (POD) | NIC | SRC IP | DST IP | SRC PORT | DST PORT |
+-----------+------+-------------------------------+-----+--------+--------+----------+----------+
|           |      |                               |     |        |        |          |          |
+-----------+------+-------------------------------+-----+--------+--------+----------+----------+
|           |      |                               |     |        |        |          |          |
+-----------+------+-------------------------------+-----+--------+--------+----------+----------+
```

That will give somebody debugging a Kubernetes networking issue concise and concrete information as to the route taken by a packet cross nodes & pods. In the case of a failed connection, it will print all steps taken, which should make investigation easier.

Given this method: we can use this tool to investigate failures on the most broken of clusters, the only requirement is node-to-node connectivity.

Note: `timestamp` here will be the packet capture timestamp as returned from libpcap.  

## Network plugin

As a first version of this project I've been looking at generic layer 2 - 4 packet parsing. Kubernetes usually doesn't manage its own network, it "out-sources" the complexity to network plugins which often get blamed for the most proposterous of problems (sometimes with valid reasons). The goal going forward will be to have network plugin specific overlay protocols implemented to be able to have a more extensive packet trace (i.e: provide vxlan, geneve, BFD, etc) protocol metadata. This should help networking engineers working on OVN-Kubernetes, openshift-sdn, Antrea and the likes. Other network plugins could implement their dependent and additional protocols later on.

## Kubectl plugin

This project has a dedicated CRD used gather the tracing data. However end-users can use the provided `kubectl` plugin `pcap` to do the same, one BPF filter can be provided for the packet filtering and a "injection" filter for tracing metadata. Ex:

```
kubectl pcap "from node $NODE_X to $IP/$DNS and port 8080 and tcp" 
```

The following "actor" and resource primitives can be used for the injection filter

**actor primitives**

```
from - designated packet origin
to - designated packet destination
```

**resource primitives**

```
node - designated origin/destination Kubernetes node 
pod - designated origin/destination Kubernetes pod 
service - designated destination Kubernetes service 
```

Note: the `service` primitive can be specified as Kubernetes DNS name or direct IP address. In the case of a DNS name, obviously the hops taken will include the path to the DNS server and (hopefully) its response. #IBlameDNS. Also note that you can't do `from service`, because that makes no frikkin' sense.

There could also be cases where no resource primitive needs to accompany the actor primitive, for ex: `kubectl pcap "from pod $POD_X_NAMESPACE/$POD_X to google.com and port 443 and tcp"`, specifically `to google.com`, this is all handled by each agent's string parser.

The `pcap` argument format must follow "from ... to ... and" and is parsed into a CR which is created on the cluster and which all node agents will watch. The first `and` primitive will mark the start of the BPF filter passed down to libpcap. The node hosting the origin endpoint (node or pod) will use `nsenter` or `ip netns exec` to exec into the network namespace (in the case of a pod) and inject a packet of the desired spec. Alternatively, a user could simply create the CR themselves. All agents then watch for the unique identifier to filter the packet data on.

Note: since all parsed packet data will be stored in a CR, it won't guarantee a sorted order on timestamp with nice looking columns like the above. You can however read an existing CR using:

```
kubectl create/apply -f my_pcap_crd.yaml
```

With an already filled out stanza: no agent will start watching for the packet, but the resource will be created on the cluster and be presentable in a nice looking pcap table, like the one you saw above, using:

```
kubectl get pcap $MY_PCAP
```

## Injecting packets

The kernel syscalls today do not allow a lot of flexibility when it comes to modifying specific TCP/IP/ethernet header fields. The kernel stack usually takes care of ensuring the implementation of those protocols. To be able to work around this, the project could use the `WritePacketData` API provided by [google/gopacket](https://github.com/google/gopacket/blob/5d8084036064d14fd0a27ef1e10ec91c2526aa6f/pcap/pcap.go#L695). The inconvience of doing this will be that the packet injector used for this project to mark each packet with a unique ID, will need to implement [TCP/IP](https://www.ietf.org/rfc/rfc793.txt), as to ensure that all packets (no matter the connection state) have the marked ID for the entire transaction. Given the initial ideas here, the protocols to be implemented for the purpose of this project would likely be: [TCP](https://www.ietf.org/rfc/rfc793.txt)/[UDP](https://datatracker.ietf.org/doc/html/rfc8085) and [ICMP](https://datatracker.ietf.org/doc/html/rfc2236).

## libpcap vs eBPF

So far, libpcap has been envisioned as the underlying library allowing packet modification. The reason is: a larger portability of that library across OS'es and systems. eBPF could however achieve the same goal.

## Imagined scope

The scope of network connections would include:

- pod to pod connectivity
- pod to service connectivity
- pod to node connectivity
- node to pod connectivity
- pod to cluster egress connectivity
- node to cluster egress connectivity
- ingress to pod/node connectivity (for this, the local `kubectl` client would need to inject packets from the client's computer towards the cluster - which would then have the indentifier observable on the cluster for all agents to monitor)

## Background

The biggest problem with debugging network failures, I feel, is the lack of a unique identifier for each packet. One can quickly end up resorting to deploying `tcpdump` across all cluster nodes and capturing Gb of data from X amounts of nodes, only to look for one specific connection at moment Y (essentially a couple of Kb of data in the end). It's nauseating. Obviously we can't re-work the existing packet header definitions, but what if we could manipulate packets to have pre-determined and deterministic fields? For example, specify that a TCP packet will have an initial sequence number of 666? Doing that already determines what the subsequent sequence and acknowledgement numbers will be, allowing us to trace the 3-way handshake. BPF, libpcap and subsequently gopacket allows us to do this. Another identifier could be to set the URGENT field in the TCP header, which is a legacy field

## Future ideas

With such functionality one could also imagine auto-reporting networking issues. I.e: instead of `readinessProbes`/`livenessProbes`, one could maybe imagine `networkCanaryProbes`, which be able to use this CRD and functionality to trigger and test connectivty to required network endpoints, should the network probe fail: the CRD would be stored on the cluster, providing the networking engineer with a detailed report of all hops from source to target and at which point the connection failed.