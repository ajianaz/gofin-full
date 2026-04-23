import { describe, it, expect, vi, beforeEach } from 'vitest';
import { groupService } from './groups.js';

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

describe('groupService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps group fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'g1', attributes: { title: 'Personal', member_count: 3, is_current: true } },
          { id: 'g2', attributes: { title: 'Business' } }
        ]
      });

      const result = await groupService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({
        id: 'g1',
        title: 'Personal',
        member_count: 3,
        is_current: true
      });
      expect(result[1]).toEqual({
        id: 'g2',
        title: 'Business',
        member_count: 0,
        is_current: false
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/groups');
    });

    it('returns empty array when no groups', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await groupService.list();

      expect(result).toEqual([]);
    });

    it('throws when api fails', async () => {
      mockedApi.get.mockRejectedValueOnce(new Error('Unauthorized'));

      await expect(groupService.list()).rejects.toThrow('Unauthorized');
    });
  });

  describe('create', () => {
    it('creates a group and returns mapped response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'g3', attributes: { title: 'Family' } }
      });

      const result = await groupService.create({ title: 'Family' });

      expect(result).toEqual({
        id: 'g3',
        title: 'Family',
        member_count: 1,
        is_current: false
      });
      expect(mockedApi.post).toHaveBeenCalledWith('/groups', { title: 'Family' });
    });

    it('falls back to input title when response has none', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'g4', attributes: {} }
      });

      const result = await groupService.create({ title: 'Test' });

      expect(result.title).toBe('Test');
    });
  });

  describe('switch', () => {
    it('calls switch endpoint with group id', async () => {
      mockedApi.post.mockResolvedValueOnce({});

      await groupService.switch('g1');

      expect(mockedApi.post).toHaveBeenCalledWith('/groups/switch', { user_group_id: 'g1' });
    });

    it('throws when api fails', async () => {
      mockedApi.post.mockRejectedValueOnce(new Error('Switch failed'));

      await expect(groupService.switch('g1')).rejects.toThrow('Switch failed');
    });
  });
});
