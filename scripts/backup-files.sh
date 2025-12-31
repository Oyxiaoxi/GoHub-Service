#!/bin/bash
# 文件备份脚本

set -e

# 配置
BACKUP_DIR="/backup/gohub/files"
DATE=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=30
SOURCE_DIR="/var/www/gohub/public/uploads"

# 创建备份目录
mkdir -p $BACKUP_DIR

echo "开始备份上传文件"
echo "备份时间: $(date)"
echo "源目录: $SOURCE_DIR"

# 检查源目录是否存在
if [ ! -d "$SOURCE_DIR" ]; then
    echo "错误: 源目录不存在: $SOURCE_DIR"
    exit 1
fi

# 执行备份
tar -czf "$BACKUP_DIR/uploads_${DATE}.tar.gz" -C "$SOURCE_DIR" .

if [ $? -eq 0 ]; then
    echo "备份成功: uploads_${DATE}.tar.gz"
    
    # 计算备份文件大小
    SIZE=$(du -h "$BACKUP_DIR/uploads_${DATE}.tar.gz" | cut -f1)
    echo "备份文件大小: $SIZE"
else
    echo "备份失败"
    exit 1
fi

# 删除旧备份
echo "清理 $KEEP_DAYS 天前的备份..."
find $BACKUP_DIR -name "uploads_*.tar.gz" -mtime +$KEEP_DAYS -delete

# 统计备份数量
COUNT=$(ls -1 $BACKUP_DIR/uploads_*.tar.gz 2>/dev/null | wc -l)
echo "当前保留备份数量: $COUNT"

echo "备份完成"
