<style>
img[src*='#left'] {
    float: left;
    margin-right: 30px;
}
br {
    clear:both;
}
</style>

[![Go Report Card](https://goreportcard.com/badge/github.com/do87/poly/src?1)](https://goreportcard.com/report/github.com/do87/poly/src)

# Poly

Poly is a project focused on implementing and managing agents that can run plans in various environments

The project consists of a mesh API which is the single source of truth, and agents that can be extended and implemented according to the developer's needs.

<br />

## Workflow

![workflow](statics/workflow.svg#left)

1. The agent registers itself with the API and retrieves an access token
2. The agent sends periodic liveness pings to the API
3. The agent is polling for new runs to perform

<br><br><br>

## Run Lifecycle

![workflow](statics/lifecycle.svg)