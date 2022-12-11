package dns

type RecordValue struct {
	Data      []string
	View      string
	Weight    uint16
	Continent string
	Country   string
}

type Record struct {
	// The zone name.
	Name string
	// The zone class.
	Type Type
	// The zone TTL.
	TTL uint32
	// The zone records.
	Value []RecordValue
}
