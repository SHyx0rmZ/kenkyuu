package keepass

import (
	"encoding/xml"
)

type Bool bool

func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s := "False"
	if b {
		s = "True"
	}
	return e.EncodeElement(s, start)
}

func (b Bool) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	a := xml.Attr{
		Name:  name,
		Value: "False",
	}
	if b {
		a.Value = "True"
	}
	return a, nil
}
