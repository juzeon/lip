package data

type IPLookupResult struct {
	Source  string
	Country string
	Region  string
	City    string
	ISP     string
}

var IPLookupResultTableHeader = []string{"Source", "Country", "Region", "City", "ISP"}
