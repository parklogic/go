package linode

import (
	"fmt"
)

const (
	PeerASN   = 65000
	ServerASN = 65001

	PeerAddressPrefix = "2600:3c0f"

	PrimaryIPCommunity   = 65000<<16 + 1
	SecondaryIPCommunity = 65000<<16 + 2
)

var (
	DatacenterIDFromRegion = map[string]int{
		"ap-northeast": 11,
		"ap-south":     9,
		"ap-southeast": 16,
		"ap-west":      14,
		"au-mel":       45,
		"br-gru":       21,
		"ca-entral":    15,
		"de-fra-2":     47,
		"es-mad":       24,
		"eu-central":   10,
		"eu-west":      7,
		"fr-par":       19,
		"gb-lon":       44,
		"id-cgk":       29,
		"in-bom-2":     46,
		"in-maa":       25,
		"it-mil":       27,
		"jp-osa":       26,
		"jp-tyo-3":     49,
		"nl-ams":       22,
		"se-sto":       23,
		"sg-sin-2":     48,
		"us-central":   2,
		"us-east":      6,
		"us-iad":       17,
		"us-lax":       30,
		"us-mia":       28,
		"us-ord":       18,
		"us-sea":       20,
		"us-southeast": 4,
		"us-west":      3,
	}

	RegionsFromDatacenterID = map[int]string{
		10: "eu-central",
		11: "ap-northeast",
		14: "ap-west",
		15: "ca-entral",
		16: "ap-southeast",
		17: "us-iad",
		18: "us-ord",
		19: "fr-par",
		20: "us-sea",
		21: "br-gru",
		22: "nl-ams",
		23: "se-sto",
		24: "es-mad",
		25: "in-maa",
		26: "jp-osa",
		27: "it-mil",
		28: "us-mia",
		29: "id-cgk",
		2:  "us-central",
		30: "us-lax",
		3:  "us-west",
		44: "gb-lon",
		45: "au-mel",
		46: "in-bom-2",
		47: "de-fra-2",
		48: "sg-sin-2",
		49: "jp-tyo-3",
		4:  "us-southeast",
		6:  "us-east",
		7:  "eu-west",
		9:  "ap-south",
	}
)

func GetBGPPeerAddresses(region string) ([]string, error) {
	dcID, ok := DatacenterIDFromRegion[region]
	if !ok {
		return nil, fmt.Errorf("%w: datacenter ID not found for region %q", ErrInvalidRegion, region)
	}

	res := make([]string, 4)
	for i := 0; i < 4; i++ {
		res[i] = fmt.Sprintf("%s:%d:34::%d", PeerAddressPrefix, dcID, i+1)
	}

	return res, nil
}
