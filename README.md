[![Go Report Card](https://goreportcard.com/badge/github.com/do87/poly/src?1)](https://goreportcard.com/report/github.com/do87/poly/src)

# Poly

Poly is a project focused on implementing and managing agents that can run plans in various environments

The project consists of a mesh API which is the single source of truth, and agents that can be extended and implemented according to the developer's needs.

<br />

## Workflow

| | |
| --- | --- |
| ![workflow](statics/workflow.svg) | 1. The agent registers itself with the API and retrieves an access token<br />2. The agent sends periodic liveness pings to the API<br />3. The agent is polling for new runs to perform |

<br><br><br>

## Run Lifecycle

![workflow](statics/lifecycle.svg)