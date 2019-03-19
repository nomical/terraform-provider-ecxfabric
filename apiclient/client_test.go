package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
	server *httptest.Server
	client *Client

	clientID            = "clientID"
	clientSecret        = "clientSecret"
	username            = "username"
	password            = "password"
	expectedAccessToken = "xxxxBxxitwxxxxx8xxRxxxxxR2xx"
)

func setup(t *testing.T) func() {
	router = mux.NewRouter()
	server = httptest.NewServer(router)
	router.HandleFunc(authEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// These responses have been tested live on api.equinix.com. Equinix documentation doesn't reflect actual responses!
		var oAuthRequest oAuthRequest
		if err := json.NewDecoder(r.Body).Decode(&oAuthRequest); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		if oAuthRequest.Username != username || oAuthRequest.Password != password {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		if oAuthRequest.ClientID != clientID || oAuthRequest.ClientSecret != clientSecret {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("oauth2/v1/token/token.json"))
	}).Methods("POST")
	router.HandleFunc(ecxV3L2ConnectionsEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		expectedAuthHeader := "Bearer " + expectedAccessToken
		actualAuthHeader := r.Header.Get("Authorization")
		if actualAuthHeader != expectedAuthHeader {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		expected := fixture("ecx/v3/l2/connections/create-request.json")
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		actual := buf.String()
		if expected != actual {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fixture("ecx/v3/l2/connections/create-error.json"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("ecx/v3/l2/connections/create-response.json"))
	}).Methods("POST")
	router.HandleFunc(ecxV3L2ConnectionsEndpoint+"/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		expectedAuthHeader := "Bearer " + expectedAccessToken
		actualAuthHeader := r.Header.Get("Authorization")
		if actualAuthHeader != expectedAuthHeader {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		expectedUUID := "c0510f9b-ca34-4bb1-8fa8-d44cc8558953"
		vars := mux.Vars(r)
		if vars["uuid"] != expectedUUID {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fixture("ecx/v3/l2/connections/{id}/notfound.json"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("ecx/v3/l2/connections/{id}/item.json"))
	}).Methods("GET")
	router.HandleFunc(ecxV3L2ConnectionsEndpoint+"/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		expectedAuthHeader := "Bearer " + expectedAccessToken
		actualAuthHeader := r.Header.Get("Authorization")
		if actualAuthHeader != expectedAuthHeader {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, fixture("oauth2/v1/token/auth-error.json"))
			return
		}

		expectedUUID := "c0510f9b-ca34-4bb1-8fa8-d44cc8558953"
		vars := mux.Vars(r)
		if vars["uuid"] != expectedUUID {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, fixture("ecx/v3/l2/connections/{id}/notfound.json"))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, fixture("ecx/v3/l2/connections/{id}/item-deleted.json"))
	}).Methods("DELETE")

	client, _ = New(BaseURL(server.URL))

	return func() {
		server.Close()
	}
}

func fixture(path string) string {
	b, err := ioutil.ReadFile("testdata/fixtures/" + path)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestAuthenticate(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	err := client.Authenticate(clientID, clientSecret, username, password)
	if err != nil {
		t.Fatal(err)
	}

	if client.accessToken != expectedAccessToken {
		t.Fatalf("Expected: %v. Actual: %v", expectedAccessToken, client.accessToken)
	}

	expectedAuthError := &AuthError{
		HTTPStatusCode: http.StatusUnauthorized,
		oAuthErrorResponse: oAuthErrorResponse{
			ErrorDomain:      "errorDomain123",
			ErrorTitle:       "errorTitle123",
			ErrorCode:        "errorCode123",
			DeveloperMessage: "developerMessage123",
			ErrorMessage:     "errorMessage123",
		},
	}
	err = client.Authenticate(clientID, clientSecret, "invalid_username", "invalid_password")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*AuthError); ok {
		if !reflect.DeepEqual(expectedAuthError, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
	}

	expectedAuthError2 := &AuthError{
		HTTPStatusCode: http.StatusInternalServerError,
		oAuthErrorResponse: oAuthErrorResponse{
			ErrorDomain:      "errorDomain123",
			ErrorTitle:       "errorTitle123",
			ErrorCode:        "errorCode123",
			DeveloperMessage: "developerMessage123",
			ErrorMessage:     "errorMessage123",
		},
	}
	err = client.Authenticate("invalid_clientID", "invalid_clientSecret", username, password)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*AuthError); ok {
		if !reflect.DeepEqual(expectedAuthError2, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedAuthError2, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedAuthError2, err)
	}
}

