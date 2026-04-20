import type { RuleGroup, Rule } from '$lib/types/index.js';

export const mockRuleGroups: RuleGroup[] = [
	{ id: 'rg1', title: 'Auto-kategorisasi', description: 'Aturkan kategori transaksi secara otomatis', active: true, stop_processing: false, order: 1, rule_count: 3 },
	{ id: 'rg2', title: 'Pindah ke Tabungan', description: 'Otomatis pindah sebagian ke rekening tabungan', active: true, stop_processing: true, order: 2, rule_count: 2 },
	{ id: 'rg3', title: 'Tag Otomatis', description: 'Tambahkan tag berdasarkan deskripsi', active: false, stop_processing: false, order: 3, rule_count: 2 }
];

export const mockRules: Rule[] = [
	{ id: 'r1', rule_group_id: 'rg1', title: 'Makanan > Warteg', active: true, strict: true, stop_processing: false, priority: 1, trigger_type: 'description_contains', trigger_value: 'warteg', action_type: 'set_category', action_value: 'Makanan' },
	{ id: 'r2', rule_group_id: 'rg1', title: 'Grab > Makanan', active: true, strict: false, stop_processing: false, priority: 2, trigger_type: 'description_contains', trigger_value: 'grab food', action_type: 'set_category', action_value: 'Makanan' },
	{ id: 'r3', rule_group_id: 'rg1', title: 'PLN > Utilitas', active: true, strict: true, stop_processing: false, priority: 3, trigger_type: 'description_contains', trigger_value: 'pln', action_type: 'set_category', action_value: 'Utilitas' },
	{ id: 'r4', rule_group_id: 'rg2', title: 'Nabung 10%', active: true, strict: true, stop_processing: true, priority: 1, trigger_type: 'deposit', trigger_value: '', action_type: 'move_to_account', action_value: '2' },
	{ id: 'r5', rule_group_id: 'rg2', title: 'Nabung 5% default', active: true, strict: false, stop_processing: true, priority: 2, trigger_type: 'deposit', trigger_value: '', action_type: 'move_to_account', action_value: '6' },
	{ id: 'r6', rule_group_id: 'rg3', title: 'Tag #freelance', active: false, strict: false, stop_processing: false, priority: 1, trigger_type: 'description_contains', trigger_value: 'freelance', action_type: 'add_tag', action_value: 'freelance' },
	{ id: 'r7', rule_group_id: 'rg3', title: 'Tag #investasi', active: false, strict: false, stop_processing: false, priority: 2, trigger_type: 'description_contains', trigger_value: 'saham', action_type: 'add_tag', action_value: 'investasi' }
];
