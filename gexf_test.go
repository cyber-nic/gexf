package gexf

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestNewGraph tests that NewGraph() returns a non-nil graph.
func TestNewGraph(t *testing.T) {
	g := NewGraph()

	if g == nil {
		t.Error("NewGraph() returned nil")
	}
}

// TestAddNode tests that AddNode() adds a node to the graph.
func TestAddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("1", "node 1", []AttrValue{})

	if len(g.Nodes) != 1 {
		t.Error("AddNode() failed to add node to graph")
	}
}

// TestSetNodeAttrs tests that SetNodeAttrs() adds an attribute to the graph.
func TestSetNodeAttrs(t *testing.T) {
	g := NewGraph()
	g.SetNodeAttrs([]Attr{
		{
			Title:   "attr1",
			Type:    String,
			Default: "attr1",
		},
		{
			Title:   "attr2",
			Type:    Float,
			Default: 1,
		},
	})

	if len(g.NodeAttrs.Attrs) != 2 {
		t.Error("SetNodeAttrs() failed to add attr to graph")
	}
}

// TestAddNode tests that AddNode() adds a node to the graph.
func TestAddNodeWithAttr(t *testing.T) {
	g := NewGraph()
	g.AddNode("1", "node 1", []AttrValue{
		{
			Title: "attr1",
			Value: "attr1",
		},
	})

	if len(g.Nodes) != 1 {
		t.Error("AddNode() failed to add node with attr to graph")
	}
}

// TestSetEdgeAttrs tests that SetEdgeAttrs() adds an attribute to the graph.
func TestSetEdgeAttrs(t *testing.T) {
	g := NewGraph()
	g.SetEdgeAttrs([]Attr{
		{
			Title:   "attr1",
			Type:    String,
			Default: "attr1",
		},
		{
			Title:   "attr2",
			Type:    Float,
			Default: 1,
		},
	})

	if len(g.EdgeAttrs.Attrs) != 2 {
		t.Error("SetEdgeAttrs() failed to add attr to graph")
	}
}

func TestEncode(t *testing.T) {
	g := NewGraph()
	// define node attributes
	g.SetNodeAttrs([]Attr{
		{
			Title:   "a0",
			Type:    String,
			Default: "foo",
		},
		{
			Title:   "a1",
			Type:    Float,
			Default: 1,
		},
	})
	// define edge attributes
	g.SetEdgeAttrs([]Attr{
		{
			Title:   "a2",
			Type:    String,
			Default: "bar",
		},
	})
	// add nodes
	g.AddNode("1", "node 1", []AttrValue{
		{
			Title: "a0",
			Value: "BAR",
		},
	})
	g.AddNode("2", "node 2", []AttrValue{
		{
			Title: "a1",
			Value: 2,
		},
	})
	// add edge
	g.AddEdge("1", "2", []AttrValue{
		{
			Title: "a2",
			Value: "FOO",
		},
	})

	var w bytes.Buffer

	err := Encode(&w, g)
	if err != nil {
		t.Error(err)
	}

	ee := `<gexf xmlns="http://www.gexf.net/1.2draft" version="1.2">
    <meta lastmodifieddate="2024-01-15">
        <creator>webscale!</creator>
        <description>so fast!</description>
    </meta>
    <graph>
        <mode>static</mode>
        <defaultedgetype>directed</defaultedgetype>
        <attributes class="node">
            <attribute id="0" title="a0" type="string">
                <default>foo</default>
            </attribute>
            <attribute id="1" title="a1" type="float">
                <default>1</default>
            </attribute>
        </attributes>
        <attributes class="edge">
            <attribute id="2" title="a2" type="string">
                <default>bar</default>
            </attribute>
        </attributes>
        <nodes>
            <node id="1" label="node 1">
                <attvalues>
                    <attvalue for="0" value="BAR"></attvalue>
                </attvalues>
            </node>
            <node id="2" label="node 2">
                <attvalues>
                    <attvalue for="1" value="2"></attvalue>
                </attvalues>
            </node>
        </nodes>
        <edges>
            <edge id="0" source="1" target="2">
                <attvalues>
                    <attvalue for="2" value="FOO"></attvalue>
                </attvalues>
            </edge>
        </edges>
    </graph>
</gexf>`

	// fmt.Println(ee)
	// fmt.Println(w.String())

	diff := cmp.Diff(w.String(), ee)
	if diff != "" {
		t.Errorf("incorrect encoded graph %v", diff)
	}
}
