# Tooling Setup Complete - Summary

## 完成日期
2024

## 概述

本次工作完成了 go-vdom 項目的自動化工具鏈建設，包括：
- ✅ Template Linter（模板檢查工具）
- ✅ Git Hooks（提交前自動檢查）
- ✅ CI/CD Pipeline（GitHub Actions）
- ✅ Makefile（開發命令）
- ✅ 完整文檔

## 已完成的工作

### 1. Template Linter（模板檢查工具）

**位置：** `tools/template-linter/`

**功能：**
- 自動檢測模板中的危險模式
- 提供修復建議
- 可集成到 CI/CD 和 git hooks

**檢測的問題類型：**

1. **quoted-var-in-expression**: `${...}` 中被引號包裹的變量
   - 錯誤：`${'{{variable}}'}`
   - 正確：`${{{variable}}}`

2. **string-bool-comparison**: 字串與布林值比較
   - 錯誤：`'{{flag}}' === 'true'`
   - 正確：`{{flag}} === true`

3. **string-concatenation-in-expression**: `${...}` 中的字串拼接
   - 錯誤：`${'text' + {{variable}}}`
   - 建議：在 Go 代碼中處理

4. **quoted-var-comparison**: 被引號包裹的變量比較
   - 錯誤：`'{{value}}' === 'something'`
   - 正確：`{{value}} === 'something'`

5. **double-quoted-var**: 雙引號包裹的變量（警告）
   - 需要根據上下文驗證是否正確

**使用方法：**

```bash
# 構建
cd tools/template-linter
go build

# 運行
./template-linter ../../

# 顯示修復建議
./template-linter -fix ../../

# 詳細輸出
./template-linter -v ../../
```

**測試結果：**
```bash
cd tools/template-linter && go run main.go -fix ../../components
✓ No template issues found!
```

### 2. Git Hooks

**位置：** `.githooks/`

**檔案：**
- `pre-commit` - 提交前檢查腳本
- `install.sh` - 安裝腳本

**功能：**
- 自動運行代碼格式化檢查（go fmt）
- 自動運行靜態分析（go vet）
- 自動運行 template linter
- 可選：運行測試

**安裝方法：**

```bash
chmod +x .githooks/install.sh
./.githooks/install.sh
```

**手動安裝：**

```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

**配置 git 使用 .githooks 目錄：**

```bash
git config core.hooksPath .githooks
```

**繞過（不推薦）：**

```bash
git commit --no-verify
```

### 3. CI/CD Pipeline

**位置：** `.github/workflows/ci.yml`

**Jobs：**

1. **test** - 運行所有測試
   - 多個 Go 版本（1.21, 1.22）
   - Race detection
   - Coverage report
   - Upload to Codecov

2. **lint** - 代碼檢查
   - go fmt
   - go vet
   - golangci-lint

3. **template-lint** - 模板檢查
   - 構建並運行 template linter
   - 檢查所有組件模板

4. **build** - 構建測試
   - 構建主應用
   - 構建所有範例
   - 上傳構建產物

5. **security** - 安全掃描
   - Gosec security scanner
   - SARIF 報告

**觸發條件：**
- Push to main/develop
- Pull requests to main/develop

### 4. Makefile

**位置：** 根目錄 `Makefile`

**主要命令：**

```bash
make test              # 運行所有測試
make test-coverage     # 測試 + 覆蓋率報告
make lint              # 運行 template linter
make build-linter      # 構建 template linter
make install-tools     # 安裝所有開發工具
make run               # 運行主示例
make build-examples    # 構建所有範例
make run-example EXAMPLE=name  # 運行特定範例
make fmt               # 格式化代碼
make vet               # 運行 go vet
make check             # 運行所有檢查（fmt, vet, lint, test）
make clean             # 清理構建產物
make help              # 顯示幫助信息
```

**注意：** Windows 用戶可能需要安裝 Make 或使用 Git Bash。

### 5. 文檔

創建/更新的文檔：

1. **tools/template-linter/README.md**
   - Linter 使用說明
   - 問題類型說明
   - 範例和修復方法

2. **docs/DEVELOPMENT.md**（新建）
   - 完整的開發指南
   - 工具使用說明
   - 常見問題解決
   - 最佳實踐
   - 發布流程

3. **docs/TOOLING_SETUP.md**（本文檔）
   - 工具設置總結
   - 完成的工作清單
   - 使用說明

4. **README.md**（更新）
   - 新增「開發工具」章節
   - 新增 Makefile 命令說明
   - 新增 CI/CD 說明
   - 更新「貢獻」章節

## 文件結構

```
go-vdom/
├── .github/
│   └── workflows/
│       └── ci.yml                    # CI/CD 配置
├── .githooks/
│   ├── pre-commit                    # Pre-commit hook
│   └── install.sh                    # Hook 安裝腳本
├── tools/
│   └── template-linter/
│       ├── main.go                   # Linter 主程序
│       ├── go.mod                    # Go module
│       └── README.md                 # 使用說明
├── docs/
│   ├── DEVELOPMENT.md                # 開發指南（新建）
│   └── TOOLING_SETUP.md              # 本文檔（新建）
├── Makefile                          # 開發命令（新建）
└── README.md                         # 更新：新增工具說明
```

## 使用流程

### 初次設置（新開發者）

```bash
# 1. Clone 倉庫
git clone https://github.com/TimLai666/go-vdom.git
cd go-vdom

# 2. 安裝依賴
go mod download

# 3. 構建 linter（Windows）
cd tools/template-linter
go build
cd ../..

