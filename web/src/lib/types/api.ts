// Auto-generated from OpenAPI spec — do not edit manually
// Generated: 2026-06-01

export interface AttachmentRequest {
	filename: string;
	mime_type: string;
	/** Base64-encoded file content */
	content?: string;
}

export interface BillRequest {
	name: string;
	amount_min: string;
	amount_max?: string;
	currency_code: string;
	date: string;
	end_date?: string;
	repeat_freq?: string;
	skip?: number;
	notes?: string;
}

export interface BudgetRequest {
	name: string;
	start_date: string;
	end_date?: string;
	amount: string;
	currency_code: string;
}

export interface CategoryRequest {
	name: string;
	parent_id?: string;
}

export interface ConfigurationRequest {
	name: string;
	value: string;
}

export interface Error {
	status?: number;
	title?: string;
	detail?: string;
}

export interface ExchangeRateRequest {
	from_currency: string;
	to_currency: string;
	rate: number;
	date: string;
}

export interface FeatureFlagRequest {
	name: string;
	enabled: boolean;
}

export interface GroupRequest {
	name: string;
	description?: string;
}

export interface LocationRequest {
	latitude?: number;
	longitude?: number;
	zoom_level?: number;
	locatable_id?: string;
	locatable_type?: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}

export interface NoteRequest {
	text: string;
	noteable_id?: string;
	noteable_type?: string;
}

export interface ObjectGroupRequest {
	name: string;
	order?: number;
}

export interface PiggyBankMoneyRequest {
	/** Amount to add or remove */
	amount: string;
}

export interface PiggyBankRequest {
	name: string;
	target_amount: string;
	start_date?: string;
	target_date?: string;
}

export interface PreferenceRequest {
	name: string;
	value: string;
}

export interface RecurrenceRequest {
	title: string;
	type: "withdrawal" | "deposit" | "transfer";
	description?: string;
	repetition_type: "daily" | "weekly" | "monthly" | "yearly";
	repetition_end: string;
	apply_rules?: boolean;
	active?: boolean;
}

export interface RefreshRequest {
	refresh_token: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	name: string;
}

export interface RuleGroupRequest {
	name: string;
	description?: string;
	active?: boolean;
}

export interface RuleRequest {
	title: string;
	rule_group_id: string;
	description?: string;
	active?: boolean;
	triggers: Record<string, unknown>[];
	actions: Record<string, unknown>[];
}

export interface SplitTransactionRequest {
	transactions: TransactionRequest[];
}

export interface TagRequest {
	tag: string;
	date?: string;
	description?: string;
}

export interface TokenResponse {
	access_token?: string;
	refresh_token?: string;
	expires_in?: number;
	token_type?: string;
}

export interface TransactionRequest {
	type: "withdrawal" | "deposit" | "transfer";
	date: string;
	amount: string;
	source_id: string;
	destination_id?: string;
	category_id?: string;
	description?: string;
	tags?: string[];
	notes?: string;
}

export interface UserResponse {
	id?: string;
	email?: string;
	name?: string;
	created_at?: string;
}

export interface WalletMemberRequest {
	user_id: string;
	role: "owner" | "manage_meta" | "manage_transactions" | "manage_budgets" | "manage_piggy_banks" | "manage_rules" | "manage_recurring" | "manage_currencies" | "manage_webhooks" | "view_memberships" | "view_reports";
}

export interface WalletRequest {
	name: string;
	/** Wallet type code */
	type: string;
	currency_code: string;
	balance?: number;
	virtual_balance?: number;
	active?: boolean;
}

export interface WebhookRequest {
	url: string;
	trigger: "STORE_TRANSACTION" | "UPDATE_TRANSACTION" | "DESTROY_TRANSACTION";
	active?: boolean;
}
