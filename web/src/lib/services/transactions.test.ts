import { describe, it, expect, vi, beforeEach } from 'vitest';
import { transactionService } from './transactions.js';

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

describe('transactionService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps transaction fields with pagination', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 't1',
            attributes: {
              amount: '100.50',
              date: '2026-01-15',
              description: 'Grocery shopping',
              type: 'withdrawal'
            }
          }
        ],
        meta: {
          pagination: {
            total: 50,
            count: 1,
            per_page: 20,
            current_page: 1,
            total_pages: 3
          }
        }
      });

      const result = await transactionService.list({ page: 1, per_page: 20 });

      expect(result.data).toHaveLength(1);
      expect(result.data[0]).toEqual(
        expect.objectContaining({
          id: 't1',
          amount: '100.50',
          date: '2026-01-15',
          description: 'Grocery shopping',
          tags: []
        })
      );
      expect(result.meta).toEqual({
        pagination: {
          total: 50,
          count: 1,
          per_page: 20,
          current_page: 1,
          total_pages: 3
        }
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/transactions?page=1&per_page=20');
    });

    it('applies defaults for missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 't2', attributes: { created_at: '2026-02-01', group_title: 'Transfer' } }
        ]
      });

      const result = await transactionService.list();

      expect(result.data[0]).toEqual(
        expect.objectContaining({
          amount: '0',
          date: '2026-02-01',
          description: 'Transfer',
          tags: []
        })
      );
    });

    it('builds query string with all params', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      await transactionService.list({
        page: 2,
        per_page: 10,
        start: '2026-01-01',
        end: '2026-01-31',
        type: 'withdrawal'
      });

      expect(mockedApi.get).toHaveBeenCalledWith(
        '/transactions?page=2&per_page=10&start=2026-01-01&end=2026-01-31&type=withdrawal'
      );
    });

    it('builds no query string when no params', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      await transactionService.list();

      expect(mockedApi.get).toHaveBeenCalledWith('/transactions');
    });
  });

  describe('create', () => {
    it('creates transaction and serializes date as ISO', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: {
          id: 't3',
          attributes: { type: 'withdrawal', description: 'Test', amount: '50' }
        }
      });

      const result = await transactionService.create({
        type: 'withdrawal',
        description: 'Test',
        amount: '50',
        source_id: 'w1',
        destination_id: 'w2',
        date: '2026-03-01'
      });

      expect(result).toEqual(
        expect.objectContaining({
          id: 't3',
          type: 'withdrawal',
          description: 'Test',
          amount: '50'
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/transactions', expect.objectContaining({
        date: new Date('2026-03-01').toISOString()
      }));
    });

    it('creates without date param', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 't4', attributes: {} }
      });

      await transactionService.create({
        type: 'deposit',
        amount: '100',
        source_id: 'w1',
        destination_id: 'w2'
      });

      const callArgs = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callArgs).not.toHaveProperty('date');
    });
  });

  describe('update', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await transactionService.update('t1', { description: 'Updated' });

      expect(mockedApi.put).toHaveBeenCalledWith('/transactions/t1', { description: 'Updated' });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await transactionService.delete('t1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/transactions/t1');
    });
  });
});
