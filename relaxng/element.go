package relaxng

type Name struct {
	Namespace string
	Local     string
}

type Context struct {
	BaseURI    string
	Namespaces map[string]string
	Default    string
}

type Attribute struct {
	Name  string
	Value string
}

type Element struct {
	Name
	Context
	Attributes []Attribute
	Children   []child
}

type Text string

type child interface{ child() }

func (Element) child() {}
func (Text) child()    {}

/*
4.2. Whitespace

For each element other than value and param, each child that is a string containing only whitespace characters is removed.

Leading and trailing whitespace characters are removed from the value of each name, type and combine attribute and from the content of each name element.
*/
