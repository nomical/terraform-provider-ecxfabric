package ecxfabric

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nomical/terraform-provider-ecxfabric/apiclient"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ECXFABRIC_CLIENT_ID", nil),
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ECXFABRIC_CLIENT_SECRET", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ECXFABRIC_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ECXFABRIC_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ecxfabric_l2_connection":              resourceL2Connection(),
			"ecxfabric_l2_connection_aws_accepter": resourceL2ConnectionAwsAccepter(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	client, err := apiclient.New()
	if err != nil {
		return nil, err
	}

	err = client.Authenticate(clientID, clientSecret, username, password)
	if err != nil {
		return nil, err
	}

	return client, err
}
