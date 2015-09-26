package main

import (
	assert "github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestDependsAA(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "a", s)
	if !CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}
}

func TestDependsAB_BA(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = DependsOn("b", "a", s)
	if !CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}
}

func TestExclusiveAB(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = AreExclusive("a", "b", s)
	if CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}
}

func TestExclusiveAB_BC(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = DependsOn("b", "c", s)
	s = AreExclusive("a", "c", s)
	if CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}
}

func TestExclusivesAB_BC_CA_DE(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = DependsOn("b", "c", s)
	s = DependsOn("c", "a", s)
	s = DependsOn("d", "e", s)
	s = AreExclusive("c", "e", s)
	if !CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}

	// Or list, array, etc.
	selected := NewSet()

	selected = Toggle(selected, "a", s)
	if Eq(selected, []string{"a", "c", "b"}) {
		t.Errorf("Toggle expected {a, c, b} got %s", selected)
		t.Log("Comment: If this is Expected to be Equal, then shouldn't be !Eq used to through t.Errorf ?")
	}

	s = DependsOn("f", "f", s)

	selected = Toggle(selected, "f", s)
	if Eq(selected, []string{"a", "c", "b", "f"}) {
		t.Errorf("Toggle expected {a, c, b, f} got %s", selected)
	}

	selected = Toggle(selected, "e", s)
	if Eq(selected, []string{"e", "f"}) {
		t.Errorf("Toggle expected {e, f} got %s", selected)
		// but this will be {e, f, d} because we have a DependsOn(d,e)
	}

	selected = Toggle(selected, "b", s)
	if Eq(selected, []string{"a", "c", "b", "f"}) {
		t.Errorf("Toggle expected {a, c, b, f} got %s", selected)
	}

	s = DependsOn("b", "g", s)
	selected = Toggle(selected, "g", s)
	selected = Toggle(selected, "b", s)
	if Eq(selected, []string{"g", "f"}) {
		t.Errorf("Toggle expected {g, f} got %s", selected)
	}
}

func TestAB_BC_Toggle(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = DependsOn("b", "c", s)
	selected := NewSet()
	selected = Toggle(selected, "c", s)
	if !Eq(selected, []string{"c"}) {
		t.Errorf("Toggle expected {c} got %s", selected)
	}
}

// Multiple dependencies and exclusions.
func TestAB_AC(t *testing.T) {
	s := NewRelationshipSet()
	s = DependsOn("a", "b", s)
	s = DependsOn("a", "c", s)
	s = AreExclusive("b", "d", s)
	s = AreExclusive("b", "e", s)
	if !CheckRelationships(s) {
		t.Error("CheckRelationships failed")
	}

	selected := NewSet()
	selected = Toggle(selected, "d", s)
	selected = Toggle(selected, "e", s)
	selected = Toggle(selected, "a", s)
	if !Eq(selected, []string{"a", "c", "b"}) {
		t.Errorf("Toggle expected {a, c, b} got %s", selected)
	}
}

func TestEq(t *testing.T) {
	aList := []string{"a", "b", "c"}
	bList := []string{"a", "c", "b"}

	if !Eq(aList, bList) {
		t.Errorf("supposes to return true when %s and %s where compared", aList, bList)
	}
}

func Eq(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)
	a = removeDuplicates(a)
	b = removeDuplicates(b)
	return reflect.DeepEqual(a, b)
}

func TestRemoveDuplicates(t *testing.T) {
	uniqueList := removeDuplicates([]string{"a", "b", "c", "c"})

	if len(uniqueList) != 3 {
		t.Error("expected to be 3 items butreturned %d", len(uniqueList))
	}
}

func removeDuplicates(input []string) (output []string) {
	if input == nil {
		return nil
	}

	if len(input) == 0 {
		return input
	}

	output = append(output, input[0])
	for index, value := range input[1:] {
		if input[index] != value {
			output = append(output, value)
		}
	}
	return
}

func TestExistsInList(t *testing.T) {
	b, idx := existsInList("a", []string{"a", "b", "c"})

	if b != true || idx != 0 {
		t.Errorf("expected exists to return true and index 0, but returned %s and %d", b, idx)
	}

	b, idx = existsInList("d", []string{"a", "b", "c"})

	if b != false || idx != -1 {
		t.Errorf("expected exists to return true and index 0, but returned %s and %d", b, idx)
	}

}

func TestDependancyList(t *testing.T) {
	assert := assert.New(t)

	dep := [][]string{
		[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"f", "f"},
	}

	depList := DependancyList(dep, "e", []string{})

	assert.Equal([]string{"e", "d"}, depList)

	dep = [][]string{
		[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"f", "f"}, []string{"d", "g"},
	}

	depList = DependancyList(dep, "e", []string{})

	assert.Equal([]string{"e", "g", "d"}, depList)

	assert.True(true)
}

func TestExclusionList(t *testing.T) {
	assert := assert.New(t)

	dep := [][]string{
		[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"f", "f"},
	}

	exc := [][]string{
		[]string{"e", "c"},
	}

	assert.Equal([]string{"c", "a", "b"}, ExclusionList(dep, exc, "e", []string{}))

	dep = [][]string{
		[]string{"a", "b", "c"}, []string{"d", "e"}, []string{"f", "f"},
	}

	exc = [][]string{
		[]string{"e", "c"},
	}

	assert.Equal([]string{"e", "d"}, DependanciesExclusions(dep, exc, "b", []string{}))

	assert.True(true)
}
