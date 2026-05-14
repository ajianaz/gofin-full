import { describe, it, expect, vi, beforeEach } from 'vitest';
import { preferenceService } from './preferences.js';

vi.mock('./client.js', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn()
  }
}));

import { api } from './client.js';
const mockedApi = vi.mocked(api);

describe('preferenceService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('maps JSON:API response to PreferenceItem array', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'p1', attributes: { name: 'language', data: 'id' } },
          { id: 'p2', attributes: { name: 'budget_indicator', data: 'true' } }
        ]
      });

      const result = await preferenceService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({ id: 'p1', name: 'language', data: 'id' });
      expect(result[1]).toEqual({ id: 'p2', name: 'budget_indicator', data: 'true' });
      expect(mockedApi.get).toHaveBeenCalledWith('/preferences');
    });

    it('returns empty array when no preferences', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await preferenceService.list();

      expect(result).toEqual([]);
    });
  });

  describe('get', () => {
    it('gets single preference by name', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: { id: 'p1', attributes: { name: 'language', data: 'en' } }
      });

      const result = await preferenceService.get('language');

      expect(result).toEqual({ id: 'p1', name: 'language', data: 'en' });
      expect(mockedApi.get).toHaveBeenCalledWith('/preferences/language');
    });
  });

  describe('set', () => {
    it('posts name and data', async () => {
      mockedApi.post.mockResolvedValueOnce({});

      await preferenceService.set('language', 'en');

      expect(mockedApi.post).toHaveBeenCalledWith('/preferences', { name: 'language', data: 'en' });
    });
  });

  describe('delete', () => {
    it('deletes preference by name', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await preferenceService.delete('language');

      expect(mockedApi.delete).toHaveBeenCalledWith('/preferences/language');
    });
  });
});
