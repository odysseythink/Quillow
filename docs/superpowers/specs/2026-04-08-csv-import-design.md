# 账单 CSV 导入设计

**日期**: 2026-04-08
**项目**: Quillow

## 目标

支持导入微信支付和支付宝的账单 CSV 文件，自动解析交易记录，结合 AI 智能分类，批量创建交易。

## 支持的格式

### 微信账单 CSV

- 编码：UTF-8 BOM
- 前 16 行为账单概览信息，实际数据从第 17 行开始
- 列：交易时间, 交易类型, 交易对方, 商品, 收/支, 金额(元), 支付方式, 当前状态, 交易单号, 商户单号, 备注
- 金额格式：`¥35.00`（带 ¥ 前缀）
- 收/支字段：`支出`、`收入`、`/`（不计收支）

### 支付宝账单 CSV

- 编码：GBK
- 前 4 行为标题信息，实际数据从第 5 行开始
- 列：交易号, 商家订单号, 交易创建时间, 付款时间, 最近修改时间, 交易来源地, 类型, 交易对方, 商品名称, 金额（元）, 收/支, 交易状态, 服务费（元）, 成功退款（元）, 备注
- 金额格式：`35.00`（无前缀）
- 收/支字段：`支出`、`收入`、空

## 格式自动检测

读取文件前几行内容判断来源：
- 包含 `微信支付账单` → 微信
- 包含 `支付宝` 或首行包含 `交易号` → 支付宝
- 否则尝试通用 CSV（首行为 header，自动映射列名）

## 导入流程

```
1. 用户上传 CSV 文件
2. 后端自动检测来源（微信/支付宝/通用）
3. 解析所有交易行
4. 对每笔交易调用 Phase 1 智能分类（AI Suggest）
5. 返回预览列表
6. 用户勾选/取消、修改分类
7. 用户确认 → 批量创建交易
```

## API 端点

### `POST /api/v1/import/preview`

Content-Type: `multipart/form-data`

表单字段：
- `file`: CSV 文件

响应：
```json
{
  "source": "wechat",
  "total": 42,
  "transactions": [
    {
      "index": 0,
      "date": "2026-04-07",
      "type": "withdrawal",
      "description": "星巴克",
      "counterparty": "星巴克中国",
      "amount": "35.00",
      "category_id": 5,
      "category_name": "餐饮",
      "status": "支付成功",
      "external_id": "420000123456789",
      "selected": true
    }
  ]
}
```

`selected: false` 的情况：
- 交易状态不是"支付成功"/"交易成功"
- 收/支为"/"或空（不计收支的记录）
- external_id 已存在（重复导入）

### `POST /api/v1/import/confirm`

请求：
```json
{
  "transactions": [
    {
      "index": 0,
      "date": "2026-04-07",
      "type": "withdrawal",
      "description": "星巴克",
      "amount": "35.00",
      "category_id": 5,
      "external_id": "420000123456789"
    }
  ]
}
```

响应：
```json
{
  "imported": 40,
  "skipped": 2
}
```

## 去重策略

使用 `交易单号`（微信）/ `交易号`（支付宝）作为 `external_id`：
- preview 阶段：查询数据库中是否已存在相同 external_id 的交易，存在则标记 `selected: false`
- confirm 阶段：再次检查，跳过已存在的记录

## 前端

### 导入页面 `/import`

侧边栏新增"导入"菜单项。页面包含三步：

**Step 1 — 上传**
- 拖拽或点击上传区域（Ant Design Upload/Dragger）
- 支持 .csv 文件
- 上传后调用 preview API

**Step 2 — 预览**
- 表格展示解析结果
- 每行有复选框（默认根据 selected 字段）
- 分类列可内联编辑（下拉选择）
- 顶部显示来源标识（微信/支付宝）和总数

**Step 3 — 确认**
- 确认导入按钮
- 成功后显示导入数量
- 提供"查看导入的交易"链接

### i18n 新增 Key

```json
{
  "import": "导入",
  "import_bills": "导入账单",
  "upload_csv": "上传 CSV 文件",
  "upload_csv_hint": "支持微信支付、支付宝账单",
  "import_preview": "预览",
  "import_confirm": "确认导入",
  "import_success": "成功导入 {count} 笔交易",
  "import_source": "来源",
  "import_counterparty": "交易对方",
  "import_status": "状态",
  "import_duplicate": "已导入",
  "import_select_all": "全选",
  "import_deselect_all": "取消全选"
}
```

## 后端新增文件

```
pkg/importer/
├── importer.go       # 入口：Detect() 自动检测格式，Parse() 分发解析
├── wechat.go         # 微信 CSV 解析
├── alipay.go         # 支付宝 CSV 解析
├── types.go          # ImportedTransaction 结构定义

internal/adapter/handler/v1/
├── import_handler.go # Preview + Confirm 端点

internal/usecase/importer/
├── usecase.go        # 解析 + AI分类 + 去重检查 + 批量创建
```

### pkg/importer/types.go

```go
type ImportedTransaction struct {
    Index        int    `json:"index"`
    Date         string `json:"date"`
    Type         string `json:"type"`         // withdrawal, deposit, transfer
    Description  string `json:"description"`
    Counterparty string `json:"counterparty"`
    Amount       string `json:"amount"`
    Status       string `json:"status"`
    ExternalID   string `json:"external_id"`
}
```

### pkg/importer/wechat.go

```go
func ParseWechat(reader io.Reader) ([]ImportedTransaction, error)
```

逻辑：
1. 跳过前 16 行
2. 读取 header 行
3. 逐行解析，处理：
   - 金额：去掉 `¥` 前缀和空格
   - 日期：解析 `2026-04-07 12:34:56` 格式，取日期部分
   - 类型：`支出` → withdrawal, `收入` → deposit, 其他跳过
   - external_id：交易单号字段

### pkg/importer/alipay.go

```go
func ParseAlipay(reader io.Reader) ([]ImportedTransaction, error)
```

逻辑：
1. 文件编码从 GBK 转 UTF-8
2. 跳过前 4 行
3. 读取 header 行
4. 逐行解析，处理：
   - 金额：直接使用，去掉空格
   - 日期：解析 `2026-04-07 12:34:56` 格式
   - 类型：`支出` → withdrawal, `收入` → deposit
   - external_id：交易号字段
   - 跳过交易状态不是"交易成功"的行

## 路由注册

```go
// 在 router.go authenticated 路由组中
auth.POST("/import/preview", h.Import.Preview)
auth.POST("/import/confirm", h.Import.Confirm)
```

## 不在此阶段做的事

- 银行账单导入（各银行格式不统一）
- 邮件自动监听导入
- 截图 OCR 导入
- 实时同步（微信/支付宝无开放 API）
- 导入历史记录管理
