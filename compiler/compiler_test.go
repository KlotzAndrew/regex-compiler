package compiler

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/awalterschulze/gographviz"
)

func TestRe2Post(t *testing.T) {
	for _, tt := range []struct {
		reg  string
		want string
	}{
		{"a", "a"},
		{"ab", "ab."},
		{"abc", "ab.c."},
		{"(ab)c", "ab.c."},
		{"ab|c", "ab.c|"},
		{"a|b|c", "abc||"},
		{"ab*", "ab*."},
	} {

		if post := Re2post(tt.reg); string(post) != tt.want {
			t.Errorf("got %v, want %v", string(post), tt.want)
		}
	}
}

func TestPFMatch(t *testing.T) {
	for _, tt := range []struct {
		pf    []rune
		input string
		exp   bool
	}{
		{[]rune("a"), "a", true},
		{[]rune("aa."), "aa", true},
		{[]rune("aa.a."), "aaa", true},
		{[]rune("a*"), "b", false},
		{[]rune("a*"), "a", true},
		{[]rune("a*"), "aaa", true},
		{[]rune("ab|"), "b", true},
	} {
		nfa := Post2nfa(tt.pf)
		// printNFA(nfa)

		if result := MatchRe(nfa, tt.input); result != tt.exp {
			t.Errorf("got %v, want %v", result, tt.exp)
		}
	}
}

func printNFA(state *State) {
	g := gographviz.NewGraph()
	if err := g.SetName("G"); err != nil {
		panic(err)
	}
	g.SetDir(true)

	addStateToG(state, g, []string{})
	s := g.String()

	fmt.Println(s)

	file, err := os.Create(`./diGraph.dot`)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(s))

	// dot -T png diGraph.dot -o diGraph.png
}

func addStateToG(s *State, g *gographviz.Graph, routeID []string) {
	if s == nil {
		return
	}

	if len(routeID) > 5 {
		return
	}

	newRouteID := append(routeID, string(s.c+100)+string(s.s))
	if strings.Join(routeID, "") == strings.Join(newRouteID, "") {
		return
	}
	if len(routeID) > 0 {
		defer g.AddEdge(strings.Join(routeID, ""), strings.Join(newRouteID, ""), true, nil)
	}

	addStateToG(s.out, g, newRouteID)
	addStateToG(s.out1, g, newRouteID)

	nodeAttrSwitch := make(map[string]string)
	nodeAttrSwitch["color"] = "blue"
	nodeAttrSwitch["fillcolor"] = "blue"
	nodeAttrSwitch["shape"] = "doublecircle"
	nodeAttrSwitch["label"] = string(s.s) + " "

	nodeAttrDefault := make(map[string]string)
	nodeAttrDefault["label"] = string(s.s) + " "

	if s.c == Split {
		g.AddNode("G", strings.Join(newRouteID, ""), nodeAttrSwitch)
	} else {
		g.AddNode("G", strings.Join(newRouteID, ""), nodeAttrDefault)
	}
}
