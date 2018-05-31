package calico

import (
	"github.com/hashicorp/terraform/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/errors"
	"github.com/projectcalico/libcalico-go/lib/options"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func resourceCalicoBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceCalicoBgpPeerCreate,
		Read:   resourceCalicoBgpPeerRead,
		Update: resourceCalicoBgpPeerUpdate,
		Delete: resourceCalicoBgpPeerDelete,

		Schema: map[string]*schema.Schema{
			"metadata": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
					},
				},
			},
			"spec": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"peer_ip": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"as_number": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// dToBgpPeerSpec return the spec of the BgpPeer
func dToBgpPeerSpec(d *schema.ResourceData) (api.BGPPeerSpec, error) {
	spec := api.BGPPeerSpec{}

	node := d.Get("spec.0.node").(string)
	spec.Node = node

	peerIp := d.Get("spec.0.peer_ip").(string)
	spec.PeerIP = peerIp

	//TODO: Reactivate this field
	//asNumber := d.Get("spec.0.as_number")
	//spec.ASNumber = asNumber.(numorstring.ASNumber)

	return spec, nil
}

// dToBgpPeerSpec return the metadata of the BgpPeer
func dToBgpPeerTypeMeta(d *schema.ResourceData) (meta.ObjectMeta, error) {
	objectMeta := meta.ObjectMeta{}

	objectMeta.Name = d.Get("metadata.0.name").(string)

	return objectMeta, nil
}

// resourceCalicoBgpPeerCreate create a new BgpPeer in Calico
func resourceCalicoBgpPeerCreate(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpPeerInterface := calicoClient.BGPPeers()

	BgpPeer, err := createBgpPeerApiRequest(d)
	if err != nil {
		return err
	}

	_, err = BgpPeerInterface.Create(ctx, BgpPeer, opts)
	if err != nil {
		return err
	}

	d.SetId(BgpPeer.ObjectMeta.Name)
	return resourceCalicoBgpPeerRead(d, m)
}

// resourceCalicoBgpPeerRead get a specific BgpPeer
func resourceCalicoBgpPeerRead(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpPeerInterface := calicoClient.BGPPeers()

	nameBgpPeer := d.Get("metadata.0.name").(string)

	BgpPeer, err := BgpPeerInterface.Get(ctx, nameBgpPeer, options.GetOptions{})
	log.Printf("Obj: %+v", d)

	// Handle endpoint does not exist
	if err != nil {
		if _, ok := err.(errors.ErrorResourceDoesNotExist); ok {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId(nameBgpPeer)
	d.Set("metadata.0.name", BgpPeer.ObjectMeta.Name)
	d.Set("spec.0.node", BgpPeer.Spec.Node)
	d.Set("spec.0.peer_ip", BgpPeer.Spec.PeerIP)
	d.Set("spec.0.as_number", BgpPeer.Spec.ASNumber.String())

	return nil
}

// resourceCalicoBgpPeerUpdate update an BgpPeer in Calico
func resourceCalicoBgpPeerUpdate(d *schema.ResourceData, m interface{}) error {
	d.Partial(false)

	calicoClient := m.(config).Client
	BgpPeerInterface := calicoClient.BGPPeers()

	BgpPeer, err := createBgpPeerApiRequest(d)
	if err != nil {
		return err
	}

	_, err = BgpPeerInterface.Update(ctx, BgpPeer, opts)
	if err != nil {
		return err
	}

	return nil
}

// resourceCalicoBgpPeerDelete delete an BgpPeer in Calico
func resourceCalicoBgpPeerDelete(d *schema.ResourceData, m interface{}) error {
	calicoClient := m.(config).Client
	BgpPeerInterface := calicoClient.BGPPeers()

	nameBgpPeer := d.Get("metadata.0.name").(string)

	_, err := BgpPeerInterface.Delete(ctx, nameBgpPeer, options.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// createApiRequest prepare the request of creation and update
func createBgpPeerApiRequest(d *schema.ResourceData) (*api.BGPPeer, error) {
	// Set Spec to BgpPeer Spec
	spec, err := dToBgpPeerSpec(d)
	if err != nil {
		return nil, err
	}

	// Set Metadata to Kubernetes Metadata
	objectMeta, err := dToBgpPeerTypeMeta(d)
	if err != nil {
		return nil, err
	}

	// Create a new BGP Peer, with TypeMeta filled in
	// Then, fill the metadata and the spec
	newBgpPeer := api.NewBGPPeer()
	newBgpPeer.ObjectMeta = objectMeta
	newBgpPeer.Spec = spec

	return newBgpPeer, nil
}
