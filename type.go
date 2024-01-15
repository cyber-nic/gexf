package gexf

// Type is an type for values Nodes can hold.
type Type string

// All the recognized GEXTTypes.
const (
	Long       Type = "long"
	Double     Type = "double"
	Float      Type = "float"
	Boolean    Type = "boolean"
	ListString Type = "liststring"
	String     Type = "string"
	AnyURI     Type = "anyURI"
)
