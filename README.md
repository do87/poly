[![Go Report Card](https://goreportcard.com/badge/github.com/do87/poly/src?1)](https://goreportcard.com/report/github.com/do87/poly/src)

# Poly

Poly is a project focused on implementing and managing agents that can run plans in various environments

The project consists of a mesh API which is the single source of truth, and agents that can be extended and implemented according to the developer's needs.

<br />

## Workflow

<img src="statics/workflow.svg" alt="workflow" align="left">
<span>

1. The agent registers itself with the API and retrieves an access token
2. The agent is polling the API periodically, and the API returns a list of runs it needs to execute
3. When an agent receives a shutdown signal, it makes an API call to deregister itself
</span>
<br><br><br><br><br><br><br><br>

## Run Lifecycle

<img src="statics/lifecycle.svg" alt="lifecycle" align="right">

1. `Created` is the initial status
2. The status changes to `Pending` when an agent is assigned 
3. The status changes to `Running` when an agent starts to execute a run
4. `Success` / `Error` are the usual outcomes of a run
5. `Canceled` status is given if the run didn't complete due to a shutdown signal
