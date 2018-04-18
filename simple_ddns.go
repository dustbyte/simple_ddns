package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func getIp() (string, error) {
	resp, err := http.Get("http://ifconfig.kreog.net/")
	if err != nil {
		fmt.Printf("Couldn't retrieve self ip address: %v\n", err)
	}

	defer resp.Body.Close()
	ipAddr, err := ioutil.ReadAll(resp.Body)

	return strings.TrimSpace(string(ipAddr[:])), err
}

func main() {
	var token string
	var assigned bool

	if token, assigned = os.LookupEnv("DNSIMPLE_TOKEN"); !assigned {
		fmt.Println("Cannot find token")
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Printf("usage: %s domain_name\n", os.Args[0])
		os.Exit(1)
	}

	domain_name := os.Args[1]

	ipAddr, err := getIp()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("IP address: %s\n", ipAddr)

	client := dnsimple.NewClient(dnsimple.NewOauthTokenCredentials(token))

	whoamiResponse, err := client.Identity.Whoami()
	if err != nil {
		fmt.Printf("Whoami() returned error: %v\n", err)
		os.Exit(1)
	}
	accountId := strconv.Itoa(whoamiResponse.Data.Account.ID)

	options := dnsimple.ZoneRecordListOptions{Type: "A"}
	resp, err := client.Zones.ListRecords(accountId, domain_name, &options)
	if err != nil {
		fmt.Printf("Domain name %s could not be found\n", domain_name)
		os.Exit(1)
	}

	record := resp.Data[0]
	if record.Content != ipAddr {
		attributes := dnsimple.ZoneRecord{Content: ipAddr}
		_, err := client.Zones.UpdateRecord(accountId, domain_name, record.ID, attributes)
		if err != nil {
			fmt.Printf("Couldn't update the A record for %s\n", domain_name)
			os.Exit(1)
		}

		fmt.Printf("Record A for %s successfully updated.\n", domain_name)
	}
}
