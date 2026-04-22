import { api } from './client.js';
import { unwrapMany, unwrapOne } from './helpers.js';
import type { RuleGroup, Rule } from '$lib/types/domain.js';

export const ruleService = {
	async list(): Promise<RuleGroup[]> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> }[] }>('/rules');
		return unwrapMany<RuleGroup>(res).map((r) => ({
			...r,
			stop_processing: false,
			rule_count: 0
		}));
	},

	async createGroup(data: { title: string }): Promise<RuleGroup> {
		const res = await api.post<{ data: { id: string; attributes: Record<string, unknown> } }>('/rule-groups', data);
		const r = unwrapOne<RuleGroup>(res);
		return { ...r, stop_processing: false, rule_count: 0 };
	},

	async get(id: string): Promise<Rule> {
		const res = await api.get<{ data: { id: string; attributes: Record<string, unknown> } }>(`/rules/${id}`);
		return unwrapOne<Rule>(res);
	},

	async updateGroup(id: string, data: { title?: string; active?: boolean }): Promise<void> {
		await api.put(`/rule-groups/${id}`, data);
	},

	async deleteGroup(id: string): Promise<void> {
		await api.delete(`/rule-groups/${id}`);
	}
};
