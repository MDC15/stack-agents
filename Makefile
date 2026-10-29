# stack-agents task runner.
#   make check                  format, vet, tests, then the drift check
#   make sync                   regenerate the per-agent rule files from AGENTS.md
#   make diff                   show what the next sync would change
#   make test                   run the test suite alone
#
# Install into other projects with `stack init --path <dir>`; there is
# intentionally no make target for it.

.PHONY: check test sync diff

GO ?= go

check:
	@bash -c 'test -z "$$($(GO)fmt -l .)"' || { echo "unformatted files:"; $(GO)fmt -l .; exit 1; }
	@$(GO) vet ./...
	@$(GO) test ./...
	@$(GO) run ./cmd/stack check

test:
	@$(GO) test ./...

sync:
	@$(GO) run ./cmd/stack sync --all

diff:
	@$(GO) run ./cmd/stack check --diff
