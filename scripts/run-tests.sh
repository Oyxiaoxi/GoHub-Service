#!/bin/bash

# GoHub-Service 测试运行脚本
# 用于运行测试并生成覆盖率报告

set -e

echo "🚀 GoHub-Service 测试套件"
echo "================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数：打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 清理之前的覆盖率文件
cleanup() {
    print_info "清理之前的测试文件..."
    rm -f coverage.out coverage.html
    go clean -testcache
    print_success "清理完成"
    echo ""
}

# 运行单元测试
run_tests() {
    print_info "运行单元测试..."
    echo ""
    
    if go test -v ./app/services/... ./app/repositories/... ./pkg/testutil/...; then
        print_success "所有测试通过！"
    else
        print_error "测试失败"
        exit 1
    fi
    echo ""
}

# 生成覆盖率报告
generate_coverage() {
    print_info "生成覆盖率报告..."
    
    # 生成覆盖率文件
    go test -coverprofile=coverage.out ./...
    
    # 显示总体覆盖率
    echo ""
    print_info "总体覆盖率："
    go tool cover -func=coverage.out | grep total
    
    # 生成HTML报告
    go tool cover -html=coverage.out -o coverage.html
    print_success "HTML覆盖率报告已生成: coverage.html"
    echo ""
}

# 显示详细的包级别覆盖率
show_package_coverage() {
    print_info "包级别覆盖率："
    echo ""
    
    echo "Service层："
    go test -cover ./app/services/... 2>&1 | grep -E "coverage:|ok|FAIL"
    
    echo ""
    echo "Repository层："
    go test -cover ./app/repositories/... 2>&1 | grep -E "coverage:|ok|FAIL"
    
    echo ""
    echo "工具包层："
    go test -cover ./pkg/testutil/... 2>&1 | grep -E "coverage:|ok|FAIL" || echo "  (无测试文件)"
    
    echo ""
}

# 检查覆盖率阈值
check_coverage_threshold() {
    print_info "检查覆盖率阈值..."
    
    # 获取总体覆盖率百分比
    total_coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    
    threshold=60
    
    echo ""
    echo "当前覆盖率: ${total_coverage}%"
    echo "目标阈值: ${threshold}%"
    echo ""
    
    # 使用bc进行浮点数比较
    if command -v bc &> /dev/null; then
        if (( $(echo "$total_coverage >= $threshold" | bc -l) )); then
            print_success "覆盖率达标！ (${total_coverage}% >= ${threshold}%)"
        else
            print_warning "覆盖率未达标 (${total_coverage}% < ${threshold}%)"
            echo ""
            print_info "建议："
            echo "  1. 为核心Service层添加更多测试"
            echo "  2. 为Repository层添加更多测试"
            echo "  3. 查看 docs/TESTING_GUIDE.md 获取测试编写指南"
            echo ""
        fi
    else
        # 如果bc不可用，使用简单的整数比较
        total_int=${total_coverage%.*}
        if [ "$total_int" -ge "$threshold" ]; then
            print_success "覆盖率达标！ (${total_coverage}% >= ${threshold}%)"
        else
            print_warning "覆盖率未达标 (${total_coverage}% < ${threshold}%)"
        fi
    fi
    echo ""
}

# 显示测试统计
show_statistics() {
    print_info "测试统计："
    echo ""
    
    # 统计测试文件数量
    test_files=$(find . -name "*_test.go" | wc -l | tr -d ' ')
    echo "测试文件数量: $test_files"
    
    # 统计测试函数数量
    test_functions=$(grep -r "^func Test" --include="*_test.go" . | wc -l | tr -d ' ')
    echo "测试函数数量: $test_functions"
    
    # 统计新增的测试文件
    new_tests=$(find ./app/services -name "*_test.go" -o -name "*_test.go" | wc -l | tr -d ' ')
    echo "新增测试文件: $new_tests (Service + Repository)"
    
    echo ""
}

# 主函数
main() {
    # 检查是否在项目根目录
    if [ ! -f "go.mod" ]; then
        print_error "请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 执行测试流程
    cleanup
    run_tests
    generate_coverage
    show_package_coverage
    check_coverage_threshold
    show_statistics
    
    print_success "测试完成！"
    print_info "查看详细报告: open coverage.html"
    echo ""
}

# 运行主函数
main
