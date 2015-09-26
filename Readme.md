// Tyba's code test.
//
// Let's say we have a set of options which the user can select. Options
// can related between them in two ways: one can depend on another, and two
// options can be mutually exclusive. That means that these equalities must
// always hold true (note: this is not code, those are logical equations):
//
// "A depends on B" or "for A to be selected, B needs to be selected"
// DependsOn(A, B) =>
//     if isSelected(A) then isSelected(B)
//
// "A and B are exclusive" or "B and A are exclusive" or "for A to be selected,
// B needs to be unselected; and for B to be selected, A needs to be unselected"
// AreExclusive(A, B) => AreExclusive(B, A) =>
//     if isSelected(A) then !isSelected(B)
//     if isSelected(B) then !isSelected(A)
//
// We say that a set of relations are coherent if the laws above are valid
// for that set. For example, this set of relations is coherent:
//
// DependsOn(A, B) -- "A depends on B"
// DependsOn(B, C)
// AreExclusive(C, D) -- "C and D are exclusive"
//
// And those sets are _not_ coherent:
//
// DependsOn(A, B)
// AreExclusive(A, B) -- A depends on B; it's a contradiction that they are
//                     -- exclusive.
//                     -- If B is selected, then A would need to be selected,
//                     -- but that's impossible because, by the exclusion rule,
//                     -- A is supposed to be unselected when B is selected.
//
// DependsOn(A, B)
// DependsOn(B, C)
// AreExclusive(A, C) -- The dependency relation is transitive; it's easy to
//                     -- see, from the rules above, that if A depends on B, and
//                     -- B depends on C, then A also depends on C. So this is a
//                     -- contradiction for the same reason as the previous
//                     -- case.
//
// So this is the problem:
//
// a. Data structure for expressing these relationships between options, ie. for
// defining a relationships set.
// You also need to define three functions:
//
// NewRelationshipSet(): Returns an empty relationship set.
// DependsOn(A, B, s): Returns a relationship set in which A depends on B. s is
//                      the previous relationship set.
// AreExclusive(A, B, s): Returns a relationship set in which A and B are
//                         exclusive. s is the previous relationship set.
//
// b. Algorithm that checks that an instance of that data structure is coherent,
// that is, that no option can depend, directly or indirectly, on another option
// and also be mutually exclusive with it.
//
// CheckRelationships(s): Returns true if s is a coherent relationship set,
//                         false otherwise.
//
// c. Algorithm that, given the relationships between options, an option,
// and a collection of selected options coherent with the relationships,
// adds the option to a collection of selected options, or removes it from the
// collection if it is already there, selecting and deselecting options
// automatically based on dependencies and exclusions.
//
// toggle(already_selected, o, s): Returns a collection (list, set, array,
// whatever) of selected options.
// - already_selected: a collection of selected options.
// - o: the option we want to add to or remove from already_selected.
// - s: the relationship set that the returned set must be coherent with.
//
// The practical issue from which this problem arose is this:
//
// You have an application for renting cars. One of the steps in the process is
// choosing the the color, extras, etc. you want. Some options required some
// other option to be included. For example, you can't have the GPS system if
// you don't also have a USB plug. (But you can have the USB without having the
// GPS system.) So, when the user clicks on the GPS checkbox, the USB checkbox
// is automatically checked, and if the USB checkbox is unchecked, the GPS
// checkbox gets unchecked too. Now, say you have another system that uses the
// USB plug, for example a music device or whatever. The USB plug is unique, so
// you can't have both the GPS and the music thing. When the user checks the
// USB checkbox, the music checkbox is disabled, and the other way around. (This
// case also illustrates how the relationships set {DependsOn(B, A),
// DependsOn(C, A), AreExclusive(B, C)} is perfectly coherent.)
//
// The algorithm for when a checkbox is selected is asked to you in section c,
// based on the data structures you define in section a. In section b you should
// provide an algorithm that 'tests' that sections a and c give a good solution.
//
// Test cases.
