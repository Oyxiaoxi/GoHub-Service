# Swagger 文档访问问题修复

## 问题
访问 `http://localhost:3000/swagger/index.html` 时出现错误：
```
Failed to load API definition.
GET /swagger/doc.json 返回 500 错误
```

## 原因
1. Swagger 文档未生成（缺少 `docs/docs.go`）
2. `main.go` 未导入 docs 包

## 解决方案

### 1. 生成 Swagger 文档
```bash
# 添加 go/bin 到 PATH
export PATH=$PATH:$HOME/go/bin

# 生成文档
swag init --parseDependency --parseInternal
```

或使用 Makefile：
```bash
make -f Makefile.swagger swagger
```

### 2. 导入 docs 包
在 `main.go` 中添加：
```go
import (
    _ "GoHub-Service/docs"  // 匿名导入，初始化 Swagger 文档
)
```

### 3. 重新编译并启动
```bash
go build -o tmp/main .
./tmp/main serve
```

## 验证
访问 `http://localhost:3000/swagger/index.html` 应该能正常显示 API 文档。

## 文件生成
运行 `swag init` 后会生成：
- `docs/docs.go` - Swagger 文档入口
- `docs/swagger.json` - JSON 格式文档
- `docs/swagger.yaml` - YAML 格式文档

## 注意事项
1. 每次修改 Swagger 注解后需要重新运行 `swag init`
2. docs 包需要匿名导入（`_ "GoHub-Service/docs"`）
3. 确保 swag 工具已安装：`go install github.com/swaggo/swag/cmd/swag@latest`

## 已修复
✅ 生成 Swagger 文档文件
✅ 在 main.go 中导入 docs 包
✅ 重新编译项目
