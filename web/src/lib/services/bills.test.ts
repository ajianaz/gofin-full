import { describe, it, expect, vi, beforeEach } from 'vitest';
import { billService } from './bills.js';

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

describe('billService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps bill fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 'bl1',
            attributes: {
              name: 'Electricity',
              date: '2026-02-01',
              currency_id: 'USD',
              amount_min: '50',
              amount_max: '100'
            }
          }
        ]
      });

      const result = await billService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'bl1',
          name: 'Electricity',
          next_date: '2026-02-01',
          currency_code: 'USD',
          currency_symbol: 'Rp'
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/bills');
    });

    it('applies defaults when date and currency_id are missing', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: 'bl2', attributes: { name: 'Internet' } }]
      });

      const result = await billService.list();

      expect(result[0]).toEqual(
        expect.objectContaining({
          next_date: '',
          currency_code: 'USD'
        })
      );
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await billService.list();

      expect(result).toEqual([]);
    });
  });

  describe('create', () => {
    it('creates bill and stringifies numeric fields', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: {
          id: 'bl3',
          attributes: { name: 'Rent', date: '2026-03-01' }
        }
      });

      const result = await billService.create({
        name: 'Rent',
        amount_min: '1500',
        amount_max: '1500',
        date: '2026-03-01'
      });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'bl3',
          next_date: '2026-03-01',
          currency_code: 'USD',
          currency_symbol: 'Rp'
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/bills', expect.objectContaining({
        name: 'Rent',
        amount_min: '1500',
        amount_max: '1500',
        date: '2026-03-01'
      }));
    });

    it('does not stringify undefined amount fields', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'bl4', attributes: { name: 'Phone' } }
      });

      await billService.create({ name: 'Phone' });

      const callBody = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callBody.amount_min).toBeUndefined();
      expect(callBody.amount_max).toBeUndefined();
    });
  });

  describe('update', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await billService.update('bl1', { name: 'Updated Bill', active: true });

      expect(mockedApi.put).toHaveBeenCalledWith('/bills/bl1', { name: 'Updated Bill', active: true });
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await billService.delete('bl1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/bills/bl1');
    });
  });
});
