package ecxfabric

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceL2ConnectionAwsAccepter() *schema.Resource {
	return &schema.Resource{
		Create: resourceL2ConnectionAwsAccepterCreate,
		Read:   resourceL2ConnectionAwsAccepterRead,
		Delete: resourceL2ConnectionAwsAccepterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
				Sensitive:   true,
				ForceNew:    true,
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
				Sensitive:   true,
				ForceNew:    true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
				ForceNew: true,
			},
			"aws_device": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_device_v2": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connection_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"jumbo_frame_capable": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_account": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"partner_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceL2ConnectionAwsAccepterCreate(d *schema.ResourceData, m interface{}) error {
	svc := createDXSession(d)
	connectionID := d.Get("connection_id").(string)

	input := &directconnect.ConfirmConnectionInput{
		ConnectionId: aws.String(connectionID),
	}
	_, err := svc.ConfirmConnection(input)
	if err != nil {
		return fmt.Errorf("Error accepting connection (%v) in AWS: %s", connectionID, err)
	}

	d.SetId(connectionID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{
			directconnect.ConnectionStateOrdering,
			directconnect.ConnectionStatePending,
		},
		Target: []string{
			directconnect.ConnectionStateAvailable,
		},
		Refresh:    l2ConnectionAwsAccepterRefreshFunc(svc, connectionID),
		Timeout:    d.Timeout(schema.TimeoutDefault),
		Delay:      1 * time.Minute,
		MinTimeout: 10 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for DX Connection (%v) to be available in AWS", connectionID)
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceL2ConnectionAwsAccepterRead(d, m)
}

func resourceL2ConnectionAwsAccepterRead(d *schema.ResourceData, m interface{}) error {
	svc := createDXSession(d)
	connectionID := d.Id()

	conn, err := getDXConnectionByID(svc, connectionID)
	if err != nil {
		return err
	}

	d.Set("aws_device", aws.StringValue(conn.AwsDevice))
	d.Set("aws_device_v2", aws.StringValue(conn.AwsDeviceV2))
	d.Set("bandwidth", aws.StringValue(conn.Bandwidth))
	d.Set("connection_id", aws.StringValue(conn.ConnectionId))
	d.Set("connection_name", aws.StringValue(conn.ConnectionName))
	d.Set("connection_state", aws.StringValue(conn.ConnectionState))
	d.Set("jumbo_frame_capable", aws.BoolValue(conn.JumboFrameCapable))
	d.Set("location", aws.StringValue(conn.Location))
	d.Set("owner_account", aws.StringValue(conn.OwnerAccount))
	d.Set("partner_name", aws.StringValue(conn.PartnerName))
	d.Set("region", aws.StringValue(conn.Region))
	d.Set("vlan", aws.Int64Value(conn.Vlan))

	return nil
}

func resourceL2ConnectionAwsAccepterDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func l2ConnectionAwsAccepterRefreshFunc(svc *directconnect.DirectConnect, connectionID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := getDXConnectionByID(svc, connectionID)
		if err != nil {
			return nil, "", fmt.Errorf("Error reading DX connection (%v): %s", connectionID, err)
		}
		connState := aws.StringValue(conn.ConnectionState)

		return conn, connState, nil
	}
}
