module github.com/cdiscount/terraform-provider-calico

go 1.16

require (
	github.com/hashicorp/terraform-plugin-sdk v1.17.2
	github.com/projectcalico/libcalico-go v1.7.2-0.20210513174936-6ccf0906db1d
	k8s.io/apimachinery v0.21.0-rc.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
