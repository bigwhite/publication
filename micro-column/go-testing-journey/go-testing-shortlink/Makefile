.PHONY: test-unit test-integration test-all build

build:
	go build github.com/bigwhite/shortlink/cmd/server

# 只运行快速的单元测试，这是本地开发中最常用的命令
test-unit:
	@echo "Running unit tests..."
	@go test -v -tags=unit ./...

# 只运行慢速的、需要环境的集成测试
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# 运行所有测试
test-all:
	@echo "Running all tests..."
	@go test -v -tags="unit,integration" ./...


.PHONY: test-e2e
test-e2e:
	@echo "Running end-to-end tests..."
	@echo "This may take a while as it builds and starts multiple Docker containers..."
	@cd test/e2e && go test -v -tags=e2e .
