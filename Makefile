# GoHub-Service Makefile

.PHONY: help test test-coverage test-services test-repositories test-elasticsearch test-es-up test-es-down test-all clean

# 默认目标
help:
	@echo "GoHub-Service 测试命令："
	@echo "  make test                  - 运行所有测试"
	@echo "  make test-coverage         - 运行测试并生成覆盖率报告"
	@echo "  make test-services         - 只测试Service层"
	@echo "  make test-repositories     - 只测试Repository层"
	@echo "  make test-elasticsearch    - 运行Elasticsearch集成测试"
	@echo "  make test-es-up            - 启动Elasticsearch集群"
	@echo "  make test-es-down          - 停止Elasticsearch集群"
	@echo "  make test-all              - 运行完整测试套件"
	@echo "  make clean                 - 清理测试缓存和覆盖率文件"

# 运行所有测试
test:
	@echo "🧪 运行所有测试..."
	go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	@echo "📊 生成测试覆盖率报告..."
	@go test -coverprofile=coverage.out ./...
	@echo "\n📈 总体覆盖率："
	@go tool cover -func=coverage.out | grep total
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ HTML报告已生成: coverage.html"

# 只测试Service层
test-services:
	@echo "🧪 测试Service层..."
	go test -v -cover ./app/services/...

# 只测试Repository层
test-repositories:
	@echo "🧪 测试Repository层..."
	go test -v -cover ./app/repositories/...

# 启动Elasticsearch集群
test-es-up:
	@echo "🚀 启动Elasticsearch集群..."
	docker-compose -f docker-compose.elasticsearch.yml up -d
	@echo "⏳ 等待Elasticsearch启动..."
	@sleep 15
	@echo "✅ Elasticsearch已启动"
	@curl -s http://localhost:9200/_cluster/health | grep -q '"status":"green"' && echo "✅ 集群状态: GREEN" || echo "⚠️  集群状态: 检查中..."

# 停止Elasticsearch集群
test-es-down:
	@echo "🛑 停止Elasticsearch集群..."
	docker-compose -f docker-compose.elasticsearch.yml down
	@echo "✅ Elasticsearch已停止"

# 运行Elasticsearch集成测试
test-elasticsearch: test-es-up
	@echo "🧪 运行Elasticsearch集成测试..."
	@go test -v -cover -timeout 120s ./pkg/elasticsearch/... || (make test-es-down && exit 1)
	@echo "\n✅ Elasticsearch测试完成"
	@make test-es-down

# 运行Elasticsearch基准测试
test-es-bench:
	@echo "🚀 运行Elasticsearch基准测试..."
	@docker-compose -f docker-compose.elasticsearch.yml up -d
	@sleep 15
	@go test -bench=. -benchmem -benchtime=10s ./pkg/elasticsearch/ || true
	@docker-compose -f docker-compose.elasticsearch.yml down


test-all:
	@echo "🚀 运行完整测试套件..."
	@echo "\n1️⃣ 清理测试缓存..."
	@go clean -testcache
	@echo "\n2️⃣ 运行Service层测试..."
	@go test -v -cover ./app/services/... || true
	@echo "\n3️⃣ 运行Repository层测试..."
	@go test -v -cover ./app/repositories/... || true
	@echo "\n4️⃣ 运行Elasticsearch集成测试..."
	@make test-elasticsearch || true
	@echo "\n5️⃣ 生成覆盖率报告..."
	@go test -coverprofile=coverage.out ./... 2>&1 | grep -v "no test files"
	@echo "\n📊 测试统计："
	@echo "测试文件数: $$(find . -name '*_test.go' | wc -l | tr -d ' ')"
	@echo "测试函数数: $$(grep -r '^func Test' --include='*_test.go' . | wc -l | tr -d ' ')"
	@echo "\n📈 覆盖率报告："
	@go tool cover -func=coverage.out | grep total
	@go tool cover -html=coverage.out -o coverage.html
	@echo "\n✅ 完成！查看详细报告: open coverage.html"

# 清理
clean:
	@echo "🧹 清理测试文件..."
	@rm -f coverage.out coverage.html
	@go clean -testcache
	@echo "✅ 清理完成"
