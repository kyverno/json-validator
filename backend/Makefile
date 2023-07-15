.PHONY: fmt
fmt:
	@echo "Running gci"
	@go run github.com/daixiang0/gci@v0.9.1 write -s standard -s default -s "prefix(github.com/kyverno/json-validator/backend)" .
	@echo "Running gofumpt"
	@go run mvdan.cc/gofumpt@v0.4.0 -w .
