import { describe, it, expect, vi, beforeEach } from 'vitest';
import { budgetService } from './budgets.js';

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

describe('budgetService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps budget fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'b1', attributes: { name: 'Monthly Budget' } }
        ]
      });

      const result = await budgetService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'b1',
          name: 'Monthly Budget',
          spend_amount: '0',
          budget_amount: '0',
          limits: []
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/budgets');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await budgetService.list();

      expect(result).toEqual([]);
    });
  });

  describe('create', () => {
    it('creates budget and maps response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'b2', attributes: { name: 'Savings Budget', order: 1 } }
      });

      const result = await budgetService.create({ name: 'Savings Budget', order: 1 });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'b2',
          name: 'Savings Budget',
          spend_amount: '0',
          budget_amount: '0',
          limits: []
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/budgets', { name: 'Savings Budget', order: 1 });
    });
  });

  describe('update', () => {
    it('calls put with correct path and data', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await budgetService.update('b1', { name: 'Updated Budget', active: false });

      expect(mockedApi.put).toHaveBeenCalledWith('/budgets/b1', { name: 'Updated Budget', active: false });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await budgetService.delete('b1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/budgets/b1');
    });
  });
});
