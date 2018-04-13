# process-supervisor

## Installation

Install using go get:

```go get github.com/serbangilvitu/process-supervisor```

## Usage

Usage of process-supervisor:
- -a string
  - Process arguments
- -i int
  - Check interval (default 5)
  - Allowed values: 1-3600
- -l bool
  - Generate Logs (default *false*)
- -p string
  - Process name (default "")
- -r int
  - Maximum retries (default 3)
- -t int
   - Wait time before restart (default 10)
   - Allowed values: 1-3600
   
## Example

This will make sure that the sleep command is always running.

```process-supervisor -l -p "sleep" -a "infinity"```

Note: if the process is already running with other arguments, it will not be restarted.
