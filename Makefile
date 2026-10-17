# stack-agents task runner.
#   make                        show this help
#   make check                  generated-file drift check, then the test suite
#   make sync                   regenerate the per-agent rule files from AGENTS.md
#   make diff                   show what sync would change, without writing
#   make test                   run the test suite alone
#   make install TARGET=<dir>   install into a project directory

.PHONY: help check test sync diff install

help:
	@echo "stack-agents"
	@echo "  make check                  check generated files, then run tests"
	@echo "  make sync                   regenerate agent rule files from AGENTS.md"
	@echo "  make diff                   show the difference for each stale file"
	@echo "  make test                   run the test suite alone"
	@echo "  make install TARGET=<dir>   install into a project directory"

check:
	@bash scripts/sync.sh --check
	@bash scripts/test.sh

test:
	@bash scripts/test.sh

sync:
	@bash scripts/sync.sh

diff:
	@bash scripts/sync.sh --diff

install:
	@test -n "$(TARGET)" || { echo "usage: make install TARGET=/path/to/project"; exit 2; }
	@bash scripts/install.sh "$(TARGET)"
