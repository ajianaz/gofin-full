import { describe, it, expect, vi, beforeEach } from 'vitest';
import { notificationService } from './notifications.js';

vi.mock('./client.js', () => ({
  api: {
    get: vi.fn(),
    put: vi.fn()
  }
}));

import { api } from './client.js';
const mockedApi = vi.mocked(api);

describe('notificationService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('maps JSON:API response to Notification array', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'n1', attributes: { channel: 'email', type: 'info', title: 'New Transaction', message: 'You received...', read: false } },
          { id: 'n2', attributes: { channel: 'push', type: 'warning', title: 'Budget Exceeded', message: 'Over budget...', read: true } }
        ]
      });

      const result = await notificationService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({ id: 'n1', channel: 'email', type: 'info', title: 'New Transaction', message: 'You received...', read: false });
      expect(result[1].read).toBe(true);
      expect(mockedApi.get).toHaveBeenCalledWith('/notifications');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await notificationService.list();

      expect(result).toEqual([]);
    });

    it('defaults missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'n1', attributes: { title: 'Test' } }
        ]
      });

      const result = await notificationService.list();

      expect(result[0]).toEqual({ id: 'n1', channel: '', type: '', title: 'Test', message: '', read: false });
    });
  });

  describe('listUnread', () => {
    it('calls unread endpoint', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      await notificationService.listUnread();

      expect(mockedApi.get).toHaveBeenCalledWith('/notifications/unread');
    });
  });

  describe('markRead', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await notificationService.markRead('n1');

      expect(mockedApi.put).toHaveBeenCalledWith('/notifications/n1/read');
    });
  });

  describe('markAllRead', () => {
    it('calls put on read-all endpoint', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await notificationService.markAllRead();

      expect(mockedApi.put).toHaveBeenCalledWith('/notifications/read-all');
    });
  });
});
