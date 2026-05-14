import { describe, it, expect, vi, beforeEach } from 'vitest';
import { tagService } from './tags.js';

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

describe('tagService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps tag fields with default date', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'tg1', attributes: { tag: 'groceries', date: '2026-01-15' } }
        ]
      });

      const result = await tagService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'tg1',
          tag: 'groceries',
          date: '2026-01-15'
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/tags');
    });

    it('applies today date when date is missing', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'tg2', attributes: { tag: 'urgent' } }
        ]
      });

      const result = await tagService.list();

      const today = new Date().toISOString().split('T')[0];
      expect(result[0].date).toBe(today);
    });
  });

  describe('create', () => {
    it('creates tag and serializes date as ISO', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'tg3', attributes: { tag: 'travel', date: '2026-03-01' } }
      });

      const result = await tagService.create({ tag: 'travel', date: '2026-03-01' });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'tg3',
          tag: 'travel',
          date: '2026-03-01'
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/tags', {
        tag: 'travel',
        date: new Date('2026-03-01').toISOString()
      });
    });

    it('creates tag without date', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'tg4', attributes: { tag: 'personal' } }
      });

      await tagService.create({ tag: 'personal' });

      const callBody = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callBody).not.toHaveProperty('date');
    });

    it('applies default date when response has no date', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'tg5', attributes: { tag: 'test' } }
      });

      const result = await tagService.create({ tag: 'test' });

      const today = new Date().toISOString().split('T')[0];
      expect(result.date).toBe(today);
    });
  });

  describe('update', () => {
    it('calls put with serialized date', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await tagService.update('tg1', { tag: 'food', date: '2026-04-01' });

      expect(mockedApi.put).toHaveBeenCalledWith('/tags/tg1', {
        tag: 'food',
        date: new Date('2026-04-01').toISOString()
      });
    });

    it('calls put without date when not provided', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await tagService.update('tg1', { tag: 'updated' });

      expect(mockedApi.put).toHaveBeenCalledWith('/tags/tg1', { tag: 'updated' });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await tagService.delete('tg1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/tags/tg1');
    });
  });
});
