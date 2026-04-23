import { describe, it, expect, vi, beforeEach } from 'vitest';
import { walletMemberService } from './wallet-members.js';

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

describe('walletMemberService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('maps JSON:API response to WalletMember array', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'm1', attributes: { wallet_id: 'w1', user_id: 'u1', role: 'owner' } },
          { id: 'm2', attributes: { wallet_id: 'w1', user_id: 'u2', role: 'editor' } }
        ]
      });

      const result = await walletMemberService.list('w1');

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({ id: 'm1', wallet_id: 'w1', user_id: 'u1', role: 'owner' });
      expect(result[1].role).toBe('editor');
      expect(mockedApi.get).toHaveBeenCalledWith('/wallets/w1/members');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await walletMemberService.list('w1');

      expect(result).toEqual([]);
    });
  });

  describe('add', () => {
    it('posts member with user_id and role', async () => {
      mockedApi.post.mockResolvedValueOnce({});

      await walletMemberService.add('w1', 'u3', 'viewer');

      expect(mockedApi.post).toHaveBeenCalledWith('/wallets/w1/members', { user_id: 'u3', role: 'viewer' });
    });
  });

  describe('updateRole', () => {
    it('puts updated role', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await walletMemberService.updateRole('w1', 'm2', 'editor');

      expect(mockedApi.put).toHaveBeenCalledWith('/wallets/w1/members/m2', { role: 'editor' });
    });
  });

  describe('remove', () => {
    it('deletes member by user_id', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await walletMemberService.remove('w1', 'u2');

      expect(mockedApi.delete).toHaveBeenCalledWith('/wallets/w1/members/u2');
    });
  });
});
