# Pathfinder

A CLI path-finding algorithm that finds the most efficient paths to move trains from one destination to another with minimal turns.

## Setup and Usage

Clone the repository using HTTPS:

```bash
git clone https://gitea.kood.tech/aliadaahirmohamed/pathfinder.git
```

Or using SSH:

```bash
git clone git@gitea.kood.tech:aliadaahirmohamed/pathfinder.git
```

Navigate to the project:

```bash
cd pathfinder
```

Run using the Edmonds-Karp algorithm:

```bash
go run . [INPUTFILEPATH] [STARTSTATION] [ENDSTATION] [TRAINNUMBER]
```

Or run using the Dinic algorithm:

```bash
go run . -alg Dinic [INPUTFILEPATH] [STARTSTATION] [ENDSTATION] [TRAINNUMBER]
```

Run unit tests:

```bash
cd unitTest
go test
```

## Proposed Bonuses

1. Implemented advanced error handling.
2. Implemented super advanced error handling.
3. Names problematic entities or specifies the line on which the error occurs.
4. Unit testing implemented for each major function.

## Algorithms Used

- Edmonds-Karp for the default path routing.
- Dinic, accessible through the `-alg Dinic` flag.

# Authors

- Ellin
- Alia