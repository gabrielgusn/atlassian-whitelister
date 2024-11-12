package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

type IpListItem struct {
	Network   string   `json:"network"`
	Mask_len  int      `json:"mask_len"`
	Cidr      string   `json:"cidr"`
	Mask      string   `json:"mask"`
	Region    []string `json:"region"`
	Product   []string `json:"product"`
	Direction []string `json:"direction"`
	Perimeter string   `json:"perimeter"`
}

type IpList struct {
	CreationDate string       `json:"creationDate"`
	SyncToken    int64        `json:"syncToken"`
	Items        []IpListItem `json:"items"`
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

func main() {
	c, err := identity.NewIdentityClientWithConfigurationProvider(common.DefaultConfigProvider())

	if err != nil {
		log.Println("Error while creating OCI Identity:", err)
		return
	}

	tenancyId, err := common.DefaultConfigProvider().TenancyOCID()
	if err != nil {
		log.Println("Error while gettint Tenancy OCID:", err)
		return
	}

	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: &tenancyId,
	}

	r, err := c.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		log.Println("Error while requesting available ADs:", err)
		return
	}
	fmt.Printf("List of available domains: %v", r.Items)

	fmt.Println(tenancyId)
	// Probably https://ip-ranges.atlassian.com/
	atlassian_ipranges_url := os.Getenv("ATLASSIAN_IPRANGES_URL")

	if atlassian_ipranges_url == "" {
		log.Fatal("Error: ATLASSIAN_IPRANGES_URL env var is nil or not defined.")
	}

	resp, err := client.Get(atlassian_ipranges_url)

	if err != nil {
		log.Fatalf("ERROR: Failed to create HTTP Get request to %s: %s\n", atlassian_ipranges_url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("REQUEST ERROR: The request to %s returned %v", atlassian_ipranges_url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("ERROR: Error while reading the request body: %s\n", err)
	}

	var ipList IpList

	err = json.Unmarshal(body, &ipList)

	if err != nil {
		log.Fatalf("JSON ERROR: Failed to unmarshal request body to structured data: %v", err)
	}

	// fmt.Print(ipList)

	// ipList.listBitbucketIps()
}

func (i *IpList) listBitbucketIps() (data []string, err any) {
	if i == nil {
		return nil, "Trying to list IPs from an empty object. Returning nil."
	} else if i.Items == nil {
		return nil, "Trying to list IPs from an empty array. Returning nil."
	}

	// var bitbucketIps []string

	//TODO: try to loop through the item list with more than 1 thread
	for pos, item := range i.Items {
		if contains(item.Product, "bitbucket") {
			fmt.Println(pos, item)
		}
	}

	return nil, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
