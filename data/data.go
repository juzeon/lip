package data

type IPLookupResult struct {
	Source  string
	Country string
	State   string
	City    string
	ISP     string
}

var IPLookupResultTableHeader = []string{"Source", "Country", "State", "City", "ISP"}
