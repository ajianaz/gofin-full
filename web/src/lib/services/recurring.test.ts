import { describe, it, expect, vi, beforeEach } from 'vitest';
import { recurringService } from './recurring.js';

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

describe('recurringService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps recurring transaction fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 'r1',
            attributes: {
              title: 'Monthly Rent',
              repeat_freq: 'monthly'
            }
          }
        ]
      });

      const result = await recurringService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'r1',
          title: 'Monthly Rent',
          type: 'withdrawal',
          amount: '0',
          currency_code: 'USD'
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/recurrences');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await recurringService.list();

      expect(result).toEqual([]);
    });
  });

  describe('create', () => {
    it('creates recurring transaction and serializes dates', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: {
          id: 'r2',
          attributes: {
            title: 'Weekly Savings',
            repeat_freq: 'weekly'
          }
        }
      });

      const result = await recurringService.create({
        title: 'Weekly Savings',
        first_date: '2026-02-01',
        repeat_freq: 'weekly',
        repeat_until: '2026-12-31',
        transactions: [{
          type: 'withdrawal',
          amount: '100',
          source_id: 'w1',
          destination_id: 'w2'
        }]
      });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'r2',
          title: 'Weekly Savings',
          type: 'withdrawal',
          amount: '0',
          currency_code: 'USD'
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/recurrences', expect.objectContaining({
        first_date: new Date('2026-02-01').toISOString(),
        repeat_until: new Date('2026-12-31').toISOString()
      }));
    });

    it('does not serialize undefined date fields', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'r3', attributes: { title: 'Test' } }
      });

      await recurringService.create({
        title: 'Test',
        first_date: '2026-01-01'
      });

      const callBody = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callBody.first_date).toBe(new Date('2026-01-01').toISOString());
      expect(callBody.repeat_until).toBeUndefined();
    });
  });

  describe('update', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await recurringService.update('r1', { title: 'Updated', repeat_freq: 'daily', active: false });

      expect(mockedApi.put).toHaveBeenCalledWith('/recurrences/r1', {
        title: 'Updated',
        repeat_freq: 'daily',
        active: false
      });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await recurringService.delete('r1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/recurrences/r1');
    });
  });
});