func TestCreateL2Connection(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	expectedAuthError := &AuthError{
		HTTPStatusCode: http.StatusUnauthorized,
		oAuthErrorResponse: oAuthErrorResponse{
			ErrorDomain:      "errorDomain123",
			ErrorTitle:       "errorTitle123",
			ErrorCode:        "errorCode123",
			DeveloperMessage: "developerMessage123",
			ErrorMessage:     "errorMessage123",
		},
	}
	_, err := client.CreateL2Connection(PostConnectionRequest{})
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*AuthError); ok {
		if !reflect.DeepEqual(expectedAuthError, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
	}

	err = client.Authenticate(clientID, clientSecret, username, password)
	if err != nil {
		t.Fatal(err)
	}

	notifications := []string{
		"user@domain.com",
	}
	input := PostConnectionRequest{
		AuthorizationKey:    StringPtr("986744318870"),
		Notifications:       &notifications,
		PrimaryName:         StringPtr("MA3-TEST"),
		PrimaryPortUUID:     StringPtr("7b5650d1-810a-10a0-66e0-30ac094f8701"),
		PrimaryVlanSTag:     IntPtr(4000),
		ProfileUUID:         StringPtr("69ee618d-be52-468d-bc99-00566f2dd2b9"),
		PurchaseOrderNumber: StringPtr("PO"),
		SellerMetroCode:     StringPtr("LD"),
		SellerRegion:        StringPtr("eu-west-2"),
		Speed:               IntPtr(50),
		SpeedUnit:           StringPtr("MB"),
	}
	actual, err := client.CreateL2Connection(input)
	if err != nil {
		t.Fatal(err)
	}

	expected := &PostConnectionResponse{
		Message:             "Connection Saved Successfully",
		PrimaryConnectionID: "552f0eab-6733-4219-804e-58fa9c77b041",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %+v. Actual: %+v", expected, actual)
	}

	expectedErrorArray := &ErrorArray{
		HTTPStatusCode: http.StatusBadRequest,
		errorResponseArray: errorResponseArray{
			{
				ErrorCode:    "IC-LAYER2-PORTS-001",
				ErrorMessage: "Port does not exist or not belong to user",
				MoreInfo:     "",
				Property:     "primaryPortUUID",
			},
		},
	}

	input.ProfileUUID = StringPtr("invalid_uuid")
	_, err = client.CreateL2Connection(input)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*ErrorArray); ok {
		if !reflect.DeepEqual(expectedErrorArray, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
	}
}

