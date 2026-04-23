import { describe, it, expect, vi, beforeEach } from 'vitest';
import { categoryService } from './categories.js';

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

describe('categoryService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps category fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'c1', attributes: { name: 'Food' } },
          { id: 'c2', attributes: { name: 'Transport' } }
        ]
      });

      const result = await categoryService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'c1',
          name: 'Food',
          type: 'expense',
          transaction_count: 0
        })
      );
      expect(result[1]).toEqual(
        expect.objectContaining({
          id: 'c2',
          name: 'Transport',
          type: 'expense',
          transaction_count: 0
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/categories');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await categoryService.list();

      expect(result).toEqual([]);
    });

    it('throws when api fails', async () => {
      mockedApi.get.mockRejectedValueOnce(new Error('Network error'));

      await expect(categoryService.list()).rejects.toThrow('Network error');
    });
  });

  describe('create', () => {
    it('creates category and maps response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'c3', attributes: { name: 'Entertainment' } }
      });

      const result = await categoryService.create({ name: 'Entertainment' });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'c3',
          name: 'Entertainment',
          type: 'expense',
          transaction_count: 0
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/categories', { name: 'Entertainment' });
    });
  });

  describe('update', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await categoryService.update('c1', { name: 'Groceries' });

      expect(mockedApi.put).toHaveBeenCalledWith('/categories/c1', { name: 'Groceries' });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await categoryService.delete('c1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/categories/c1');
    });
  });
});
