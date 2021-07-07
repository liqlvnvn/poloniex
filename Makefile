.PHONY: coverage ci clear

lint:
	set -euo pipefail
	golangci-lint run --out-format code-climate | jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"'

prepare: lint