func TestReadL2Connection(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	uuid := "c0510f9b-ca34-4bb1-8fa8-d44cc8558953"
	expectedAuthError := &AuthError{
		HTTPStatusCode: http.StatusUnauthorized,
		oAuthErrorResponse: oAuthErrorResponse{
			ErrorDomain:      "errorDomain123",
			ErrorTitle:       "errorTitle123",
			ErrorCode:        "errorCode123",
			DeveloperMessage: "developerMessage123",
			ErrorMessage:     "errorMessage123",
		},
	}
	_, err := client.ReadL2Connection(uuid)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*AuthError); ok {
		if !reflect.DeepEqual(expectedAuthError, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
	}

	err = client.Authenticate(clientID, clientSecret, username, password)
	if err != nil {
		t.Fatal(err)
	}

	expected := &GetConnectionByUUIDResponse{
		BuyerOrganizationName:  "Test Networks",
		UUID:                   "c0510f9b-ca34-4bb1-8fa8-d44cc8558953",
		Name:                   "LD8_TEST",
		VlanSTag:               4000,
		PortUUID:               "7b050214-18c8-8c80-c6e0-30ac094f8be9",
		PortName:               "TEST-ECX-LD8-ECX-PRI-20912671",
		AsideEncapsulation:     "dot1q",
		MetroCode:              "LD",
		MetroDescription:       "London",
		ProviderStatus:         "DEPROVISIONED",
		Status:                 "DEPROVISIONED",
		BillingTier:            "Up to 50 MB",
		AuthorizationKey:       "986744318870",
		Speed:                  50,
		SpeedUnit:              "MB",
		RedundancyType:         "primary",
		SellerRegion:           "eu-west-2",
		SellerMetroCode:        "LD",
		SellerMetroDescription: "London",
		SellerServiceName:      "AWS Direct Connect - Redundant - East London - LD8",
		SellerServiceUUID:      "8b5cbdbd-53da-4d37-9bfc-cc988eb522f1",
		SellerOrganizationName: "EQUINIX-AWS",
		Notifications: []string{
			"user@domain.com",
		},
		PurchaseOrderNumber: "PO",
		ActionDetails: []ActionDetail{
			ActionDetail{
				ActionType:    "EQUINIX_EXECUTE_ACTION",
				OperationID:   "CONFIRM_CONNECTION",
				ActionMessage: "Please provide the following credentials so that we can confirm the connection on your behalf. You also have the option to do this operation on AWS Console",
				ActionRequiredData: []ActionRequiredDataItem{
					ActionRequiredDataItem{
						Key:               "accessKey",
						Label:             "Amazon Access Key",
						Value:             "",
						Editable:          true,
						ValidationPattern: "{RegexForString50chars}",
					},
					ActionRequiredDataItem{
						Key:               "secretKey",
						Label:             "Amazon Secret Key",
						Value:             "",
						Editable:          true,
						ValidationPattern: "{RegexForString24chars}",
					},
					ActionRequiredDataItem{
						Key:               "awsConnectionId",
						Label:             "AWS Hosted Connection Id",
						Value:             "dxcon-fgx9sq3x",
						Editable:          false,
						ValidationPattern: "",
					},
				},
			},
		},
		CreatedDate:           "2019-03-12T12:07:44.058Z",
		CreatedBy:             "user@domain.com",
		CreatedByFullName:     "User",
		CreatedByEmail:        "user@domain.com",
		LastUpdatedBy:         "user@domain.com",
		LastUpdatedDate:       "2019-03-12T12:11:27.846Z",
		LastUpdatedByFullName: "User",
		LastUpdatedByEmail:    "user@domain.com",
		DeletedBy:             "user@domain.com",
		DeletedDate:           "2019-03-12T12:11:12.268Z",
		DeletedByEmail:        "user@domain.com",
		ZSidePortName:         "dxcon-ffrlneok",
		ZSidePortUUID:         "69bb30fb-0e00-e000-6ee0-30ac094f85df",
		ZSideVlanSTag:         307,
		Remote:                false,
		Private:               false,
		Self:                  false,
	}
	actual, err := client.ReadL2Connection(expected.UUID)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected: %+v. Actual: %+v", expected, actual)
	}

	expectedErrorArray := &ErrorArray{
		HTTPStatusCode: http.StatusBadRequest,
		errorResponseArray: errorResponseArray{
			{
				ErrorCode:    "IC-LAYER2-4023",
				ErrorMessage: "Connection does not exist or do not belong to user,Please check connection Id",
				MoreInfo:     "",
				Property:     "uuid",
			},
		},
	}
	_, err = client.ReadL2Connection("invalid_uuid")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*ErrorArray); ok {
		if !reflect.DeepEqual(expectedErrorArray, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
	}
}

func TestDeleteL2Connection(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	uuid := "c0510f9b-ca34-4bb1-8fa8-d44cc8558953"
	expectedAuthError := &AuthError{
		HTTPStatusCode: http.StatusUnauthorized,
		oAuthErrorResponse: oAuthErrorResponse{
			ErrorDomain:      "errorDomain123",
			ErrorTitle:       "errorTitle123",
			ErrorCode:        "errorCode123",
			DeveloperMessage: "developerMessage123",
			ErrorMessage:     "errorMessage123",
		},
	}
	err := client.DeleteL2Connection(uuid)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*AuthError); ok {
		if !reflect.DeepEqual(expectedAuthError, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedAuthError, err)
	}

	err = client.Authenticate(clientID, clientSecret, username, password)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteL2Connection(uuid)
	if err != nil {
		t.Fatal(err)
	}

	expectedErrorArray := &ErrorArray{
		HTTPStatusCode: http.StatusBadRequest,
		errorResponseArray: errorResponseArray{
			{
				ErrorCode:    "IC-LAYER2-4023",
				ErrorMessage: "Connection does not exist or do not belong to user,Please check connection Id",
				MoreInfo:     "",
				Property:     "uuid",
			},
		},
	}
	err = client.DeleteL2Connection("invalid_uuid")
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
	if err, ok := err.(*ErrorArray); ok {
		if !reflect.DeepEqual(expectedErrorArray, err) {
			t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
		}
	} else {
		t.Fatalf("Expected: %v. Actual: %v", expectedErrorArray, err)
	}
}
