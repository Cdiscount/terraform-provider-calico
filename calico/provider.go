package calico

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/projectcalico/libcalico-go/lib/apiconfig"
	client "github.com/projectcalico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/libcalico-go/lib/options"
	"log"
	"os"
)

type config struct {
	Client client.Interface
}

var ctx = context.Background()
var opts = options.SetOptions{}

// Provider is the provider for terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend_type": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_BACKEND_TYPE", ""),
				Description: "Type of the backend (etcdv3 or kubernetes)",
			},
			"etcd_endpoints": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_ENDPOINTS", ""),
				Description: "List of ETCD endpoints, separate by \",\"",
			},
			"etcd_username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_USERNAME", ""),
				Description: "Etcd username",
			},
			"etcd_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_PASSWORD", ""),
				Description: "Etcd password",
			},
			"etcd_ca_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_CA_CERT_FILE", ""),
				Description: "Cert of the CA (Only if SSL is enabled on your ETCD Server)",
			},
			"etcd_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_CERT_FILE", ""),
				Description: "Client cert file to connect to ETCD (Only if SSL is enabled on your ETCD Server)",
			},

			"etcd_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_ETCD_KEY_FILE", ""),
				Description: "Key file to connect to ETCD (Only if SSL is enabled on your ETCD Server)",
			},
			"kube_config": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_CONFIG", ""),
				Description: "Kube config to connect to Kubernetes server",
			},
			"kube_api_server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_API_SERVER", ""),
				Description: "APIServer to fetch (Optional if the apiserver is fill in the kube config)",
			},
			"kube_ca_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_CA_CERT_FILE", ""),
				Description: "Cert of the CA (Only if SSL is enabled on your Kubernetes)",
			},
			"kube_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_CERT_FILE", ""),
				Description: "Client cert file to connect to APIServer (Only if SSL is enabled on your Kubernetes)",
			},
			"kube_key_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_KEY_FILE", ""),
				Description: "Key file to connect to APIServer (Only if SSL is enabled on your Kubernetes)",
			},

			"kube_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncWithDefault("CALICO_KUBE_TOKEN", ""),
				Description: "Token to connect to the APIServer",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"calico_ippool":           resourceCalicoIpPool(),
			"calico_bgppeer":          resourceCalicoBgpPeer(),
			"calico_bgpconfiguration": resourceCalicoBgpConfiguration(),
		},

		ConfigureFunc: providerConfigure,
	}
}

// envDefaultFuncWithDefault get the env var to populate the config of the provider
func envDefaultFuncWithDefault(key string, defaultValue string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(key); v != "" {
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			return v, nil
		}
		return defaultValue, nil
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	calicoConfig := apiconfig.CalicoAPIConfig{}

	backendType := d.Get("backend_type").(string)

	switch backendType {
	case "etcdv3":
		calicoConfig.Spec.DatastoreType = apiconfig.DatastoreType(backendType)

		calicoConfig.Spec.EtcdEndpoints = d.Get("etcd_endpoints").(string)
		calicoConfig.Spec.EtcdUsername = d.Get("etcd_username").(string)
		calicoConfig.Spec.EtcdPassword = d.Get("etcd_password").(string)
		calicoConfig.Spec.EtcdCACertFile = d.Get("etcd_ca_cert_file").(string)
		calicoConfig.Spec.EtcdCertFile = d.Get("etcd_cert_file").(string)
		calicoConfig.Spec.EtcdKeyFile = d.Get("etcd_key_file").(string)

	case "kubernetes":
		calicoConfig.Spec.DatastoreType = apiconfig.DatastoreType(backendType)

		calicoConfig.Spec.Kubeconfig = d.Get("kube_config").(string)
		calicoConfig.Spec.K8sAPIEndpoint = d.Get("kube_api_server").(string)
		calicoConfig.Spec.K8sAPIToken = d.Get("kube_token").(string)
		calicoConfig.Spec.K8sCAFile = d.Get("kube_ca_cert_file").(string)
		calicoConfig.Spec.K8sCertFile = d.Get("kube_cert_file").(string)
		calicoConfig.Spec.K8sKeyFile = d.Get("kube_key_file").(string)

	default:
		return nil, fmt.Errorf("backend_type etcdv3 and kubernetes is the only supported backend at the moment")
	}

	calicoClient, err := client.New(calicoConfig)
	if err != nil {
		return nil, err
	}

	config := config{
		Client: calicoClient,
	}

	log.Printf("Configured: %#v", config)

	return config, nil
}
