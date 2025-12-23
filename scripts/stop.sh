#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  停止 SIEM 日志审计系统${NC}"
echo -e "${GREEN}========================================${NC}"

echo -e "${YELLOW}停止所有容器...${NC}"
docker-compose down

echo -e "${YELLOW}清理未使用的镜像...${NC}"
docker image prune -f

echo -e "${GREEN}服务已停止${NC}"
