package xlsx

import (
	"encoding/xml"
)

// XMLBase is a XML base entity.
type XMLBase struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	Children []*XMLBase `xml:",any"`
	Chardata []byte     `xml:",chardata"`
	Comment  []byte     `xml:",comment"`
	Cdata    []byte     `xml:",cdata"`
}
