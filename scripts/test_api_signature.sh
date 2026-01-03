#!/bin/bash

# API 签名验证测试脚本
# 使用方法: ./test_signature.sh

BASE_URL="http://localhost:3000"
SECRET="your-secret-key-12345678"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=================================================="
echo "  API 签名验证测试"
echo "=================================================="
echo ""

# 生成签名的函数
generate_signature() {
    local method=$1
    local path=$2
    local timestamp=$3
    local nonce=$4
    local body=$5
    
    # 构建签名字符串
    local sign_string="${method}\n${path}\n${timestamp}\n${nonce}\n${body}"
    
    # 使用 openssl 生成 HMAC-SHA256 签名
    echo -ne "$sign_string" | openssl dgst -sha256 -hmac "$SECRET" | awk '{print $2}'
}

# 生成随机 Nonce
generate_nonce() {
    cat /dev/urandom | LC_ALL=C tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1
}

# 测试1: 正确的签名验证（修改密码）
echo -e "${YELLOW}测试1: 正确的签名 - 修改密码${NC}"
TIMESTAMP=$(date +%s)
NONCE=$(generate_nonce)
BODY='{"old_password":"123456","new_password":"newpass123"}'
METHOD="PUT"
PATH="/api/v1/users/password"
SIGNATURE=$(generate_signature "$METHOD" "$PATH" "$TIMESTAMP" "$NONCE" "$BODY")

echo "请求参数:"
echo "  Method: $METHOD"
echo "  Path: $PATH"
echo "  Timestamp: $TIMESTAMP"
echo "  Nonce: $NONCE"
echo "  Signature: $SIGNATURE"
echo ""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X PUT "${BASE_URL}${PATH}" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: $TIMESTAMP" \
    -H "X-Nonce: $NONCE" \
    -H "X-Signature: $SIGNATURE" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
    -d "$BODY")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BODY_RESPONSE=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "401" ] && echo "$BODY_RESPONSE" | grep -q "未登录"; then
    echo -e "${GREEN}✅ 签名验证通过（需要 JWT token）${NC}"
elif [ "$HTTP_CODE" = "200" ]; then
    echo -e "${GREEN}✅ 请求成功${NC}"
else
    echo -e "${RED}❌ 测试失败 (HTTP $HTTP_CODE)${NC}"
    echo "Response: $BODY_RESPONSE"
fi
echo ""
echo "------------------------------------------------"
echo ""

# 测试2: 签名不匹配
echo -e "${YELLOW}测试2: 签名不匹配${NC}"
TIMESTAMP=$(date +%s)
NONCE=$(generate_nonce)
BODY='{"email":"test@example.com"}'
SIGNATURE="invalid_signature_1234567890abcdef1234567890abcdef1234567890abcd"

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X PUT "${BASE_URL}/api/v1/users/email" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: $TIMESTAMP" \
    -H "X-Nonce: $NONCE" \
    -H "X-Signature: $SIGNATURE" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
    -d "$BODY")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BODY_RESPONSE=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "401" ] && echo "$BODY_RESPONSE" | grep -q "签名验证失败"; then
    echo -e "${GREEN}✅ 正确拒绝（签名不匹配）${NC}"
else
    echo -e "${RED}❌ 测试失败 - 应该拒绝无效签名${NC}"
    echo "Response: $BODY_RESPONSE"
fi
echo ""
echo "------------------------------------------------"
echo ""

# 测试3: 时间戳过期
echo -e "${YELLOW}测试3: 时间戳过期（10分钟前）${NC}"
TIMESTAMP=$(($(date +%s) - 600))  # 10分钟前
NONCE=$(generate_nonce)
BODY='{"phone":"13800138000"}'
SIGNATURE=$(generate_signature "PUT" "/api/v1/users/phone" "$TIMESTAMP" "$NONCE" "$BODY")

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X PUT "${BASE_URL}/api/v1/users/phone" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: $TIMESTAMP" \
    -H "X-Nonce: $NONCE" \
    -H "X-Signature: $SIGNATURE" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
    -d "$BODY")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BODY_RESPONSE=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "401" ] && echo "$BODY_RESPONSE" | grep -q "timestamp expired"; then
    echo -e "${GREEN}✅ 正确拒绝（时间戳过期）${NC}"
