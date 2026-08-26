# faraday-corr：Go 均匀腐蚀 Faraday 速率 Web 服务（当量公式 + Stern–Geary + 前端控制台）

由 i_corr、M、n、ρ 计算质量损失率与年腐蚀深度；提供 `/api/rate`、`/api/polarization`、`/api/temperature` 与嵌入网页。

## 构建 / 运行 / 测试

```text
go build ./...
./faraday-corr -http :8080
curl -s http://127.0.0.1:8080/api/example
go run . rate example/fe-seawater.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name faraday-corr-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port faraday-corr-b14 8080 | cut -d: -f2)/api/example
docker rm -f faraday-corr-b14
```
