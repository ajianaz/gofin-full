import { describe, it, expect, vi, beforeEach } from 'vitest';
import { apiKeyService } from './api-keys.js';

vi.mock('./client.js', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
  }
}));

import { api } from './client.js';
const mockedApi = vi.mocked(api);

describe('apiKeyService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('maps API response to ApiKeyListItem array', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'ak1', name: 'Mobile App', key_prefix: 'gofin_mob', last_used: '2026-04-15T09:00:00Z', created_at: '2026-04-10T14:00:00Z' },
          { id: 'ak2', name: 'Script', key_prefix: 'gofin_scr', last_used: '', created_at: '2026-03-15T10:00:00Z' }
        ]
      });

      const result = await apiKeyService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({ id: 'ak1', name: 'Mobile App', key_prefix: 'gofin_mob', last_used: '2026-04-15T09:00:00Z', created_at: '2026-04-10T14:00:00Z' });
      expect(result[1].last_used).toBe('');
      expect(mockedApi.get).toHaveBeenCalledWith('/api-keys');
    });

    it('handles empty data array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await apiKeyService.list();

      expect(result).toEqual([]);
    });
  });

  describe('create', () => {
    it('returns created API key with raw key', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'ak3', name: 'Test', key_prefix: 'gofin_tst', key: 'gofin_tst_rawkey123', created_at: '2026-04-20T10:00:00Z' }
      });

      const result = await apiKeyService.create('Test');

      expect(result).toEqual({ id: 'ak3', name: 'Test', key_prefix: 'gofin_tst', key: 'gofin_tst_rawkey123', created_at: '2026-04-20T10:00:00Z' });
      expect(mockedApi.post).toHaveBeenCalledWith('/api-keys', { name: 'Test' });
    });
  });

  describe('delete', () => {
    it('calls delete with correct id', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await apiKeyService.delete('ak1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/api-keys/ak1');
    });
  });
});
