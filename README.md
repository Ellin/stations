# pathfinder
A cli path-finding algorithm, that find the most efficient paths to move trains from one destination to another, with minimal turns.

## Setup and Usage <br>

clone the repository with https: <br>
`git clone https://gitea.kood.tech/aliadaahirmohamed/pathfinder.git`

or with a ssh: <br>
`git clone git@gitea.kood.tech:aliadaahirmohamed/pathfinder.git`

navigate to the project with:
`cd pathfinder`

Then run for Edmond Karp algorithm <br>
`go run . [INPUTFILEPATH] [STARTSTATION] [ENDSTATION] [TRAINNUMBER]` 


or  for Dinic algorithm <br>
`go run . -alg Dinic [INPUTFILEPATH] [STARTSTATION] [ENDSTATION] [TRAINNUMBER]`


running Uning Test: <br>
`cd unitTest` then run `go Test`



## Proposed bonuses:

1. implemented advanced Error handling.

2. implemented super advanced Error handling.

4. It names problematic entities, or specifies the line on which the error occurs.

5. Unit Testing implemented for each major functions.


## Algorithms used
- Edmonds Karp for the defoult path routing
- Dinic accessible through the `-alg Dinic flag`

# Authors

Ellin and Alia