export interface MockUser {
	id: string;
	email: string;
	name: string;
	role: 'admin' | 'manager' | 'user';
	is_active: boolean;
	created_at: string;
}

export const mockUsers: MockUser[] = [
	{ id: 'u1', email: 'admin@gofin.id', name: 'Admin Gofin', role: 'admin', is_active: true, created_at: '2025-01-01' },
	{ id: 'u2', email: 'budi@gofin.id', name: 'Budi Santoso', role: 'manager', is_active: true, created_at: '2025-03-15' },
	{ id: 'u3', email: 'sari@gofin.id', name: 'Sari Dewi', role: 'user', is_active: true, created_at: '2025-06-20' },
	{ id: 'u4', email: 'andi@gofin.id', name: 'Andi Pratama', role: 'user', is_active: true, created_at: '2025-09-10' },
	{ id: 'u5', email: 'rina@gofin.id', name: 'Rina Wati', role: 'user', is_active: false, created_at: '2025-11-05' },
	{ id: 'u6', email: 'doni@gofin.id', name: 'Doni Saputra', role: 'user', is_active: true, created_at: '2026-01-18' }
];
