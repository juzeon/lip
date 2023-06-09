package data

import "net"

type IPLookupResult struct {
	IP         net.IP `json:"ip"`
	Source     string `json:"source"`
	Country    string `json:"country"`
	Region     string `json:"region"`
	City       string `json:"city"`
	ISP        string `json:"isp"`
	Additional string `json:"additional"`
}

var IPLookupResultTableHeader = []string{"Source", "Country", "Region", "City", "ISP", "Additional"}
