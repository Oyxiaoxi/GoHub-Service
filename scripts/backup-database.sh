#!/bin/bash
# 数据库备份脚本

set -e

# 配置
BACKUP_DIR="/backup/gohub/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=30

# 从 .env 加载数据库配置
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "错误: 未找到 .env 文件"
    exit 1
fi

# 创建备份目录
mkdir -p $BACKUP_DIR

echo "开始备份数据库: $DB_DATABASE"
echo "备份时间: $(date)"

# 执行备份
mysqldump \
    -h"$DB_HOST" \
    -P"$DB_PORT" \
    -u"$DB_USERNAME" \
    -p"$DB_PASSWORD" \
    --single-transaction \
    --routines \
    --triggers \
    --events \
    --set-gtid-purged=OFF \
    --databases "$DB_DATABASE" | gzip > "$BACKUP_DIR/backup_${DATE}.sql.gz"

if [ $? -eq 0 ]; then
    echo "备份成功: backup_${DATE}.sql.gz"
    
    # 计算备份文件大小
    SIZE=$(du -h "$BACKUP_DIR/backup_${DATE}.sql.gz" | cut -f1)
    echo "备份文件大小: $SIZE"
else
    echo "备份失败"
    exit 1
fi

# 删除旧备份
echo "清理 $KEEP_DAYS 天前的备份..."
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +$KEEP_DAYS -delete

# 统计备份数量
COUNT=$(ls -1 $BACKUP_DIR/backup_*.sql.gz 2>/dev/null | wc -l)
echo "当前保留备份数量: $COUNT"

echo "备份完成"
