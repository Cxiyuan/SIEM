#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  SIEM 日志审计系统部署脚本${NC}"
echo -e "${GREEN}========================================${NC}"

if [ -f .env ]; then
    echo -e "${YELLOW}加载环境变量文件...${NC}"
    source .env
fi

IMAGE_TAG=${IMAGE_TAG:-latest}
DOCKER_REGISTRY=${DOCKER_REGISTRY:-ghcr.io}
GITHUB_REPOSITORY=${GITHUB_REPOSITORY:-your-org/siem}

echo -e "${YELLOW}镜像配置:${NC}"
echo "  Registry: $DOCKER_REGISTRY"
echo "  Repository: $GITHUB_REPOSITORY"
echo "  Tag: $IMAGE_TAG"

if [ "$1" == "pull" ]; then
    echo -e "${YELLOW}拉取 Docker 镜像...${NC}"
    docker pull ${DOCKER_REGISTRY}/${GITHUB_REPOSITORY}:backend-${IMAGE_TAG}
    docker pull ${DOCKER_REGISTRY}/${GITHUB_REPOSITORY}:frontend-${IMAGE_TAG}
fi

echo -e "${YELLOW}停止现有容器...${NC}"
docker-compose down || true

echo -e "${YELLOW}启动服务...${NC}"
export IMAGE_TAG
export DOCKER_REGISTRY
export GITHUB_REPOSITORY

docker-compose up -d

echo -e "${GREEN}等待服务就绪...${NC}"
sleep 10

if docker-compose ps | grep -q "Up"; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  部署成功！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "访问地址:"
    echo -e "  前端: ${GREEN}http://localhost${NC}"
    echo -e "  后端: ${GREEN}http://localhost:8080${NC}"
    echo -e "  OpenSearch: ${GREEN}http://localhost:9200${NC}"
    echo ""
    echo -e "查看日志: ${YELLOW}docker-compose logs -f${NC}"
    echo -e "停止服务: ${YELLOW}docker-compose down${NC}"
else
    echo -e "${RED}部署失败，请检查日志${NC}"
    docker-compose logs
    exit 1
fi
