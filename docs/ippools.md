# IP Pools
An IP pool resource (IPPool) represents a collection of IP addresses from which Calico expects endpoint IPs to be assigned.

You can declare an ippool with this configuration : 

```hcl
resource "calico_ippool" "default" {
  metadata {
    name = "ippool-192.168.0.0"
  }
  spec {
    cidr = "192.168.0.0/24"
    ipip_mode = "Never"
    nat_outgoing = true
    disabled = false
  }
}
```

More informations on configuration [here](https://docs.projectcalico.org/v3.1/reference/calicoctl/resources/ippool).