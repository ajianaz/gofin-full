import { describe, it, expect, vi, beforeEach } from 'vitest';
import { reportService } from './reports.js';

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

describe('reportService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('spendingByCategory', () => {
    it('unwraps and maps category spending fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: '1', attributes: { category_id: 'cat1', category_name: 'Food', total: '150.50', count: '12' } },
          { id: '2', attributes: { category_id: 'cat2', category_name: 'Transport', total: '80', count: '5' } }
        ]
      });

      const result = await reportService.spendingByCategory('2026-01-01', '2026-01-31');

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({
        category_id: 'cat1',
        category_name: 'Food',
        total: 150.50,
        count: 12
      });
      expect(result[1]).toEqual({
        category_id: 'cat2',
        category_name: 'Transport',
        total: 80,
        count: 5
      });
      expect(mockedApi.get).toHaveBeenCalledWith(
        '/analytics/spending-by-category?start=2026-01-01&end=2026-01-31'
      );
    });

    it('builds URL without query params when no dates', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      await reportService.spendingByCategory();

      expect(mockedApi.get).toHaveBeenCalledWith('/analytics/spending-by-category');
    });

    it('applies defaults for missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: '3', attributes: {} }]
      });

      const result = await reportService.spendingByCategory();

      expect(result[0]).toEqual({
        category_id: '',
        category_name: '',
        total: 0,
        count: 0
      });
    });
  });

  describe('spendingByPeriod', () => {
    it('unwraps and maps period spending fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: '1', attributes: { period: '2026-01', income: '5000', expense: '3200' } }
        ]
      });

      const result = await reportService.spendingByPeriod('2026-01-01', '2026-01-31');

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual({
        period: '2026-01',
        income: 5000,
        expense: 3200
      });
      expect(mockedApi.get).toHaveBeenCalledWith(
        '/analytics/spending-by-period?start=2026-01-01&end=2026-01-31'
      );
    });

    it('applies defaults for missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: '2', attributes: { period: '2026-02' } }]
      });

      const result = await reportService.spendingByPeriod();

      expect(result[0]).toEqual({
        period: '2026-02',
        income: 0,
        expense: 0
      });
    });
  });

  describe('netWorth', () => {
    it('unwraps single data response and parses fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: {
          type: 'analytics',
          attributes: {
            total_income: '10000',
            total_expense: '7500',
            net_income: '2500',
            transaction_count: '45'
          }
        }
      });

      const result = await reportService.netWorth('2026-01-01', '2026-12-31');

      expect(result).toEqual({
        total_income: 10000,
        total_expense: 7500,
        net_income: 2500,
        transaction_count: 45
      });
      expect(mockedApi.get).toHaveBeenCalledWith(
        '/analytics/net-worth?start=2026-01-01&end=2026-12-31'
      );
    });

    it('applies defaults when no dates provided', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: { type: 'analytics', attributes: {} }
      });

      const result = await reportService.netWorth();

      expect(result).toEqual({
        total_income: 0,
        total_expense: 0,
        net_income: 0,
        transaction_count: 0
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/analytics/net-worth');
    });
  });
});
