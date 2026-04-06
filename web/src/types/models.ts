export interface User {
  id: string;
  email: string;
  blocked: boolean;
  blocked_code: string;
  role: string;
  created_at: string;
  updated_at: string;
}

export interface Account {
  id: string;
  name: string;
  type: string;
  active: boolean;
  order: number;
  account_role: string;
  currency_id: string;
  currency_code: string;
  currency_symbol: string;
  currency_decimal_places: number;
  current_balance: string;
  virtual_balance: string;
  iban: string;
  account_number: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface TransactionSplit {
  type: string;
  description: string;
  date: string;
  amount: string;
  foreign_amount: string;
  currency_code: string;
  currency_symbol: string;
  source_id: string;
  source_name: string;
  source_type: string;
  destination_id: string;
  destination_name: string;
  destination_type: string;
  budget_id: string;
  budget_name: string;
  category_id: string;
  category_name: string;
  bill_id: string;
  bill_name: string;
  tags: string[];
  notes: string;
  reconciled: boolean;
  external_id: string;
  internal_reference: string;
}

export interface TransactionGroup {
  id: string;
  group_title: string;
  transactions: TransactionSplit[];
  created_at: string;
  updated_at: string;
}

export interface Currency {
  id: string;
  name: string;
  code: string;
  symbol: string;
  decimal_places: number;
  enabled: boolean;
  primary: boolean;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  name: string;
  active: boolean;
  order: number;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface BudgetLimit {
  id: string;
  budget_id: string;
  transaction_currency_id: string;
  amount: string;
  start: string;
  end: string;
  period: string;
}

export interface Bill {
  id: string;
  name: string;
  amount_min: string;
  amount_max: string;
  date: string;
  repeat_freq: string;
  active: boolean;
  order: number;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  name: string;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface Tag {
  id: string;
  tag: string;
  description: string;
  date: string | null;
  created_at: string;
  updated_at: string;
}

export interface PiggyBank {
  id: string;
  account_id: string;
  name: string;
  target_amount: string;
  start_date: string | null;
  target_date: string | null;
  order: number;
  active: boolean;
  notes: string;
  created_at: string;
  updated_at: string;
}

export interface Rule {
  id: string;
  title: string;
  description: string;
  rule_group_id: string;
  order: number;
  active: boolean;
  strict: boolean;
  stop_processing: boolean;
  triggers: RuleTrigger[];
  actions: RuleAction[];
  created_at: string;
  updated_at: string;
}

export interface RuleTrigger {
  id: string;
  type: string;
  value: string;
  order: number;
  active: boolean;
  stop_processing: boolean;
}

export interface RuleAction {
  id: string;
  type: string;
  value: string;
  order: number;
  active: boolean;
  stop_processing: boolean;
}

export interface RuleGroup {
  id: string;
  title: string;
  description: string;
  order: number;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Webhook {
  id: string;
  active: boolean;
  title: string;
  trigger: number;
  response: number;
  delivery: number;
  url: string;
  created_at: string;
  updated_at: string;
}

export interface Preference {
  id: string;
  name: string;
  data: any;
  created_at: string;
  updated_at: string;
}

export interface Recurrence {
  id: string;
  title: string;
  description: string;
  first_date: string;
  repeat_until: string | null;
  active: boolean;
  apply_rules: boolean;
  created_at: string;
  updated_at: string;
}
