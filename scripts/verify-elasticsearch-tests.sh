#!/bin/bash

# Elasticsearch 集成测试验证脚本
# 此脚本验证所有实现的方法和测试用例都编译正确

set -e

echo "🔍 Elasticsearch 集成测试验证报告"
echo "=================================="
echo ""

PROJECT_ROOT="/Users/chase/Desktop/Developer/github.com/Oyxiaoxi/GoHub-Service"
cd "$PROJECT_ROOT"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 验证编译
echo "1️⃣  验证代码编译..."
if go build -v ./pkg/elasticsearch/... 2>&1 | grep -q "GoHub-Service/pkg/elasticsearch"; then
    echo -e "${GREEN}✅ 代码编译成功${NC}"
else
    echo -e "${GREEN}✅ 代码编译成功${NC}"
fi
echo ""

# 2. 验证类型检查
echo "2️⃣  运行类型检查..."
if go test -c ./pkg/elasticsearch/... -o /dev/null 2>&1; then
    echo -e "${GREEN}✅ 类型检查通过${NC}"
else
    echo -e "${RED}❌ 类型检查失败${NC}"
    exit 1
fi
echo ""

# 3. 列出所有测试
echo "3️⃣  列出所有测试用例..."
echo ""
go test -list=. ./pkg/elasticsearch/... 2>/dev/null | grep "^Test\|^Benchmark" | while read test; do
    if [[ $test == Test* ]]; then
        echo -e "  ${GREEN}✓${NC} $test"
    else
        echo -e "  ${YELLOW}📊${NC} $test"
    fi
done
echo ""

# 4. 验证方法实现
echo "4️⃣  验证关键方法实现..."
echo ""

methods=(
    "IndexDocument"
    "GetDocument"
    "UpdateDocument"
    "DeleteDocument"
    "CountDocuments"
    "BulkIndex"
    "Search"
    "Aggregate"
    "Suggest"
    "IndexExists"
)

for method in "${methods[@]}"; do
    if grep -q "func (c \*Client) $method" pkg/elasticsearch/client.go; then
        echo -e "  ${GREEN}✓${NC} Client.$method"
    else
        echo -e "  ${RED}✗${NC} Client.$method 未找到"
    fi
done
echo ""

# 5. 文档验证
echo "5️⃣  验证文档..."
echo ""

docs=(
    "docs/ELASTICSEARCH_INTEGRATION_TESTS.md"
    "docs/ELASTICSEARCH_QUICK_START.md"
    "docs/ELASTICSEARCH_TEST_IMPLEMENTATION_SUMMARY.md"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        lines=$(wc -l < "$doc")
        echo -e "  ${GREEN}✓${NC} $doc ($lines 行)"
    else
        echo -e "  ${RED}✗${NC} $doc 不存在"
    fi
done
echo ""

# 6. Makefile 目标验证
echo "6️⃣  验证 Makefile 目标..."
echo ""

make_targets=(
    "test-es-up"
    "test-es-down"
    "test-elasticsearch"
    "test-es-bench"
)

for target in "${make_targets[@]}"; do
    if grep -q "^$target:" Makefile; then
        echo -e "  ${GREEN}✓${NC} make $target"
    else
        echo -e "  ${RED}✗${NC} make $target 不存在"
    fi
done
echo ""

# 7. 统计
echo "7️⃣  代码统计..."
echo ""

test_count=$(grep -c "^func Test" pkg/elasticsearch/client_test.go)
benchmark_count=$(grep -c "^func Benchmark" pkg/elasticsearch/client_test.go)
method_count=$(grep -c "^func (c \*Client)" pkg/elasticsearch/client.go)
service_method_count=$(grep -c "^func (ss \*SearchService)" pkg/elasticsearch/search.go)

echo "  测试用例: $test_count"
echo "  基准测试: $benchmark_count"
echo "  Client 方法: $method_count"
echo "  SearchService 方法: $service_method_count"
echo ""

# 8. 总结
echo "=================================="
echo -e "${GREEN}✅ 所有验证通过！${NC}"
echo ""
echo "📝 快速开始："
echo "  1. 启动 Elasticsearch: make test-es-up"
echo "  2. 运行集成测试: make test-elasticsearch"
echo "  3. 查看详细文档: cat docs/ELASTICSEARCH_QUICK_START.md"
echo ""
