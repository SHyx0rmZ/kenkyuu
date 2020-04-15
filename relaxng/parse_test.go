package relaxng

import (
	"testing"
)

func TestParse(t *testing.T) {
	x, err := Parse("http://www.example.com/doc.xml", `<?xml version="1.0"?>
<foo><pre1:bar1 xmlns:pre1="http://www.example.com/n1"/><pre2:bar2 xmlns:pre2="http://www.example.com/n2"/></foo>`)
	y := Element{
		Name: Name{
			Namespace: "",
			Local:     "foo",
		},
		Context: Context{
			BaseURI: "http://www.example.com/doc.xml",
			Namespaces: map[string]string{
				"xml": "http://www.w3.org/XML/1998/namespace",
			},
			Default: "",
		},
		Attributes: nil,
		Children: []child{
			Element{
				Name: Name{
					Namespace: "http://www.example.com/n1",
					Local:     "bar1",
				},
				Context: Context{
					BaseURI: "http://www.example.com/doc.xml",
					Namespaces: map[string]string{
						"pre1": "http://www.example.com/n1",
						"xml":  "http://www.w3.org/XML/1998/namespace",
					},
					Default: "",
				},
				Attributes: nil,
				Children:   nil,
			},
			Element{
				Name: Name{
					Namespace: "http://www.example.com/n2",
					Local:     "bar2",
				},
				Context: Context{
					BaseURI: "http://www.example.com/doc.xml",
					Namespaces: map[string]string{
						"pre2": "http://www.example.com/n2",
						"xml":  "http://www.w3.org/XML/1998/namespace",
					},
					Default: "",
				},
				Attributes: nil,
				Children:   nil,
			},
		},
	}
}

const _31 = `<?xml version="1.0"?>
<element name="foo"
         xmlns="https://relaxng.org/ns/structure/1.0"
		 xmlns:a="http://relaxng.org/ns/annotation/1.0"
		 xmlns:ex1="http://www.example.com/n1"
         xmlns:ex2="http://www.example.com/n2">
  <a:documentation>A foo element.</a:document>
  <element name="ex1:bar1">
    <empty/>
  </element>
  <element name="ex2:bar2">
    <empty/>
  </element>
</element>`
