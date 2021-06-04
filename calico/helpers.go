package calico

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	api "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/numorstring"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// Generate the metadata of the resource
func dToTypeMeta(d *schema.ResourceData) (meta.ObjectMeta, error) {
	objectMeta := meta.ObjectMeta{}

	objectMeta.Name = dToString(d, "metadata.0.name")

	return objectMeta, nil
}

// Convert field to IPIPMode type
func dToIpIpMode(d *schema.ResourceData, field string) api.IPIPMode {
	f := api.IPIPMode(dToString(d, field))
	return f
}

// Convert field to ASNumber type
func dToAsNumber(d *schema.ResourceData, field string) numorstring.ASNumber {
	f := numorstring.ASNumber(uint(d.Get(field).(int)))
	return f
}

// Convert field to string
func dToString(d *schema.ResourceData, field string) string {
	f := d.Get(field).(string)
	return f
}

// Convert field to bool
func dToBool(d *schema.ResourceData, field string) bool {
	f := d.Get(field).(bool)
	return f
}

// Convert field to int
func dToInt(d *schema.ResourceData, field string) int {
	f := d.Get(field).(int)
	return f
}

// Convert field to uint32
func dToUint32(d *schema.ResourceData, field string) uint32 {
	f := uint32(d.Get(field).(uint))
	return f
}

// Convert field to duration
func dToDuration(d *schema.ResourceData, field string) meta.Duration {
	f := time.Duration(dToInt(d, field))
	f2 := meta.Duration{f}

	return f2
}
