package main

import (
	"fmt"
)

var (
	dependsID   = "depends"
	exclusiveID = "exclusive"
)

func NewRelationshipSet() map[string][][]string {
	return map[string][][]string{
		dependsID:   [][]string{},
		exclusiveID: [][]string{},
	}
}

func existsInList(a string, aList []string) (bool, int) {
	for idx, item := range aList {
		if item == a {
			return true, idx
		}
	}

	return false, -1
}

func existsInIntList(a int, aList []int) (bool, int) {
	for idx, item := range aList {
		if item == a {
			return true, idx
		}
	}

	return false, -1
}

func DependsOn(A, B string, s map[string][][]string) map[string][][]string {
	depends := s[dependsID]

	nDepends := [][]string{}

	if len(depends) == 0 {
		list := []string{A, B}
		s[dependsID] = [][]string{list}
		return s
	}

	aOccurance := 0
	aOccuranceIdx := []int{}

	for idx, list := range depends {
		if ok, _ := existsInList(A, list); ok {
			aOccurance = aOccurance + 1
			aOccuranceIdx = append(aOccuranceIdx, idx)
		}
	}

	bOccurance := 0
	bOccuranceIdx := []int{}

	for idx, list := range depends {
		if ok, _ := existsInList(B, list); ok {
			bOccurance = bOccurance + 1
			bOccuranceIdx = append(bOccuranceIdx, idx)
		}
	}

	// the list of index returned; means all of those elements inside each list are dependant on each other.
	// So merge those items in indexes; and
	// append all indexes that doesn't have these elements
	totalOccuranceIdx := aOccuranceIdx
	for _, item := range bOccuranceIdx {
		totalOccuranceIdx = append(totalOccuranceIdx, item)
	}

	totalOccuranceIdx = RemoveIntDuplicates(totalOccuranceIdx)

	theUnaffectedListOfList := [][]string{}
	theAffectedList := []string{}

	for idx, list := range depends {
		if ok, _ := existsInIntList(idx, totalOccuranceIdx); ok {
			for _, item := range list {
				// if the index with item(A/B) in list is found, append each item in the list to affected list
				theAffectedList = append(theAffectedList, item)
			}
		} else {
			// if either of the A/B item is not in list; append it to the list of unaffected list
			theUnaffectedListOfList = append(theUnaffectedListOfList, list)
		}
	}

	theAffectedList = RemoveDuplicates(theAffectedList)
	// nDepends is still empty
	nDepends = theUnaffectedListOfList
	if len(theAffectedList) != 0 {
		nDepends = append(nDepends, theAffectedList)
	}

	depends = nDepends

	if aOccurance > 0 && bOccurance > 0 {
		s[dependsID] = nDepends
		return s
	}

	if aOccurance == 0 && bOccurance == 0 {
		nDepends = append(nDepends, []string{A, B})
		s[dependsID] = nDepends
		return s
	}

	// setting back nDepends to nil
	nDepends = [][]string{}

	if aOccurance > 0 {
		for _, list := range depends {
			if ok, _ := existsInList(A, list); ok {
				list = append(list, B)
			}
			nDepends = append(nDepends, list)
		}
	}

	if bOccurance > 0 {
		for _, list := range depends {
			if ok, _ := existsInList(B, list); ok {
				list = append(list, A)
			}
			nDepends = append(nDepends, list)
		}
	}

	s[dependsID] = nDepends
	return s
}

func AreExclusive(A, B string, s map[string][][]string) map[string][][]string {
	nExclusive := s[exclusiveID]

	list := []string{A, B}
	nExclusive = append(nExclusive, list)
	s[exclusiveID] = nExclusive
	return s
}

func NewSet() []string {
	return []string{}
}

func DependancyList(D [][]string, A string, traversedList []string) []string {
	traversedList = append(traversedList, A)

	for _, list := range D {
		if ok, _ := existsInList(A, list); ok {
			// if exists
			for _, item := range list {
				if ok, _ := existsInList(item, traversedList); !ok {
					nTraversedList := DependancyList(D, item, traversedList)

					for _, ntlItem := range nTraversedList {
						traversedList = append(traversedList, ntlItem)
					}
				}
			}
		}
	}

	traversedList = RemoveDuplicates(traversedList)

	return traversedList
}

