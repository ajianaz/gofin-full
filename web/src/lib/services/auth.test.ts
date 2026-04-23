import { describe, it, expect, vi, beforeEach } from 'vitest';
import { authService } from './auth.js';

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

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('login', () => {
    it('returns token response from api', async () => {
      const tokenResponse = {
        access_token: 'abc123',
        refresh_token: 'xyz789',
        token_type: 'Bearer'
      };
      mockedApi.post.mockResolvedValueOnce(tokenResponse);

      const result = await authService.login({
        email: 'test@test.com',
        password: 'password123'
      });

      expect(result).toEqual(tokenResponse);
      expect(mockedApi.post).toHaveBeenCalledWith('/auth/login', {
        email: 'test@test.com',
        password: 'password123'
      });
    });
  });

  describe('register', () => {
    it('returns token response from api', async () => {
      const tokenResponse = {
        access_token: 'new-token',
        refresh_token: 'new-refresh',
        token_type: 'Bearer'
      };
      mockedApi.post.mockResolvedValueOnce(tokenResponse);

      const result = await authService.register({
        email: 'new@test.com',
        password: 'pass123',
        name: 'New User'
      });

      expect(result).toEqual(tokenResponse);
      expect(mockedApi.post).toHaveBeenCalledWith('/auth/register', {
        email: 'new@test.com',
        password: 'pass123',
        name: 'New User'
      });
    });
  });

  describe('logout', () => {
    it('sends refresh token to logout endpoint', async () => {
      mockedApi.post.mockResolvedValueOnce({ message: 'Logged out' });

      const result = await authService.logout('refresh-123');

      expect(result).toEqual({ message: 'Logged out' });
      expect(mockedApi.post).toHaveBeenCalledWith('/auth/logout', { refresh_token: 'refresh-123' });
    });

    it('works without refresh token', async () => {
      mockedApi.post.mockResolvedValueOnce({ message: 'Logged out' });

      await authService.logout();

      expect(mockedApi.post).toHaveBeenCalledWith('/auth/logout', { refresh_token: undefined });
    });
  });

  describe('refresh', () => {
    it('sends refresh token and returns new tokens', async () => {
      const newTokens = {
        access_token: 'new-access',
        refresh_token: 'new-refresh',
        token_type: 'Bearer'
      };
      mockedApi.post.mockResolvedValueOnce(newTokens);

      const result = await authService.refresh('old-refresh');

      expect(result).toEqual(newTokens);
      expect(mockedApi.post).toHaveBeenCalledWith('/auth/refresh', { refresh_token: 'old-refresh' });
    });
  });

  describe('getMe', () => {
    it('unwraps user data and derives name from email', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: {
          id: 'u1',
          attributes: {
            email: 'john.doe@test.com',
            created_at: '2026-01-15T10:00:00Z'
          }
        }
      });

      const result = await authService.getMe();

      expect(result).toEqual({
        id: 'u1',
        email: 'john.doe@test.com',
        name: 'john.doe',
        created_at: '2026-01-15T10:00:00Z'
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/users/me');
    });

    it('defaults name to empty string when email is empty', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: {
          id: 'u2',
          attributes: { email: '' }
        }
      });

      const result = await authService.getMe();

      expect(result.name).toBe('');
      expect(result.email).toBe('');
    });

    it('defaults created_at to current ISO date when missing', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: {
          id: 'u3',
          attributes: { email: 'test@test.com' }
        }
      });

      const result = await authService.getMe();

      expect(result.created_at).toBeTruthy();
      expect(typeof result.created_at).toBe('string');
    });
  });

  describe('getProvider', () => {
    it('returns auth provider info', async () => {
      mockedApi.get.mockResolvedValueOnce({ provider: 'keycloak' });

      const result = await authService.getProvider();

      expect(result).toEqual({ provider: 'keycloak' });
      expect(mockedApi.get).toHaveBeenCalledWith('/auth/provider');
    });
  });
});
