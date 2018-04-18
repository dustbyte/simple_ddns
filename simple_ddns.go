package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/mota/klash"
)

const (
	defaultTtl = 60
)

type Args struct {
	Token string `klash-alias:"t" klash-help:"DNSimple API token"`
	TTL   int    `klash-help:"TTL of the record in seconds"`
}

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
	args := Args{Token: os.Getenv("DNSIMPLE_TOKEN"), TTL: defaultTtl}
	leftover := klash.Parse(&args, "DynDNS for mere mortals")

	if len(leftover) == 0 {
		fmt.Printf("usage: %s domain_name\n", os.Args[0])
		os.Exit(1)
	}

	domain_name := leftover[0]

	if args.Token == "" {
		fmt.Println("A token must be provided")
		os.Exit(1)
	}

	ipAddr, err := getIp()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("IP address: %s\n", ipAddr)

	client := dnsimple.NewClient(dnsimple.NewOauthTokenCredentials(args.Token))

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

	if len(resp.Data) == 0 {
		fmt.Printf("Record A doesn't exist for domain %s, creating...\n", domain_name)
		record := dnsimple.ZoneRecord{Type: "A", Name: "", TTL: args.TTL, Content: ipAddr}

		_, err := client.Zones.CreateRecord(accountId, domain_name, record)
		if err != nil {
			fmt.Printf("Couldn't create the A record for domain %s\n", domain_name)
			os.Exit(1)
		}
		fmt.Printf("Record for domain name %s created.\n", domain_name)
	} else {
		record := resp.Data[0]
		if record.Content != ipAddr {
			attributes := dnsimple.ZoneRecord{Content: ipAddr}
			_, err := client.Zones.UpdateRecord(accountId, domain_name, record.ID, attributes)
			if err != nil {
				fmt.Printf("Couldn't update the A record for %s\n", domain_name)
				os.Exit(1)
			}

			fmt.Printf("Record A for %s successfully updated.\n", domain_name)
		} else {
			fmt.Println("No change to perform.")
		}
	}
}
