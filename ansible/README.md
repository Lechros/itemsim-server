# Ansible Deployment

## 사전 요구사항

- 대상 서버에 Docker 설치
- Swarm 모드 초기화

## 사용법

### 환경 변수 설정
```bash
export DEPLOY_HOST="your-server-ip"
export DEPLOY_USER="deploy"
export DEPLOY_KEY="ssh-key"
export DOCKER_IMAGE_NAME="ghcr.io/lechros/itemsim-server"
export DOCKER_TAG="latest"
export DOCKER_DIGEST="sha256:..."
export METRICS_PASSWORD="your-secure-metrics-password"
export RESOURCES_PATH="resources"
```