func DependanciesExclusions(D [][]string, E [][]string, A string, exclusionList []string) []string {
	// removing exclusions of the dependancies
	exclusionDependancy := [][]string{}
	for _, dep := range DependancyList(D, A, []string{}) {
		list := ExclusionList(D, E, dep, []string{})
		exclusionDependancy = append(exclusionDependancy, list)
	}

	eNestedExclusivity := []string{}

	for _, list := range exclusionDependancy {
		for _, item := range list {
			eNestedExclusivity = append(eNestedExclusivity, item)
		}
	}

	eNestedExclusivity = RemoveDuplicates(eNestedExclusivity)

	return eNestedExclusivity
}

func ExclusionList(D [][]string, E [][]string, A string, exclusionList []string) []string {
	exclusionDependancy := [][]string{}

	// removing dependancies of the exclusions
	for _, list := range E {
		if list[0] == A {
			exclusionList = append(exclusionList, list[1])

			exclusionDependancy = append(exclusionDependancy, DependancyList(D, list[1], []string{}))
		}

		if list[1] == A {
			exclusionList = append(exclusionList, list[0])

			exclusionDependancy = append(exclusionDependancy, DependancyList(D, list[0], []string{}))
		}
	}

	for _, list := range exclusionDependancy {
		for _, item := range list {
			exclusionList = append(exclusionList, item)
		}
	}

	return RemoveDuplicates(exclusionList)
}

func Toggle(selected []string, A string, s map[string][][]string) []string {
	nSelected := []string{}

	// if the element exists in selected; remove it
	// which includes removing any dependancy in dependants list; which will fire a chain of dependancy ejections
	if ok, _ := existsInList(A, selected); !ok {
		// the element has to be added
		eNestedDependancy := DependancyList(s[dependsID], A, []string{})
		eNestedExclusivity := ExclusionList(s[dependsID], s[exclusiveID], A, []string{})

		eNestedDependancyExclusion := DependanciesExclusions(s[dependsID], s[exclusiveID], A, []string{})

		for _, eNDE := range eNestedDependancyExclusion {
			eNestedExclusivity = append(eNestedExclusivity, eNDE)
		}

		eNestedExclusivity = RemoveDuplicates(eNestedExclusivity)

		nSelected = append(nSelected, A)

		for _, s := range selected {
			if ok, _ := existsInList(s, eNestedExclusivity); !ok {
				nSelected = append(nSelected, s)
			}
		}

		for _, d := range eNestedDependancy {
			nSelected = append(nSelected, d)
		}
	} else {
		// remove the element and dependancies from the list
		eNestedDependancy := DependancyList(s[dependsID], A, []string{})

		for _, s := range selected {
			if ok, _ := existsInList(s, eNestedDependancy); !ok {
				nSelected = append(nSelected, s)
			}
		}
	}

	nSelected = RemoveDuplicates(nSelected)

	fmt.Println("\nstate:", s, "\ntoggle", A, "prev:", selected, "now:", nSelected)

	return nSelected
}

func CheckRelationships(s map[string][][]string) bool {
	for _, vlist := range s[dependsID] {
		for _, evList := range s[exclusiveID] {
			counter := 0
			for _, v := range vlist {
				for _, ev := range evList {
					if v == ev {
						counter = counter + 1
					}
				}
				if counter > 1 {
					return false
				}
			}
		}
	}

	return true
}

func RemoveDuplicates(this []string) []string {
	length := len(this) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if this[i] == this[j] {
				this[j] = this[length]
				this = this[0:length]
				length--
				j--
			}
		}
	}
	return this
}

func RemoveIntDuplicates(this []int) []int {
	length := len(this) - 1
	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if this[i] == this[j] {
				this[j] = this[length]
				this = this[0:length]
				length--
				j--
			}
		}
	}
	return this
}
