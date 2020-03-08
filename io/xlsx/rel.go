package xlsx

// Rel is a rel.
type Rel struct {
	ID     string `xml:"Id,attr"`
	Target string `xml:"Target,attr"`
}

// Rels is the rels.
type Rels struct {
	Rels []Rel `xml:"Relationship"`
}