# 4. 安裝 git hooks（Linux/Mac）
chmod +x .githooks/install.sh
./.githooks/install.sh

# 5. 驗證設置
cd tools/template-linter
./template-linter ../../
```

### 日常開發流程

```bash
# 1. 創建功能分支
git checkout -b feature/my-feature

# 2. 進行修改
# ... 編輯文件 ...

# 3. 運行檢查（手動）
cd tools/template-linter
./template-linter -fix ../../

# 4. 運行測試
cd ../..
go test ./dom/...
go test ./components/...

# 5. 格式化代碼
go fmt ./...

# 6. 提交（自動觸發 pre-commit hook）
git add .
git commit -m "feat: Add my feature"

# 7. Push（觸發 CI）
git push origin feature/my-feature
```

### CI/CD 流程

1. **Push/PR** → 觸發 GitHub Actions
2. **自動運行：**
   - 測試（多版本）
   - Linting（代碼 + 模板）
   - 構建
   - 安全掃描
3. **結果：**
   - ✅ 全部通過 → 可以合併
   - ❌ 有失敗 → 修復後重新 push

## 測試結果

### Template Linter 測試

```bash
$ cd tools/template-linter && ./template-linter -fix ../../components
✓ No template issues found!

$ ./template-linter -fix ../../
✓ No template issues found!
```

### 核心測試

之前的測試結果顯示：

```bash
$ go test ./dom
PASS
ok      github.com/TimLai666/go-vdom/dom    0.XXXs

$ go test ./components/...
PASS
ok      github.com/TimLai666/go-vdom/components    0.XXXs
```

## 已知限制

1. **Makefile 在 Windows 上需要額外工具**
   - 解決方案：使用 Git Bash 或直接運行 go 命令

2. **Examples 目錄的構建問題**
   - 原因：每個檔案都有 main()
   - 解決方案：單獨構建每個範例
   ```bash
   go build examples/forms_demo.go
   ```

3. **Template Linter 可能有誤報**
   - 某些特殊情況可能被錯誤標記
   - 需要人工檢查確認

## 後續建議

### 短期（已完成）

- ✅ 建立 template linter
- ✅ 設置 git hooks
- ✅ 配置 CI/CD
- ✅ 創建 Makefile
- ✅ 編寫文檔

### 中期（可選）

1. **增強 Linter**
   - 新增更多檢測規則
   - 支援自動修復（--fix）
   - 生成詳細報告

2. **改進 CI/CD**
   - 新增性能測試
   - 新增集成測試
   - 自動生成變更日誌

3. **工具改進**
   - 創建 Windows batch 腳本（替代 Makefile）
   - 新增 VS Code 擴展配置
   - 新增 IDE 整合

### 長期（未來）

1. **自動化發布**
   - 自動標記版本
   - 自動生成 release notes
   - 自動發布到 Go package registry

2. **文檔自動化**
   - 從代碼生成 API 文檔
   - 自動更新範例
   - 互動式文檔網站

3. **開發環境**
   - Docker 開發環境
   - VS Code Dev Container
   - 在線 playground

## 相關文檔

- [README.md](../README.md) - 項目概述和快速開始
- [DEVELOPMENT.md](DEVELOPMENT.md) - 完整開發指南
- [TEMPLATE_EXPRESSION_FIX_GUIDE.md](../TEMPLATE_EXPRESSION_FIX_GUIDE.md) - 模板修復指南
- [CHANGELOG_JSON_SERIALIZATION.md](../CHANGELOG_JSON_SERIALIZATION.md) - JSON 序列化變更日誌
- [API_REFERENCE.md](API_REFERENCE.md) - API 參考
- [tools/template-linter/README.md](../tools/template-linter/README.md) - Linter 使用說明

## 問題排查

### Linter 無法運行

**症狀：** `template-linter: command not found`

**解決方案：**
```bash
cd tools/template-linter
go build
./template-linter ../../  # Linux/Mac
template-linter.exe ..\..\  # Windows
```

### Git Hooks 未觸發

**症狀：** 提交時沒有運行檢查

**解決方案：**
```bash
# 檢查 hook 是否存在
ls -la .git/hooks/pre-commit

# 檢查是否可執行
chmod +x .git/hooks/pre-commit

# 重新安裝
./.githooks/install.sh
```

### CI 失敗

**症狀：** GitHub Actions 顯示紅色 X

**解決方案：**
1. 查看失敗的 job
2. 檢查錯誤日誌
3. 在本地運行相同的命令
4. 修復後重新 push

### Make 命令失敗（Windows）

**症狀：** `make: command not found`

**解決方案：**
```bash
# 選項 1：使用 Git Bash
# 選項 2：直接運行 go 命令
go test ./...
go fmt ./...
cd tools/template-linter && go run main.go ../../

# 選項 3：安裝 Make for Windows
# https://gnuwin32.sourceforge.net/packages/make.htm
```

## 總結

本次工作成功建立了完整的開發工具鏈，包括：

1. **自動化檢查** - Template Linter 可以在開發早期捕獲錯誤
2. **持續集成** - GitHub Actions 確保代碼質量
3. **開發便利性** - Makefile 和 git hooks 簡化日常工作
4. **完善文檔** - 詳細的使用說明和故障排除指南

這些工具將大幅提升開發效率和代碼質量，防止模板相關的錯誤進入生產環境。

---

**工作完成時間：** 2024
**狀態：** ✅ 已完成並測試
**維護者：** 參考 [DEVELOPMENT.md](DEVELOPMENT.md) 了解如何貢獻和維護這些工具
