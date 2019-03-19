package ecxfabric

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/terraform/helper/schema"
)

func createDXSession(d *schema.ResourceData) *directconnect.DirectConnect {
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	region := d.Get("region").(string)

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
	}))
	svc := directconnect.New(sess)

	return svc
}

func getDXConnectionByID(svc *directconnect.DirectConnect, id string) (*directconnect.Connection, error) {
	output, err := svc.DescribeConnections(&directconnect.DescribeConnectionsInput{
		ConnectionId: aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("Error describing connection (%v) in AWS: %s", id, err)
	}

	if len(output.Connections) != 1 {
		return nil, fmt.Errorf("Error describing connection (%v) in AWS. Expected 1 connection. Got %v", id, len(output.Connections))
	}

	return output.Connections[0], nil
}
