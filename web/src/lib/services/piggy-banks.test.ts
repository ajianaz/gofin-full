import { describe, it, expect, vi, beforeEach } from 'vitest';
import { piggyBankService } from './piggy-banks.js';

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

describe('piggyBankService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps piggy bank fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          {
            id: 'pb1',
            attributes: {
              name: 'Vacation',
              wallet_id: 'w1',
              target_amount: '5000',
              current_amount: '2000'
            }
          }
        ]
      });

      const result = await piggyBankService.list('w1');

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'pb1',
          name: 'Vacation',
          account_id: 'w1',
          account_name: '',
          current_amount: '0',
          status: 'active'
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/wallets/w1/piggy_banks');
    });

    it('uses walletId when wallet_id attribute is missing', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'pb2', attributes: { name: 'Emergency' } }
        ]
      });

      const result = await piggyBankService.list('w2');

      expect(result[0].account_id).toBe('w2');
    });
  });

  describe('create', () => {
    it('creates piggy bank and maps response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: {
          id: 'pb3',
          attributes: { name: 'New Car', wallet_id: 'w1', target_amount: '20000' }
        }
      });

      const result = await piggyBankService.create({
        wallet_id: 'w1',
        name: 'New Car',
        target_amount: '20000'
      });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'pb3',
          name: 'New Car',
          account_id: 'w1',
          account_name: '',
          current_amount: '0',
          status: 'active'
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith(
        '/wallets/w1/piggy_banks',
        expect.objectContaining({
          wallet_id: 'w1',
          name: 'New Car',
          target_amount: '20000'
        })
      );
    });

    it('stringifies target_amount', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'pb4', attributes: { name: 'Test' } }
      });

      await piggyBankService.create({
        wallet_id: 'w1',
        name: 'Test',
        target_amount: '5000'
      });

      const callBody = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callBody.target_amount).toBe('5000');
    });

    it('omits target_amount when not provided', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'pb5', attributes: { name: 'Test' } }
      });

      await piggyBankService.create({
        wallet_id: 'w1',
        name: 'Test'
      });

      const callBody = vi.mocked(mockedApi.post).mock.calls[0][1] as Record<string, unknown>;
      expect(callBody.target_amount).toBeUndefined();
    });
  });

  describe('update', () => {
    it('calls put with correct path and stringifies target_amount', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await piggyBankService.update('w1', 'pb1', { name: 'Updated', target_amount: '7000' });

      expect(mockedApi.put).toHaveBeenCalledWith('/wallets/w1/piggy_banks/pb1', {
        name: 'Updated',
        target_amount: '7000'
      });
    });

    it('omits undefined fields from payload', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await piggyBankService.update('w1', 'pb1', {});

      expect(mockedApi.put).toHaveBeenCalledWith('/wallets/w1/piggy_banks/pb1', {});
    });
  });

  describe('delete', () => {
    it('calls delete with correct path including walletId', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await piggyBankService.delete('w1', 'pb1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/wallets/w1/piggy_banks/pb1');
    });
  });
});
