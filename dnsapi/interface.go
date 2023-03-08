package dnsapi

type APICreator func(string) (*API, error)

type API interface {
	// Create creates new Record of Zone
	Create(zoneID string, record Record) error
	// Update updates Record of Zone
	Update(zoneID string, record Record) error
	Exists(zoneID string, record Record) bool
}

type Record struct {
	// Type represents record's type
	// allowed values are: A, AAAA, CNAME, HTTPS, TXT, SRV, LOC, MX, NS, CERT, DNSKEY, DS, NAPTR, SMIMEA, SSHFP, SVCB, TLSA, URI
	Type string
	// Name is name of record, example: abc.example.com
	Name string
	// Content is content of record, example: 93.184.216.34, maximum length is 255
	Content string
	// TTL is time-to-live value of record, example: 3600(60 minutes), must be in 60 and 86400
	// (-1 points to the auto-select value according to the DNS provider if provider supports)
	TTL int
}

type Zone struct {
	// ID represents identifier of zone
	ID string
	// Name is name of zone, example: example.com (without trailing dot)
	Name string
}
