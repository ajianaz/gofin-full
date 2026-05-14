import { describe, it, expect, vi, beforeEach } from 'vitest';
import { adminService } from './admin.js';

vi.mock('./client.js', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn()
  }
}));

import { api } from './client.js';

const mockedApi = vi.mocked(api);

describe('adminService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('listUsers', () => {
    it('unwraps and maps user fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'u1', attributes: { email: 'admin@test.com', name: 'Admin', role: 'admin', is_active: true, created_at: '2026-01-01' } },
          { id: 'u2', attributes: { email: 'user@test.com' } }
        ]
      });

      const result = await adminService.listUsers();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'u1',
          email: 'admin@test.com',
          name: 'Admin',
          role: 'admin',
          is_active: true,
          created_at: '2026-01-01'
        })
      );
      // Defaults for missing fields
      expect(result[1]).toEqual(
        expect.objectContaining({
          id: 'u2',
          email: 'user@test.com',
          name: '',
          role: 'user',
          is_active: true,
          created_at: ''
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/admin/users');
    });

    it('returns empty array for no users', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await adminService.listUsers();

      expect(result).toEqual([]);
    });

    it('throws when api fails', async () => {
      mockedApi.get.mockRejectedValueOnce(new Error('Forbidden'));

      await expect(adminService.listUsers()).rejects.toThrow('Forbidden');
    });
  });

  describe('listAuditLogs', () => {
    it('unwraps and maps audit log fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 'a1',
            attributes: {
              action: 'create',
              user_id: 'u1',
              entity_type: 'transaction',
              entity_id: 't1',
              new_value: '{"amount":"100"}',
              created_at: '2026-01-15T10:00:00Z'
            }
          }
        ]
      });

      const result = await adminService.listAuditLogs();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual({
        id: 'a1',
        action: 'create',
        user_email: 'u1',
        entity_type: 'transaction',
        entity_id: 't1',
        changes: '{"amount":"100"}',
        created_at: '2026-01-15T10:00:00Z'
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/audit-logs');
    });

    it('passes entity_type query param', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      await adminService.listAuditLogs('transaction');

      expect(mockedApi.get).toHaveBeenCalledWith('/audit-logs?entity_type=transaction');
    });

    it('applies defaults for all fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: 'a2', attributes: {} }]
      });

      const result = await adminService.listAuditLogs();

      expect(result[0]).toEqual({
        id: 'a2',
        action: '',
        user_email: '',
        entity_type: '',
        entity_id: '',
        changes: '',
        created_at: ''
      });
    });

    it('prefers user_id over user_email and new_value over changes', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{
          id: 'a3',
          attributes: {
            user_email: 'admin@test.com',
            user_id: 'u1',
            changes: 'updated name',
            new_value: '{"name":"old"}',
            entity_id: 'e1'
          }
        }]
      });

      const result = await adminService.listAuditLogs();

      // user_id takes precedence over user_email
      expect(result[0].user_email).toBe('u1');
      // new_value takes precedence over changes
      expect(result[0].changes).toBe('{"name":"old"}');
    });
  });
});
