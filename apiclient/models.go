package apiclient

type oAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"user_name"`
	Password     string `json:"user_password"`
	GrantType    string `json:"grant_type"`
}

type oAuthResponse struct {
	AccessToken         string `json:"access_token"`
	TokenTimeout        string `json:"token_timeout"`
	Username            string `json:"user_name"`
	TokenType           string `json:"token_type"`
	RefreshToken        string `json:"refresh_token"`
	RefreshTokenTimeout string `json:"refresh_token_timeout"`
}

type PostConnectionResponse struct {
	Message               string `json:"message"`               // example: Connection created successfully
	PrimaryConnectionID   string `json:"primaryConnectionId"`   // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	SecondaryConnectionID string `json:"secondaryConnectionId"` // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	Status                string `json:"status"`                // example: SUCCESS
}

type PostConnectionRequest struct {
	AuthorizationKey       *string   `json:"authorizationKey,omitempty"` // example: 444111000222
	NamedTag               *string   `json:"namedTag,omitempty"`         // example: Private
	Notifications          *[]string `json:"notifications,omitempty"`
	PrimaryName            *string   `json:"primaryName,omitempty"`            // example: v3-api-test-pri
	PrimaryPortUUID        *string   `json:"primaryPortUUID,omitempty"`        // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	PrimaryVlanCTag        *int      `json:"primaryVlanCTag,omitempty"`        // example: 602
	PrimaryVlanSTag        *int      `json:"primaryVlanSTag,omitempty"`        // example: 601
	PrimaryZSidePortUUID   *string   `json:"primaryZSidePortUUID,omitempty"`   // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	PrimaryZSideVlanCTag   *int      `json:"primaryZSideVlanCTag,omitempty"`   // example: 101
	PrimaryZSideVlanSTag   *int      `json:"primaryZSideVlanSTag,omitempty"`   // example: 301
	ProfileUUID            *string   `json:"profileUUID,omitempty"`            // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	PurchaseOrderNumber    *string   `json:"purchaseOrderNumber,omitempty"`    // example: 312456323
	SecondaryName          *string   `json:"secondaryName,omitempty"`          // example: v3-api-test-sec1
	SecondaryPortUUID      *string   `json:"secondaryPortUUID,omitempty"`      // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	SecondaryVlanCTag      *int      `json:"secondaryVlanCTag,omitempty"`      // example: 501
	SecondaryVlanSTag      *int      `json:"secondaryVlanSTag,omitempty"`      // example: 501
	SecondaryZSidePortUUID *string   `json:"secondaryZSidePortUUID,omitempty"` // example: xxxxx192-xx70-xxxx-xx04-xxxxxxxa37xx
	SecondaryZSideVlanCTag *int      `json:"secondaryZSideVlanCTag,omitempty"` // example: 102
	SecondaryZSideVlanSTag *int      `json:"secondaryZSideVlanSTag,omitempty"` // example: 302
	SellerMetroCode        *string   `json:"sellerMetroCode,omitempty"`        // example: SV
	SellerRegion           *string   `json:"sellerRegion,omitempty"`           // example: us-west-1
	Speed                  *int      `json:"speed,omitempty"`                  // example: 50
	SpeedUnit              *string   `json:"speedUnit,omitempty"`              // example: MB
}

type GetConnectionByUUIDResponse struct {
	UUID                   string        `json:"uuid"`                  // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	AsideEncapsulation     string        `json:"asideEncapsulation"`    // example: dot1q
	AuthorizationKey       string        `json:"authorizationKey"`      // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	BillingTier            string        `json:"billingTier"`           // example: Up to 500MB
	BuyerOrganizationName  string        `json:"buyerOrganizationName"` // example: Forsythe Solutions Group, Inc.
	CreatedBy              string        `json:"createdBy"`             // example: sandboxuser@example-company.com
	CreatedByEmail         string        `json:"createdByEmail"`        // example: sandboxuser@example-company.com
	CreatedByFullName      string        `json:"createdByFullName"`     // example: Sandbox User
	CreatedDate            string        `json:"createdDate"`           // example: 2017-09-26T22:46:24.312Z
	DeletedBy              string        `json:"deletedBy"`             // example: user@domain.com
	DeletedByEmail         string        `json:"deletedByEmail"`        // example: user@domain.com
	DeletedDate            string        `json:"deletedDate"`           // example: 2017-09-26T22:46:24.312Z
	LastUpdatedBy          string        `json:"lastUpdatedBy"`         // example: sandboxuser@example-company.com
	LastUpdatedByEmail     string        `json:"lastUpdatedByEmail"`    // example: sandboxuser@example-company.com
	LastUpdatedByFullName  string        `json:"lastUpdatedByFullName"` // example: Sandbox User
	LastUpdatedDate        string        `json:"lastUpdatedDate"`       // example: 2017-09-26T23:01:46Z
	MetroCode              string        `json:"metroCode"`             // example: CH
	MetroDescription       string        `json:"metroDescription"`      // example: Chicago
	Name                   string        `json:"name"`                  // example: Test-123
	NamedTag               string        `json:"namedTag"`              // example: Private
	Notifications          []string      `json:"notifications"`
	PortName               string        `json:"portName"`               // example: TEST-CH2-CX-SEC-01
	PortUUID               string        `json:"portUUID"`               // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	Private                bool          `json:"private"`                // example: false,
	ProviderStatus         string        `json:"providerStatus"`         // example: PROVISIONED
	PurchaseOrderNumber    string        `json:"purchaseOrderNumber"`    // example: O-1234567890
	RedundancyGroup        string        `json:"redundancyGroup"`        // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	RedundancyType         string        `json:"redundancyType"`         // example: secondary
	RedundantUUID          string        `json:"redundantUUID"`          // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	Remote                 bool          `json:"remote"`                 // example: false
	Self                   bool          `json:"self"`                   // example: false,
	SellerMetroCode        string        `json:"sellerMetroCode"`        // example: CH
	SellerMetroDescription string        `json:"sellerMetroDescription"` // example: Chicago
	SellerOrganizationName string        `json:"sellerOrganizationName"` // example: EQUINIX-CLOUD-EXCHANGE
	SellerServiceName      string        `json:"sellerServiceName"`      // example: XYZ Cloud Service
	SellerServiceUUID      string        `json:"sellerServiceUUID"`      // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	SellerRegion           string        `json:"sellerRegion"`           // example: eu-west-1
	Speed                  int           `json:"speed"`                  // example: 500
	SpeedUnit              string        `json:"speedUnit"`              // example: MB
	Status                 string        `json:"status"`                 // example: PROVISIONED
	VlanSTag               int           `json:"vlanSTag"`               // example: 1015
	ZSidePortName          string        `json:"zSidePortName"`          // example: TEST-CHG-06GMR-Tes-2-TES-C
	ZSidePortUUID          string        `json:"zSidePortUUID"`          // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
	ZSideVlanCTag          int           `json:"zSideVlanCTag"`          // example: 515
	ZSideVlanSTag          int           `json:"zSideVlanSTag"`          // example: 2
	ActionDetails          ActionDetails `json:"actionDetails"`
}

type ActionDetails []ActionDetail

type ActionDetail struct {
	ActionType         string                   `json:"actionType"`
	OperationID        string                   `json:"operationId"`
	ActionMessage      string                   `json:"actionMessage"`
	ActionRequiredData []ActionRequiredDataItem `json:"actionRequiredData"`
}

type ActionRequiredDataItem struct {
	Key               string `json:"key"`
	Label             string `json:"label"`
	Value             string `json:"value"`
	Editable          bool   `json:"editable"`
	ValidationPattern string `json:"validationPattern"`
}

type DeleteConnectionResponse struct {
	Message             string `json:"message"`             // example: Message
	PrimaryConnectionID string `json:"primaryConnectionId"` // example: xxxxx191-xx70-xxxx-xx04-xxxxxxxa37xx
}
