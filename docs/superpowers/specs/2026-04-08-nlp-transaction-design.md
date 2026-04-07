# 自然语言记账设计

**日期**: 2026-04-08
**项目**: Quillow
**阶段**: Phase 2 of 3 (智能分类 → **自然语言记账** → 财务洞察助手)

## 目标

用户在聊天气泡中输入自然语言（如"昨天午饭35元"），AI 解析并自动创建交易。高置信度直接创建，低置信度显示预览卡片让用户确认。

## 交互流程

页面右下角浮动聊天气泡，点开对话窗口。

1. 用户输入自然语言
2. 调用 `POST /api/v1/ai/parse-transaction`
3. 后端执行解析（本地正则 → LLM fallback）
4. 根据置信度决定行为：
   - **high**（必填字段齐全：description, amount, date, type）→ 直接创建交易，返回 `created: true`，聊天中显示"已记录"+ 撤销按钮
   - **low**（缺少字段或不确定）→ 返回 `created: false`，前端展示可编辑预览卡片，用户修改后点确认创建

## AI 解析

### 3 级 Fallback

```
1. 本地正则解析（nlp.go）— 提取金额、日期、类型关键词
2. 智能分类（Phase 1 已有）— 根据描述匹配分类
3. LLM API — 完整的自然语言理解，返回结构化 JSON
```

### 本地正则规则（nlp.go）

| 字段 | 匹配规则 |
|---|---|
| 金额 | `(\d+\.?\d*)\s*元`, `¥(\d+\.?\d*)`, `(\d+\.?\d*)\s*(块\|刀)` |
| 日期 | `今天` → today, `昨天` → yesterday, `前天` → day-before-yesterday, `\d{1,2}月\d{1,2}日`, `\d{4}-\d{2}-\d{2}` |
| 类型 | 默认 `withdrawal`；出现 `收入/工资/奖金/转入` → `deposit`；出现 `转账/转给` → `transfer` |
| 描述 | 去掉金额和日期后的剩余文本 |

置信度判定：能提取金额 + 描述 → `high`，否则 → `low`。

### LLM Prompt

```
你是一个记账助手。解析用户输入，提取交易信息。

用户输入: "{message}"
可用分类: {categories}
今天日期: {today}

返回 JSON:
{
  "type": "withdrawal|deposit|transfer",
  "description": "交易描述",
  "amount": "金额数字",
  "date": "YYYY-MM-DD",
  "category": "分类名",
  "source_name": "来源账户（可选）",
  "destination_name": "目标账户（可选）"
}
仅返回 JSON。
```

## API 端点

### `POST /api/v1/ai/parse-transaction`

请求：
```json
{"message": "昨天午饭35元"}
```

响应（高置信度，已创建）：
```json
{
  "parsed": {
    "type": "withdrawal",
    "description": "午饭",
    "amount": "35",
    "date": "2026-04-07",
    "category": "餐饮",
    "category_id": 5,
    "source_name": "",
    "destination_name": ""
  },
  "confidence": "high",
  "created": true,
  "transaction_id": 42
}
```

响应（低置信度，未创建）：
```json
{
  "parsed": {
    "type": "withdrawal",
    "description": "买了个东西",
    "amount": "",
    "date": "2026-04-08",
    "category": "",
    "category_id": 0,
    "source_name": "",
    "destination_name": ""
  },
  "confidence": "low",
  "created": false,
  "transaction_id": 0
}
```

### 确认创建（复用已有端点）

前端预览卡片确认后，调用已有的 `POST /api/v1/transactions` 创建。

### 撤销（复用已有端点）

调用 `DELETE /api/v1/transactions/:id`。

## 前端组件

### ChatBubble.tsx

新增 `web/src/components/ChatBubble.tsx`，嵌入 `AppLayout.tsx`。

**结构：**
- 浮动按钮：右下角圆形按钮，聊天图标
- 聊天面板：展开后 300px 宽，400px 高，固定定位
- 消息列表：滚动区域
- 输入框：底部，回车发送

**消息类型：**
1. **用户消息** — 右对齐气泡
2. **AI 文本回复** — 左对齐气泡（"已记录：餐饮 ¥35.00 午饭"）
3. **交易预览卡片** — 左对齐卡片，包含表单字段（type, description, amount, date, category），底部有确认/取消按钮
4. **撤销按钮** — 在"已记录"消息下方

**状态管理：** 组件内部 useState，不需要 Redux（聊天记录不需要全局共享）。

### i18n 新增 Key

```json
{
  "chat_placeholder": "输入记账内容...",
  "chat_recorded": "已记录",
  "chat_confirm": "确认记账",
  "chat_undo": "撤销",
  "chat_parse_failed": "无法识别，请重新输入"
}
```

## 后端新增/修改文件

```
pkg/ai/
├── nlp.go                    # 新增：本地正则解析

internal/adapter/handler/v1/
├── ai_handler.go             # 修改：加入 ParseTransaction 方法

internal/usecase/ai/
├── usecase.go                # 修改：加入 ParseAndCreate 方法
```

## 代码结构

### pkg/ai/nlp.go

```go
type ParsedTransaction struct {
    Type            string `json:"type"`
    Description     string `json:"description"`
    Amount          string `json:"amount"`
    Date            string `json:"date"`
    Category        string `json:"category"`
    SourceName      string `json:"source_name"`
    DestinationName string `json:"destination_name"`
}

func ParseLocal(message string, today string) (*ParsedTransaction, string)
// 返回 (解析结果, "high"|"low")
```

### usecase 新增方法

```go
func (uc *UseCase) ParseAndCreate(ctx, userID, message) → (parsed, confidence, created, txID, error)
```

逻辑：
1. `ParseLocal()` 提取金额/日期/描述
2. 如果有描述 → 调用 Phase 1 的 `Suggest()` 获取分类
3. 如果置信度 high → 调用交易创建逻辑
4. 返回结果

## 不在此阶段做的事

- 多轮对话（"帮我记一笔" → "什么？" → "午饭35"）
- 语音输入
- 图片/OCR 识别
- 财务洞察对话（Phase 3）
