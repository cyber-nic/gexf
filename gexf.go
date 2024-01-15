// Package gexf provides GEXF formatting and encoding
package gexf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"time"
)

// MarshalXML marshals a GEXF graph. This custom marshaler is needed to support duplicate `attributes` elements for node and edge.
func (g Graph) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "graph"
	e.EncodeToken(start)

	if err := e.EncodeElement(g.Mode, xml.StartElement{Name: xml.Name{Local: "mode"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(g.EdgeType, xml.StartElement{Name: xml.Name{Local: "defaultedgetype"}}); err != nil {
		return err
	}

	// handle NodeAttributes
	nodeAttrsStartElement := xml.StartElement{Name: xml.Name{Local: "attributes"}}
	e.EncodeElement(g.NodeAttrs, nodeAttrsStartElement)

	// handle EdgeAttributes
	edgeAttrsStartElement := xml.StartElement{Name: xml.Name{Local: "attributes"}}
	e.EncodeElement(g.EdgeAttrs, edgeAttrsStartElement)

	// Start the <nodes> element
	nodesStartElement := xml.StartElement{Name: xml.Name{Local: "nodes"}}
	if err := e.EncodeToken(nodesStartElement); err != nil {
		return err
	}

	// Encode each Node within the <nodes> element
	for _, node := range g.Nodes {
		nodeStartElement := xml.StartElement{Name: xml.Name{Local: "node"}}
		if err := e.EncodeElement(node, nodeStartElement); err != nil {
			return err
		}
	}

	// End the <nodes> element
	if err := e.EncodeToken(xml.EndElement{Name: nodesStartElement.Name}); err != nil {
		return err
	}

	// Start the <edges> element
	edgesStartElement := xml.StartElement{Name: xml.Name{Local: "edges"}}
	if err := e.EncodeToken(edgesStartElement); err != nil {
		return err
	}

	// Encode each Edge within the <edges> element
	for _, edge := range g.Edges {
		edgeStartElement := xml.StartElement{Name: xml.Name{Local: "edge"}}
		if err := e.EncodeElement(edge, edgeStartElement); err != nil {
			return err
		}
	}

	// End the <edges> element
	if err := e.EncodeToken(xml.EndElement{Name: edgesStartElement.Name}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}

	return nil
}

// Encode encodes a graph to GEXF.
func Encode(w io.Writer, g *Graph) error {
	gx := gexf{
		Namespace: "http://www.gexf.net/1.2draft",
		Version:   "1.2",
		Meta: &meta{
			LastModified: time.Now().Format("2006-01-02"),
			Creator:      "webscale!",
			Desc:         "so fast!",
		},
		Graph: g,
	}

	data, err := xml.MarshalIndent(gx, "", "    ")
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	_, err = io.Copy(w, buf)

	return err
	// return xml.NewEncoder(w).Encode(gx)
}

// Attr is an attribute for a node or edge.
type Attr struct {
	Title   string
	Type    Type
	Default interface{}
}

// AttrValue is a value for an attribute.
type AttrValue struct {
	Title string
	Value interface{}
}

// Graph is a GEXF graph.
type Graph struct {
	XMLName xml.Name `xml:"graph"`

	Mode     string `xml:"mode,attr,omitempty"`
	EdgeType string `xml:"defaultedgetype,attr"`

	Nodes     []node `xml:"nodes>node"`
	NodeAttrs *attributes

	Edges     []edge `xml:"edges>edge"`
	EdgeAttrs *attributes

	attrTitleToID map[string]string
	featureToID   map[interface{}]string
}

// NewGraph returns a new Graph.
func NewGraph() *Graph {
	return &Graph{
		Mode:          "static",
		EdgeType:      "directed",
		attrTitleToID: make(map[string]string),
		featureToID:   make(map[interface{}]string),
	}
}

// GetNodeAttrs returns the attributes for nodes.
func (g *Graph) GetNodeAttrs() []Attr {
	var attrs []Attr
	for _, a := range g.NodeAttrs.Attrs {
		attrs = append(attrs, Attr{
			Title:   a.Title,
			Type:    Type(a.Type),
			Default: a.Default,
		})
	}
	return attrs
}

// SetNodeAttrs sets the attributes for nodes.
func (g *Graph) SetNodeAttrs(attrs []Attr) error {
	g.NodeAttrs = &attributes{
		Class: "node",
	}
	for _, a := range attrs {
		if _, ok := g.attrTitleToID[a.Title]; ok {
			return fmt.Errorf("attr '%s' defined multiple times", a.Title)
		}

		id := len(g.attrTitleToID)
		attr := attribute{
			ID:      strconv.Itoa(id),
			Title:   a.Title,
			Type:    string(a.Type),
			Default: a.Default,
		}

		g.NodeAttrs.Attrs = append(g.NodeAttrs.Attrs, attr)
		g.attrTitleToID[attr.Title] = attr.ID
	}
	return nil
}

// AddNode adds a node to the graph.
func (g *Graph) AddNode(id, label string, attr []AttrValue) {
	n := node{
		ID:    id,
		Label: label,
	}

	var values []attrValue
	for _, a := range attr {
		av := attrValue{
			For:   g.attrTitleToID[a.Title],
			Value: a.Value,
		}
		values = append(values, av)
	}

	if len(values) > 0 {
		n.Attr = &values
	}

	g.Nodes = append(g.Nodes, n)
}

// AddEdge adds an edge to the graph.
func (g *Graph) AddEdge(from, to string, attr []AttrValue) {
	e := edge{
		ID:     strconv.Itoa(len(g.Edges)),
		Source: from,
		Target: to,
	}

	var values []attrValue
	for _, a := range attr {
		av := attrValue{
			For:   g.attrTitleToID[a.Title],
			Value: a.Value,
		}
		values = append(values, av)
	}

	if len(values) > 0 {
		e.Attr = &values
	}

	g.Edges = append(g.Edges, e)
}

// SetEdgeAttrs sets the attributes for edges.
func (g *Graph) SetEdgeAttrs(attrs []Attr) error {
	g.EdgeAttrs = &attributes{
		Class: "edge",
	}
	for _, a := range attrs {
		if _, ok := g.attrTitleToID[a.Title]; ok {
			return fmt.Errorf("edge '%s' defined multiple times", a.Title)
		}

		id := len(g.attrTitleToID)
		attr := attribute{
			ID:      strconv.Itoa(id),
			Title:   a.Title,
			Type:    string(a.Type),
			Default: a.Default,
		}

		g.EdgeAttrs.Attrs = append(g.EdgeAttrs.Attrs, attr)
		g.attrTitleToID[attr.Title] = attr.ID
	}
	return nil
}

// GetID returns the ID for a feature.
func (g *Graph) GetID(feature interface{}) string {
	if id, ok := g.featureToID[feature]; ok {
		return id
	}

	newID := strconv.Itoa(len(g.featureToID))
	g.featureToID[feature] = newID
	return newID
}
