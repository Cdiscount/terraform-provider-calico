# BGP Peers
A BGP configuration resource (BGPConfiguration) represents BGP specific configuration options for the cluster.

You can declare a BGP Peer with this configuration : 

```hcl
resource "calico_bgpconfiguration" "default" {
  metadata {
    name = "Config1"
  }
  spec {
    log_severity_screen = "Warning"
    node_to_node_mesh_enabled = false
    as_number = 62523
  }
}
```

More informations on configuration [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/bgpconfig)