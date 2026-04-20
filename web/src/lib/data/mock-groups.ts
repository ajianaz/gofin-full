import type { UserGroup } from '$lib/types/index.js';

export const mockGroups: UserGroup[] = [
	{ id: 'g1', title: 'Keuangan Pribadi', member_count: 1, is_current: true },
	{ id: 'g2', title: 'Keluarga', member_count: 4, is_current: false },
	{ id: 'g3', title: 'Bisnis Side Project', member_count: 2, is_current: false }
];