else
    echo -e "${RED}❌ 测试失败 - 应该拒绝过期时间戳${NC}"
    echo "Response: $BODY_RESPONSE"
fi
echo ""
echo "------------------------------------------------"
echo ""

# 测试4: 重放攻击检测
echo -e "${YELLOW}测试4: 重放攻击检测（使用相同 Nonce 两次）${NC}"
TIMESTAMP=$(date +%s)
NONCE=$(generate_nonce)
BODY='{"test":"data"}'
SIGNATURE=$(generate_signature "POST" "/api/v1/messages" "$TIMESTAMP" "$NONCE" "$BODY")

echo "第一次请求（应该通过或需要JWT）:"
RESPONSE1=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST "${BASE_URL}/api/v1/messages" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: $TIMESTAMP" \
    -H "X-Nonce: $NONCE" \
    -H "X-Signature: $SIGNATURE" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
    -d "$BODY")

HTTP_CODE1=$(echo "$RESPONSE1" | grep "HTTP_CODE" | cut -d: -f2)
echo "  HTTP Code: $HTTP_CODE1"

sleep 1

echo "第二次请求（相同 Nonce，应该被拒绝）:"
RESPONSE2=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST "${BASE_URL}/api/v1/messages" \
    -H "Content-Type: application/json" \
    -H "X-Timestamp: $TIMESTAMP" \
    -H "X-Nonce: $NONCE" \
    -H "X-Signature: $SIGNATURE" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
    -d "$BODY")

HTTP_CODE2=$(echo "$RESPONSE2" | grep "HTTP_CODE" | cut -d: -f2)
BODY_RESPONSE2=$(echo "$RESPONSE2" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE2" = "401" ] && echo "$BODY_RESPONSE2" | grep -q "重放攻击"; then
    echo -e "${GREEN}✅ 正确检测到重放攻击${NC}"
else
    echo -e "${RED}❌ 测试失败 - 应该检测到重放攻击${NC}"
    echo "Response: $BODY_RESPONSE2"
fi
echo ""
echo "------------------------------------------------"
echo ""

# 测试5: 可选签名验证（登录接口）
echo -e "${YELLOW}测试5: 可选签名验证 - 登录接口（无签名应该放行）${NC}"
BODY='{"login_id":"testuser","password":"123456"}'

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST "${BASE_URL}/api/v1/auth/login/using-password" \
    -H "Content-Type: application/json" \
    -d "$BODY")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d: -f2)
BODY_RESPONSE=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')

# 登录接口无签名应该不会因为签名失败而拒绝（可能因为用户名密码错误）
if [ "$HTTP_CODE" != "401" ] || ! echo "$BODY_RESPONSE" | grep -q "签名验证失败"; then
    echo -e "${GREEN}✅ 可选签名验证正常（无签名时放行）${NC}"
else
    echo -e "${RED}❌ 测试失败 - 可选签名不应该强制要求签名${NC}"
    echo "Response: $BODY_RESPONSE"
fi
echo ""
echo "------------------------------------------------"
echo ""

# 测试6: 查看 Prometheus 指标
echo -e "${YELLOW}测试6: 检查 Prometheus 监控指标${NC}"
METRICS=$(curl -s "${BASE_URL}/metrics" | grep "gohub_api_signature")

if [ -n "$METRICS" ]; then
    echo -e "${GREEN}✅ 监控指标已记录${NC}"
    echo ""
    echo "相关指标:"
    echo "$METRICS" | head -20
else
    echo -e "${YELLOW}⚠️  未找到签名验证指标（可能还没有请求）${NC}"
fi
echo ""

echo "=================================================="
echo "  测试完成"
echo "=================================================="
echo ""
echo "注意事项:"
echo "1. 某些测试需要有效的 JWT token（替换 YOUR_JWT_TOKEN_HERE）"
echo "2. SECRET 需要与服务器配置的 SIGNATURE_SECRET 一致"
echo "3. 确保 Redis 服务已启动（用于 Nonce 防重放）"
echo "4. 查看完整日志: tail -f storage/logs/gohub.log"
echo ""
