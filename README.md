# Provider Calico for Terraform

Compatible Calico V3.

## Installation
You can find all the releases [here](https://github.com/cdiscount/terraform-provider-calico/releases). The project is available for Linux, Mac and Windows as binaries and as deb and rpm packages.

If you want to build the project by yourself, you can. The prerequisites are :

* Go (version 1.9+)
* Glide

```bash
git clone https://github.com/cdiscount.com/terraform-provider-calico
make build
```

## Documentation

* [Configure the provider](docs/provider.md)
* Resources
  * [IP Pools](docs/ippools.md)
  * [BGP Peers](docs/bgppeers.md) 
  