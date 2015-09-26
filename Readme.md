A library that handles the must-include and must-omit relations in a smart way
Ever wanted to maintain a list of elements, where when you add one element,
you have to add its dependancies, and remove its exclusivity obstacles.

Also, maintain a relationship matrix, of who depends on whom and who wants exclusivity
that too without overlapping.

A depends on B, and B depends on D, so when A is included both B and C should also be included
G exclusive to A, so if G is added, A,B and C should all be removed from the list.

But is somewhere there is a relation G depends on B, then this whole relationship matrix fails

Coz G cannot have A and need to have B but B would need A.
