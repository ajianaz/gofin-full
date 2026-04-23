import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { RuleGroup, Rule } from '$lib/types/domain.js';

export const ruleService = {
	async listGroups(): Promise<RuleGroup[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/rule-groups');
		return unwrapMany<RuleGroup>(res).map((r) => ({
			...r,
			title: (r as any).title || '',
			active: !!((r as any).active),
			order: Number((r as any).order) || 0,
			stop_processing: !!((r as any).stop_processing),
			rule_count: Number((r as any).rule_count) || 0
		}));
	},

	async createGroup(data: { title: string }): Promise<RuleGroup> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/rule-groups', data);
		const r = unwrapOne<RuleGroup>(res);
		return { ...r, stop_processing: false, rule_count: 0 };
	},

	async getGroup(id: string): Promise<RuleGroup> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> } }>(`/rule-groups/${id}`);
		const r = unwrapOne<RuleGroup>(res);
		return { ...r, stop_processing: !!((r as any).stop_processing), rule_count: 0 };
	},

	async updateGroup(id: string, data: { title?: string; active?: boolean }): Promise<void> {
		await api.put(`/rule-groups/${id}`, data);
	},

	async deleteGroup(id: string): Promise<void> {
		await api.delete(`/rule-groups/${id}`);
	},

	async listRules(): Promise<Rule[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/rules');
		return unwrapMany<Rule>(res).map((r) => ({
			...r,
			title: (r as any).title || '',
			active: !!((r as any).active),
			strict: !!((r as any).strict),
			stop_processing: !!((r as any).stop_processing),
			priority: Number((r as any).priority) || 0,
			trigger_type: (r as any).trigger_type || '',
			action_type: (r as any).action_type || ''
		}));
	},

	async getRule(id: string): Promise<Rule> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> } }>(`/rules/${id}`);
		const r = unwrapOne<Rule>(res);
		return { ...r, active: !!r.active, strict: !!r.strict, stop_processing: !!r.stop_processing };
	},

	async createRule(data: {
		title: string;
		priority?: number;
		rule_group_id?: string;
		triggers?: Array<{ trigger_type: string; trigger_value: string }>;
		actions?: Array<{ action_type: string; action_value: string }>;
	}): Promise<Rule> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/rules', data);
		return unwrapOne<Rule>(res);
	},

	async updateRule(id: string, data: {
		title?: string;
		active?: boolean;
		strict?: boolean;
		stop_processing?: boolean;
		triggers?: Array<{ trigger_type: string; trigger_value: string }>;
		actions?: Array<{ action_type: string; action_value: string }>;
	}): Promise<void> {
		await api.put(`/rules/${id}`, data);
	},

	async deleteRule(id: string): Promise<void> {
		await api.delete(`/rules/${id}`);
	}
};
