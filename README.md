# magpie-lite

[magpie](https://github.com/yetone/magpie) 的组织分支（SophiaData）。功能与上游一致，只改分发方式和几个默认行为。

## 与上游的差异

| 项 | 本仓库 |
| --- | --- |
| 更新来源 | 只读本仓库的 Releases（`feed.json`），不连上游站点 |
| 自动更新 | 关闭：不会自己检查、下载或重启；`magpie update` 或设置里的"检查"才会 |
| 合作伙伴 / 赞助商 | 没有：不拉取列表、不显示那一栏、不上报计数，环境变量也打不开 |
| macOS 安装包 | 除 zip（`magpie update` 用）外，另出 dmg（拖进"应用程序"手动安装） |
| CI | Test 跑 ubuntu 和 Windows；macOS 全量测试每天定时跑，不阻塞发版 |

改动细节、上游同步与发版流程见 [docs/fork.md](docs/fork.md)。

## 下载

从本仓库的 [Releases](https://github.com/SophiaData/magpie-lite/releases/latest)：

| 系统 | 文件 |
| --- | --- |
| Windows | `magpie-windows-amd64.exe`（另有 `arm64`） |
| macOS · Apple Silicon | `magpie-darwin-arm64.dmg` 或 `magpie-darwin-arm64.zip` |
| macOS · Intel | `magpie-darwin-amd64.dmg` 或 `magpie-darwin-amd64.zip` |
| Linux 桌面 | `magpie-linux-amd64`（需 GTK 3 与 WebKitGTK 4.1） |
| 终端版 | `magpie-cli-<os>-<arch>` |

dmg 用于手动安装，zip 是 `magpie update` 取的文件。

## macOS 签名

ad-hoc 签名，未做 Developer ID 签名与公证。首次打开可能被 Gatekeeper 拦下：右键选择"打开"，或到"系统设置 → 隐私与安全性"里允许。更新时会校验 feed 的 SHA-256 和签名完整性，但不会识别开发者或团队身份。

## 更新

自动更新关闭。手动更新用 `magpie update`，或设置页的"检查更新"。

## 从源码构建

见 [docs/fork.md](docs/fork.md)：Windows 用 `make release-windows`，macOS 用 `make release-darwin`，Linux 桌面用 `make release-linux`，终端版用 `make release-cli`。

## 上游文档

功能、配置项与 CLI 参考以上游为准：[上游 README](https://github.com/yetone/magpie#readme)、[reference](docs/reference.md)、[文档站](https://usemagpie.ai/docs/start)。

上游文档里"从官网下载"、"自动更新"、"合作伙伴"这些描述不适用于本仓库。

## 隐私

上游的匿名用量统计仍在本 build 里（每天一次：随机 id、版本、系统，以及使用的 agent、provider、model 名称），可在设置 → 隐私关闭，或设 `DO_NOT_TRACK=1`。合作伙伴的展示与点击计数在本 build 中不存在。

## 许可证

[MIT](LICENSE)，与上游相同。

---

magpie 与 [Apache Magpie](https://magpie.apache.org) 及 Apache 软件基金会无关。
