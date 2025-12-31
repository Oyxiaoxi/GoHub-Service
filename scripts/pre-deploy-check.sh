#!/bin/bash
# 生产环境部署检查脚本

set -e

echo "=================================="
echo "GoHub-Service 生产环境部署检查"
echo "=================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查函数
check_pass() {
    echo -e "${GREEN}✓${NC} $1"
}

check_fail() {
    echo -e "${RED}✗${NC} $1"
}

check_warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# 错误计数
ERRORS=0
WARNINGS=0

echo "1. 检查环境变量配置"
echo "-----------------------------------"

# 检查 .env 文件是否存在
if [ ! -f .env ]; then
    check_fail ".env 文件不存在"
    ERRORS=$((ERRORS + 1))
else
    check_pass ".env 文件存在"
    
    # 加载环境变量
    set -a
    source .env
    set +a
    
    # 检查关键环境变量
    if [ "$APP_ENV" != "production" ]; then
        check_fail "APP_ENV 必须设置为 production (当前: $APP_ENV)"
        ERRORS=$((ERRORS + 1))
    else
        check_pass "APP_ENV = production"
    fi
    
    if [ "$APP_DEBUG" != "false" ]; then
        check_fail "APP_DEBUG 必须设置为 false (当前: $APP_DEBUG)"
        ERRORS=$((ERRORS + 1))
    else
        check_pass "APP_DEBUG = false"
    fi
    
    if [ -z "$APP_KEY" ] || [ "$APP_KEY" == "CHANGE_THIS_32_CHARACTER_KEY_HERE" ]; then
        check_fail "APP_KEY 必须设置为强随机字符串"
        ERRORS=$((ERRORS + 1))
    else
        check_pass "APP_KEY 已设置"
    fi
    
    if [ -z "$JWT_KEY" ] || [ "$JWT_KEY" == "CHANGE_JWT_SECRET_KEY_HERE" ]; then
        check_fail "JWT_KEY 必须设置为强随机字符串"
        ERRORS=$((ERRORS + 1))
    else
        check_pass "JWT_KEY 已设置"
    fi
    
    if [ "$DB_USERNAME" == "root" ]; then
        check_warn "建议不要使用 root 作为数据库用户"
        WARNINGS=$((WARNINGS + 1))
    else
        check_pass "数据库用户名不是 root"
    fi
    
    if [ -z "$REDIS_PASSWORD" ]; then
        check_fail "REDIS_PASSWORD 必须设置"
        ERRORS=$((ERRORS + 1))
    else
        check_pass "Redis 密码已设置"
    fi
    
    if [ "$LOG_LEVEL" == "debug" ]; then
        check_warn "生产环境建议使用 info 或 error 日志级别"
        WARNINGS=$((WARNINGS + 1))
    else
        check_pass "日志级别设置正确"
    fi
fi

echo ""
echo "2. 检查数据库连接"
echo "-----------------------------------"

if command -v mysql &> /dev/null; then
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USERNAME" -p"$DB_PASSWORD" -e "USE $DB_DATABASE" &> /dev/null; then
        check_pass "数据库连接成功"
    else
        check_fail "无法连接到数据库"
        ERRORS=$((ERRORS + 1))
    fi
else
    check_warn "未安装 mysql 客户端，跳过数据库连接检查"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "3. 检查 Redis 连接"
echo "-----------------------------------"

if command -v redis-cli &> /dev/null; then
    if [ -n "$REDIS_PASSWORD" ]; then
        if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" ping &> /dev/null; then
            check_pass "Redis 连接成功"
        else
            check_fail "无法连接到 Redis"
            ERRORS=$((ERRORS + 1))
        fi
    else
        if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping &> /dev/null; then
            check_pass "Redis 连接成功"
        else
            check_fail "无法连接到 Redis"
            ERRORS=$((ERRORS + 1))
        fi
    fi
else
    check_warn "未安装 redis-cli，跳过 Redis 连接检查"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "4. 检查文件权限"
echo "-----------------------------------"

# 检查 storage 目录
if [ -d "storage" ]; then
    if [ -w "storage" ]; then
        check_pass "storage 目录可写"
    else
        check_fail "storage 目录不可写"
        ERRORS=$((ERRORS + 1))
    fi
else
    check_warn "storage 目录不存在"
    WARNINGS=$((WARNINGS + 1))
fi

# 检查 public/uploads 目录
if [ -d "public/uploads" ]; then
    if [ -w "public/uploads" ]; then
        check_pass "public/uploads 目录可写"
    else
        check_fail "public/uploads 目录不可写"
        ERRORS=$((ERRORS + 1))
    fi
else
    check_warn "public/uploads 目录不存在"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "5. 检查 Go 环境"
echo "-----------------------------------"

if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    check_pass "Go 已安装 ($GO_VERSION)"
    
    # 检查依赖
    if go mod verify &> /dev/null; then
        check_pass "Go 模块验证通过"
    else
        check_fail "Go 模块验证失败"
        ERRORS=$((ERRORS + 1))
    fi
else
    check_fail "未安装 Go"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "6. 检查编译"
echo "-----------------------------------"

if go build -o /tmp/gohub_test_build &> /dev/null; then
    check_pass "项目编译成功"
    rm -f /tmp/gohub_test_build
else
    check_fail "项目编译失败"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "7. 检查测试"
echo "-----------------------------------"

if go test ./... &> /dev/null; then
    check_pass "所有测试通过"
else
    check_warn "部分测试失败"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "8. 检查端口占用"
echo "-----------------------------------"

if command -v lsof &> /dev/null; then
    if lsof -i:"$APP_PORT" &> /dev/null; then
        check_warn "端口 $APP_PORT 已被占用"
        WARNINGS=$((WARNINGS + 1))
    else
        check_pass "端口 $APP_PORT 可用"
    fi
else
    check_warn "未安装 lsof，跳过端口检查"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "9. 检查 SSL 证书（如果使用 HTTPS）"
echo "-----------------------------------"

if [ -f "ssl/cert.pem" ] && [ -f "ssl/key.pem" ]; then
    check_pass "SSL 证书文件存在"
    
    # 检查证书有效期
    if command -v openssl &> /dev/null; then
        CERT_EXPIRY=$(openssl x509 -enddate -noout -in ssl/cert.pem | cut -d= -f2)
        check_pass "证书过期时间: $CERT_EXPIRY"
    fi
else
    check_warn "未找到 SSL 证书文件（如果使用反向代理，可忽略）"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "10. 检查数据库迁移状态"
echo "-----------------------------------"

if [ -f "./main" ]; then
    check_pass "可执行文件存在"
else
    check_warn "未找到可执行文件，请先编译: go build -o main"
    WARNINGS=$((WARNINGS + 1))
fi

echo ""
echo "=================================="
echo "检查完成"
echo "=================================="
echo ""
echo -e "错误: ${RED}$ERRORS${NC}"
echo -e "警告: ${YELLOW}$WARNINGS${NC}"
echo ""

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}发现 $ERRORS 个错误，请修复后再部署${NC}"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}发现 $WARNINGS 个警告，建议检查后再部署${NC}"
    exit 0
else
    echo -e "${GREEN}所有检查通过，可以部署！${NC}"
    exit 0
fi
