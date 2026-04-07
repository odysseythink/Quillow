# 财务洞察助手设计

**日期**: 2026-04-08
**项目**: Quillow
**阶段**: Phase 3 of 3 (智能分类 → 自然语言记账 → **财务洞察助手**)

## 目标

用户在聊天气泡中用自然语言提问（如"这个月餐饮花了多少"），AI 查询数据并以自然语言回答，支持查询和同比/环比分析。

## 交互方式

复用 Phase 2 的聊天气泡，统一入口。AI 根据意图自动区分记账和查询：
- 包含金额/记账动词 → Phase 2 自然语言记账
- 包含疑问词/查询意图 → 本 Phase 洞察查询

### 统一聊天 API

`POST /api/v1/ai/chat`

前端只调这一个端点，后端内部判断意图并分发。

请求：
```json
{"message": "这个月比上个月多花了多少"}
```

响应（查询意图）：
```json
{
  "intent": "query",
  "answer": "这个月总支出 ¥5,320.00，上个月 ¥4,180.00，多花了 ¥1,140.00（+27.3%）。",
  "data": {
    "function": "monthly_comparison",
    "result": { "month1_total": "5320.00", "month2_total": "4180.00" }
  }
}
```

响应（记账意图）：
```json
{
  "intent": "record",
  "parsed": {
    "type": "withdrawal",
    "description": "午饭",
    "amount": "35",
    "date": "2026-04-08",
    "category": "餐饮"
  },
  "confidence": "high",
  "created": true,
  "transaction_id": 42
}
```

## 意图识别

在 `pkg/ai/nlp.go` 中新增 `DetectIntent(message) → "record" | "query"`。

### 本地规则

| 意图 | 匹配规则 |
|---|---|
| query | 包含：多少、花了、支出、收入、趋势、最大、最小、平均、比较、哪个、几笔、统计、分析、总共、合计、余额、预算 |
| record | 包含金额模式（`\d+元`、`¥\d+`）且不含疑问词 |
| 不确定 | 交给 LLM 判断 |

## 查询执行架构

```
用户提问
  → 意图识别（本地规则 / LLM）
  → LLM 第 1 次调用：选择查询函数 + 填参数
  → 后端执行安全查询
  → LLM 第 2 次调用：将结果组织为自然语言
  → 返回给用户
```

**安全原则：不让 LLM 直接写 SQL。** LLM 只从预定义的查询函数中选择，后端执行。

## 预定义查询函数

| 函数名 | 说明 | 参数 |
|---|---|---|
| `total_spending` | 时间段内总支出 | start, end |
| `total_income` | 时间段内总收入 | start, end |
| `sum_by_category` | 按分类汇总支出/收入 | start, end, type |
| `top_transactions` | 最大/最小 N 笔交易 | start, end, type, order(desc/asc), limit |
| `daily_trend` | 每日支出/收入趋势 | start, end, type |
| `monthly_comparison` | 两个月的总额对比 | month1(YYYY-MM), month2(YYYY-MM), type |
| `category_comparison` | 某分类的月度对比 | category, month1, month2 |
| `account_balance` | 账户当前余额 | account_name |
| `budget_status` | 预算使用情况 | budget_name, month |

### 函数签名

```go
type QueryFunction struct {
    Name        string
    Description string
    Execute     func(ctx context.Context, db *gorm.DB, userID uint, params map[string]string) (any, error)
}
```

所有函数返回 JSON-serializable 结果，传给 LLM 生成自然语言回答。

## LLM Prompt

### 第 1 次调用 — 解析 + 选函数

```
你是一个财务数据助手。根据用户问题，选择合适的查询函数并填写参数。

用户问题: "{message}"
今天日期: {today}
可用查询函数:
- total_spending(start, end): 时间段内总支出
- total_income(start, end): 时间段内总收入
- sum_by_category(start, end, type): 按分类汇总
- top_transactions(start, end, type, order, limit): 最大/最小N笔
- daily_trend(start, end, type): 每日趋势
- monthly_comparison(month1, month2, type): 月度对比
- category_comparison(category, month1, month2): 分类月度对比
- account_balance(account_name): 账户余额
- budget_status(budget_name, month): 预算状态

可用分类: {categories}
可用账户: {accounts}

返回 JSON: {"function": "函数名", "params": {"key": "value"}}
仅返回 JSON。
```

### 第 2 次调用 — 生成回答

```
用户问题: "{message}"
查询结果: {query result JSON}

用简洁友好的中文回答用户问题。如果涉及对比，指出变化金额和百分比。不要编造数据。
```

### 错误处理

- LLM 选择了不存在的函数 → 返回"暂不支持该查询"
- 查询结果为空 → LLM 回答"该时间段没有相关记录"
- LLM 调用失败 → 返回"服务暂时不可用，请稍后再试"

## 后端新增/修改文件

```
pkg/ai/
├── nlp.go              # 修改：加入 DetectIntent()
├── insight.go          # 新增：QueryFunction 注册表 + 执行器

internal/usecase/ai/
├── usecase.go          # 修改：加入 Chat() 统一入口、Insight() 方法

internal/adapter/handler/v1/
├── ai_handler.go       # 修改：加入 Chat 统一端点
```

### pkg/ai/insight.go 结构

```go
// QueryRegistry holds all available query functions.
type QueryRegistry struct {
    functions map[string]QueryFunction
}

func NewQueryRegistry(db *gorm.DB) *QueryRegistry
func (r *QueryRegistry) Execute(ctx, userID, functionName, params) (any, error)
func (r *QueryRegistry) Describe() string  // 返回函数列表描述，供 LLM prompt 使用
```

### usecase 新增方法

```go
// Chat 统一入口，自动判断意图分发
func (uc *UseCase) Chat(ctx, userID, message) → ChatResponse

// Insight 处理查询意图
func (uc *UseCase) Insight(ctx, userID, message) → (answer string, data any, error)
```

## 前端修改

### ChatBubble.tsx

修改发送逻辑，统一调用 `POST /api/v1/ai/chat`，根据响应 `intent` 字段渲染：
- `intent: "record"` → 同 Phase 2（已记录消息 / 预览卡片）
- `intent: "query"` → 显示 AI 文本回答（`answer` 字段）

### i18n 新增 Key

```json
{
  "chat_thinking": "正在查询...",
  "chat_query_failed": "查询失败，请稍后再试"
}
```

## API 端点汇总

| 端点 | 说明 | Phase |
|---|---|---|
| `POST /api/v1/ai/suggest` | 智能分类建议 | Phase 1 |
| `POST /api/v1/ai/classify-batch` | 批量分类 | Phase 1 |
| `POST /api/v1/ai/learn` | 学习分类模式 | Phase 1 |
| `POST /api/v1/ai/parse-transaction` | 自然语言解析交易 | Phase 2 |
| `POST /api/v1/ai/insight` | 财务洞察查询 | Phase 3 |
| `POST /api/v1/ai/chat` | **统一聊天入口** | Phase 3 |

Phase 3 完成后，前端聊天气泡统一使用 `/ai/chat`，其他端点仍保留供直接调用。

## 不在此阶段做的事

- 图表可视化（聊天中嵌入图表）
- 定时推送洞察通知
- 多轮对话上下文记忆
- 理财建议/预测
