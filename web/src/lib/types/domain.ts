export interface Account {
	id: string;
	name: string;
	type: 'asset' | 'cash' | 'liability' | 'expense' | 'revenue';
	active: boolean;
	balance: string;
	currency_code: string;
	currency_symbol: string;
	currency_decimal_places: number;
	iban?: string;
}

export interface Transaction {
	id: string;
	type: 'withdrawal' | 'deposit' | 'transfer' | 'opening-balance';
	description: string;
	amount: string;
	date: string;
	source_account?: string;
	source_account_name?: string;
	destination_account?: string;
	destination_account_name?: string;
	category?: string;
	category_name?: string;
	budget?: string;
	budget_name?: string;
	tags: string[];
	currency_code: string;
	currency_symbol: string;
	currency_decimal_places: number;
	note?: string;
}

export interface Budget {
	id: string;
	name: string;
	active: boolean;
	order: number;
	auto_budget_type?: 'none' | 'reset' | 'rollover' | 'fixed' | 'adjust';
	auto_budget_period?: 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'yearly';
	spend_amount: string;
	budget_amount: string;
	limits: BudgetLimit[];
}

export interface BudgetLimit {
	id: string;
	budget_id: string;
	start_date: string;
	end_date: string;
	amount: string;
	spend: string;
	category?: string;
	category_name?: string;
}

export interface PiggyBank {
	id: string;
	name: string;
	account_id: string;
	account_name: string;
	target_amount: string;
	current_amount: string;
	start_date: string;
	target_date?: string;
	status: 'active' | 'completed' | 'cancelled';
}

export interface Bill {
	id: string;
	name: string;
	amount_min: string;
	amount_max: string;
	next_date: string;
	repeat_freq: 'weekly' | 'monthly' | 'quarterly' | 'yearly';
	active: boolean;
	currency_code: string;
	currency_symbol: string;
}

export interface RecurringTransaction {
	id: string;
	title: string;
	type: 'withdrawal' | 'deposit' | 'transfer';
	first_date: string;
	repeat_freq: 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'yearly';
	repeat_until?: string;
	active: boolean;
	amount: string;
	currency_code: string;
	description?: string;
	source_account_id?: string;
	source_account_name?: string;
	destination_account_id?: string;
	destination_account_name?: string;
	category_id?: string;
	category_name?: string;
}

export interface Category {
	id: string;
	name: string;
	type: 'expense' | 'income' | 'transfer';
	icon?: string;
	transaction_count: number;
}

export interface Tag {
	id: string;
	tag: string;
	date: string;
	description?: string;
}

export interface RuleGroup {
	id: string;
	title: string;
	description?: string;
	active: boolean;
	stop_processing: boolean;
	order: number;
	rule_count: number;
}

export interface Rule {
	id: string;
	rule_group_id: string;
	title: string;
	active: boolean;
	strict: boolean;
	stop_processing: boolean;
	priority: number;
	trigger_type: string;
	trigger_value?: string;
	action_type: string;
	action_value?: string;
}

export interface Currency {
	id: string;
	code: string;
	name: string;
	symbol: string;
	decimal_places: number;
	enabled: boolean;
}

export interface ExchangeRate {
	id: string;
	from_code: string;
	to_code: string;
	rate: number;
	date: string;
}

export interface UserGroup {
	id: string;
	title: string;
	member_count: number;
	is_current: boolean;
}

export interface AuditLogEntry {
	id: string;
	action: string;
	user_email: string;
	entity_type: string;
	entity_id: string;
	changes: string;
	created_at: string;
}

export interface ApiKey {
	id: string;
	title: string;
	key: string;
	created_at: string;
	last_used_at?: string;
}

export interface Preference {
	name: string;
	value: string | boolean | number;
	description: string;
	type: 'text' | 'boolean' | 'select' | 'number';
	options?: string[];
}

export interface NotificationSetting {
	id: string;
	title: string;
	description: string;
	enabled: boolean;
}
