package vsphere

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVSphereHostThumbprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereHostThumbprintRead,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The address of the ESXi to extract the thumbprint from.",
			},
		},
	}
}

func dataSourceVSphereHostThumbprintRead(d *schema.ResourceData, meta interface{}) error {
	config := &tls.Config{}
	config.InsecureSkipVerify = true
	conn, err := tls.Dial("tcp", d.Get("address").(string)+":443", config)
	if err != nil {
		return err
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	fingerprint := sha1.Sum(cert.Raw)

	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	d.SetId(buf.String())
	return nil
}
