# faraday-corr：均匀腐蚀 Faraday 速率核算命令行工具

faraday-corr 根据腐蚀电流密度 i_corr、金属摩尔质量 M、价数 n 与密度 ρ，
用法拉第当量关系计算质量损失率、年腐蚀深度（mm/y 或 μm/y）与累计失重。

## 构建 / 运行 / 测试

```text
go build ./...
go run . rate example/fe-seawater.json
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
