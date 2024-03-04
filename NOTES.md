# About the grammar.

## Caveats

There are no function calls after an assignment, i.e. there is no `variable = function("something")`
in this programming language you can't assign a variable as the result of a function
call, to do that I just need to add a function call construct and parse it, and
that's pretty much it.
