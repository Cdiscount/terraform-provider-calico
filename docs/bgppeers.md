# BGP Peers

A BGP peer resource (BGPPeer) represents a remote BGP peer with which the node(s) in a Calico cluster will peer.
Configuring BGP peers allows you to peer a Calico network with your datacenter fabric (e.g. ToR). For more information on cluster layouts, see Calico’s documentation on L3 Topologies.

You can declare a BGP Peer with this configuration : 

```hcl
resource "calico_bgppeer" "default" {
  metadata {
    name = "Router1"
  }
  spec {
    node = "global"
    peer_ip = "192.168.0.5"
    as_number = 62523
  }
}
```

More informations [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/bgppeer) 