# Usage of Special Parameters #

In this example DataMover 🐼 uses special parameters in query.

Normal parameters are indicated with "@".

Special parameters are indicated with "@@".

Both syntax lookup for parameter names into VARIABLES.

But "@@" parameters are not used as standard SQL parameters, but are replaced directly in command string.

This implementation allow you to execute a statement like this:

`SELECT * FROM table1 WHERE name IN ( @@array )`

Where the @@array variable is an array and not a string or a standard SQL type.

In this test I used javascript to create the parameter (se at the job [here](./job-sqlite-javascript))