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

<img src="statics/lifecycle.svg" alt="lifecycle" align="left">

1. When a run is created it's not assigned to any agent
2. When an agent is assigned to a run, the status changes to pending
3. When an agent starts to execute a run plan, the status changes to running
4. Success / Error are the usual outcomes of a run
5. Cancelled status is given if the run didn't complete due to a shutdown signal

