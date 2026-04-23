import { describe, it, expect, vi, beforeEach } from 'vitest';
import { walletService } from './wallets.js';

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

describe('walletService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps wallet fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 'w1',
            attributes: {
              name: 'Cash',
              wallet_type: 'cash',
              virtual_balance: '5000',
              currency_id: 'USD'
            }
          }
        ]
      });

      const result = await walletService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'w1',
          name: 'Cash',
          type: 'cash',
          balance: '5000',
          currency_code: 'USD',
          currency_symbol: 'Rp',
          currency_decimal_places: 0
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/wallets');
    });

    it('applies defaults when wallet_type and virtual_balance are missing', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'w2', attributes: { name: 'Bank' } }
        ]
      });

      const result = await walletService.list();

      expect(result[0]).toEqual(
        expect.objectContaining({
          type: 'asset',
          balance: '0',
          currency_code: 'USD'
        })
      );
    });

    it('returns empty array when no wallets', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await walletService.list();

      expect(result).toEqual([]);
    });

    it('throws when api fails', async () => {
      mockedApi.get.mockRejectedValueOnce(new Error('Network error'));

      await expect(walletService.list()).rejects.toThrow('Network error');
    });
  });

  describe('create', () => {
    it('creates wallet and maps response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: {
          id: 'w3',
          attributes: { name: 'Savings', wallet_type: 'savings' }
        }
      });

      const result = await walletService.create({ name: 'Savings', wallet_type: 'savings' });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'w3',
          name: 'Savings',
          type: 'savings',
          balance: '0',
          currency_code: 'USD',
          currency_symbol: 'Rp',
          currency_decimal_places: 0
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/wallets', {
        name: 'Savings',
        wallet_type: 'savings'
      });
    });
  });

  describe('update', () => {
    it('calls put with correct path and data', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await walletService.update('w1', { name: 'Updated', active: true });

      expect(mockedApi.put).toHaveBeenCalledWith('/wallets/w1', { name: 'Updated', active: true });
    });

    it('throws when api fails', async () => {
      mockedApi.put.mockRejectedValueOnce(new Error('Update failed'));

      await expect(walletService.update('w1', { name: 'X' })).rejects.toThrow('Update failed');
    });
  });

  describe('delete', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await walletService.delete('w1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/wallets/w1');
    });

    it('throws when api fails', async () => {
      mockedApi.delete.mockRejectedValueOnce(new Error('Delete failed'));

      await expect(walletService.delete('w1')).rejects.toThrow('Delete failed');
    });
  });

  describe('types', () => {
    it('returns wallet types from raw data', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'cash', attributes: { type: 'Cash wallet' } },
          { id: 'bank', attributes: { type: 'Bank account' } }
        ]
      });

      const result = await walletService.types();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({ id: 'cash', attributes: { type: 'Cash wallet' } });
      expect(result[1]).toEqual({ id: 'bank', attributes: { type: 'Bank account' } });
      expect(mockedApi.get).toHaveBeenCalledWith('/wallet-types');
    });
  });
});
