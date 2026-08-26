# faraday-corr

faraday-corr 按法拉第当量核算均匀腐蚀速率。输入腐蚀电流密度 i_corr、摩尔质量 M、价数 n 与密度 ρ，给出质量损失率与年腐蚀深度：

```
mdot  = M · i_corr / (n · F)
CR    = K · (M / n) · i_corr / ρ
```

F 钉在 CODATA 值 96485.33212 C/mol；K 把 μA/cm² 与 g/cm³ 折成 mm/y。电流加倍则两种速率加倍，价数加倍则减半，密度加倍时质量率不变、深度率减半。另可由极化电阻（Stern–Geary）反推 i_corr、按 Arrhenius 把 i_corr 变到另一温度、或按涂层破损面积把总电流从几何面积上扣下来。它不是腐蚀数据库，也不输出维修工单。

## 用法

最直接的一条命令（铁在海水中 10 μA/cm²）：

```text
go run . rate example/fe-seawater.json
```

输出腐蚀深度约 0.1159 mm/y、年失重约 0.09126 g/(cm²·y)。其它入口：

```text
faraday-corr rate <算例.json>
faraday-corr reverse / schedule / life / chain / json / validate / metals / template
faraday-corr -http :8080
faraday-corr help
```

浏览器打开 http://localhost:8080 可改参数并打 `/api/rate`。容器默认拉起同一服务。

## 算例格式

`example/fe-seawater.json` 可直接运行：

```json
{
  "metal": "Fe",
  "i_corr": 10.0,
  "M": 55.845,
  "n": 2,
  "rho": 7.874,
  "area": 1.0,
  "duration_y": 1.0
}
```

- `i_corr` 单位 μA/cm²，允许 0（此时 CR=0）。
- 给出 `metal` 且省略 M/n/ρ 时从内置登记表补全。
- `reverse` 另要 `target_cr_mm_y` 或 `target_annual_loss` 之一。

## 关键约定

- 两条深度率路径必须闭合：质量路径 `10·annual/ρ` 与直接 `K·(M/n)·i/ρ` 一致。
- Stern–Geary：`i_corr = B / Rp`，`B = βa βc / (2.303 (βa+βc))`；高过电位下 Butler–Volmer 回到 Tafel。
- Arrhenius：`i(T) = i_ref · exp(-Ea/R · (1/T − 1/Tref))`，再代入法拉第公式；先变温再求 CR 与先求 CR 再乘同一温度因子必须一致。
- 涂层破损：局部 CR 仍由 i_corr 决定，总电流与总失重按破损面积，不是按几何面积。
- 点蚀因子 `pf≥1`：局部穿透率 = 均匀率 × pf，穿孔寿命按比例缩短。电偶混合电流随阴极面积增大而增大；扩散膜越薄极限电流越高。

## API

- `GET /health` 与 `GET /api/health` — 探活
- `GET /api/example` — 预置铁海水算例
- `POST /api/rate` — 法拉第速率（请求体即算例 JSON）
- `POST /api/polarization` — 由 Rp 与 Tafel 斜率求速率
- `POST /api/temperature` — Arrhenius 变温后再求速率
- `POST /api/coating` — 破损面积下的总电流与失重

非法参数返回 `{"error":"…"}` 与 400。

## 构建与测试

```text
go build ./...
go test ./...
go run . -http :8080
```

纯标准库，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。
