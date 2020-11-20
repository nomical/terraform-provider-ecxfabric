package ecxfabric

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nomical/terraform-provider-ecxfabric/apiclient"
)

func resourceL2Connection() *schema.Resource {
	return &schema.Resource{
		Create: resourceL2ConnectionCreate,
		Read:   resourceL2ConnectionRead,
		Delete: resourceL2ConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"authorization_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"purchase_order_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notifications": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_port_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"primary_vlan_s_tag": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "Service VLAN tag",
				Required:    true,
				ForceNew:    true,
			},
			"profile_uuid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"seller_metro_code": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"seller_region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"speed": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"speed_unit": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_connection_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceL2ConnectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Client)
	input := apiclient.PostConnectionRequest{}

	if v, ok := d.GetOk("authorization_key"); ok {
		input.AuthorizationKey = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("notifications"); ok {
		var notifications []string
		for _, n := range v.([]interface{}) {
			notifications = append(notifications, n.(string))
		}

		input.Notifications = &notifications
	}
	if v, ok := d.GetOk("primary_name"); ok {
		input.PrimaryName = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("purchase_order_number"); ok {
		input.PurchaseOrderNumber = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("primary_port_uuid"); ok {
		input.PrimaryPortUUID = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("primary_vlan_s_tag"); ok {
		input.PrimaryVlanSTag = apiclient.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("profile_uuid"); ok {
		input.ProfileUUID = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("seller_metro_code"); ok {
		input.SellerMetroCode = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("seller_region"); ok {
		input.SellerRegion = apiclient.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("speed"); ok {
		input.Speed = apiclient.IntPtr(v.(int))
	}
	if v, ok := d.GetOk("speed_unit"); ok {
		input.SpeedUnit = apiclient.StringPtr(v.(string))
	}

	resp, err := client.CreateL2Connection(input)
	if err != nil {
		return fmt.Errorf("Error creating L2 connection: %s", err)
	}
	d.SetId(resp.PrimaryConnectionID)

	l2ConnStateConf := &resource.StateChangeConf{
		Pending: []string{
			apiclient.L2ConnectionStatusProvisioning,
		},
		Target: []string{
			apiclient.L2ConnectionStatusProvisioned,
		},
		Refresh:    l2ConnectionStatusRefreshFunc(client, resp.PrimaryConnectionID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for L2 Connection (%v) to be provisioned", resp.PrimaryConnectionID)
	_, err = l2ConnStateConf.WaitForState()
	if err != nil {
		return err
	}

	// Added provider status check here because if an l2_connection_aws_accepter is used there can be a race condition where the return aws_connection_id is not populated as AWS hasn't completed creating the DX connection in time.
	l2ProviderConnStateConf := &resource.StateChangeConf{
		Pending: []string{
			apiclient.L2ConnectionProviderStatusNotAvailable,
			apiclient.L2ConnectionProviderStatusProvisioning
		},
		Target: []string{
			apiclient.L2ConnectionProviderStatusPendingApproval,
		},
		Refresh: l2ConnectionProviderStatusRefreshFunc(client, resp.PrimaryConnectionID),
		Timeout: d.Timeout(schema.TimeoutCreate),
	}

	log.Printf("[DEBUG] Waiting for L2 Connection (%v) provider to create connection in state 'Pending Approval'", resp.PrimaryConnectionID)
	_, err = l2ProviderConnStateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceL2ConnectionRead(d, m)
}

func resourceL2ConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Client)
	uuid := d.Id()
	conn, err := client.ReadL2Connection(uuid)
	if err != nil {
		return err
	}

	d.Set("authorization_key", conn.AuthorizationKey)
	d.Set("notifications", conn.Notifications)
	d.Set("primary_name", conn.Name)
	d.Set("purchase_order_number", conn.PurchaseOrderNumber)
	d.Set("primary_port_uuid", conn.PortUUID)
	d.Set("primary_vlan_s_tag", conn.VlanSTag)
	d.Set("profile_uuid", conn.SellerServiceUUID)
	d.Set("seller_metro_code", conn.SellerMetroCode)
	d.Set("speed", conn.Speed)
	d.Set("speed_unit", conn.SpeedUnit)
	d.Set("aws_connection_id", conn.ActionDetails.ExtractValueFromActionRequiredDataItem("awsConnectionId"))

	return nil
}

func resourceL2ConnectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Client)
	uuid := d.Id()
	err := client.DeleteL2Connection(uuid)
	if err != nil {
		return err
	}

	l2ConnStateConf := &resource.StateChangeConf{
		Pending: []string{
			apiclient.L2ConnectionProviderStatusDeprovisioning,
		},
		Target: []string{
			apiclient.L2ConnectionProviderStatusDeprovisioned,
		},
		Refresh:    l2ConnectionProviderStatusRefreshFunc(client, uuid),
		Timeout:    d.Timeout(schema.TimeoutDefault),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for L2 Connection (%v) to be deprovisioned", uuid)
	_, err = l2ConnStateConf.WaitForState()
	if err != nil {
		return err
	}

	return err
}

func l2ConnectionStatusRefreshFunc(client *apiclient.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := client.ReadL2Connection(uuid)
		if err != nil {
			return nil, "", fmt.Errorf("Error reading L2 connection (%v): %s", uuid, err)
		}

		return conn, conn.Status, nil
	}
}

func l2ConnectionProviderStatusRefreshFunc(client *apiclient.Client, uuid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := client.ReadL2Connection(uuid)
		if err != nil {
			return nil, "", fmt.Errorf("Error reading L2 connection (%v): %s", uuid, err)
		}

		return conn, conn.ProviderStatus, nil
	}
}
