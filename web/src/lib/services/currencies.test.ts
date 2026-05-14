import { describe, it, expect, vi, beforeEach } from 'vitest';
import { currencyService } from './currencies.js';

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

describe('currencyService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps JSON:API response and maps fields correctly', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'USD', attributes: { name: 'US Dollar', symbol: '$', decimal_places: 2, enabled: true } },
          { id: 'EUR', attributes: { name: 'Euro', symbol: '\u20AC', decimal_places: 2, enabled: false } }
        ]
      });

      const result = await currencyService.list();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual({
        id: 'USD',
        name: 'US Dollar',
        symbol: '$',
        code: 'USD',
        decimal_places: 2,
        enabled: true
      });
      expect(result[1]).toEqual({
        id: 'EUR',
        name: 'Euro',
        symbol: '\u20AC',
        code: 'EUR',
        decimal_places: 2,
        enabled: false
      });
      expect(mockedApi.get).toHaveBeenCalledWith('/currencies');
    });

    it('applies defaults for missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'JPY', attributes: { name: 'Japanese Yen', symbol: '\u00A5' } }
        ]
      });

      const result = await currencyService.list();

      expect(result[0]).toEqual(
        expect.objectContaining({
          code: 'JPY',
          decimal_places: 2,
          enabled: true
        })
      );
    });

    it('returns empty array when no currencies', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await currencyService.list();

      expect(result).toEqual([]);
    });

    it('throws when api.get fails', async () => {
      mockedApi.get.mockRejectedValueOnce(new Error('Network error'));

      await expect(currencyService.list()).rejects.toThrow('Network error');
    });
  });

  describe('exchangeRates', () => {
    it('unwraps and maps exchange rate fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: '1', attributes: { from_currency_id: 'USD', to_currency_id: 'EUR', rate: '0.85', date: '2026-01-01' } },
          { id: '2', attributes: { from_code: 'GBP', to_code: 'USD', rate: '1.25', date: '2026-01-02' } }
        ]
      });

      const result = await currencyService.exchangeRates();

      expect(result).toHaveLength(2);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: '1',
          from_code: 'USD',
          to_code: 'EUR',
          rate: 0.85,
          date: '2026-01-01'
        })
      );
      expect(result[1]).toEqual(
        expect.objectContaining({
          id: '2',
          from_code: 'GBP',
          to_code: 'USD',
          rate: 1.25,
          date: '2026-01-02'
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/exchange-rates');
    });

    it('applies defaults for missing fields', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: '3', attributes: {} }]
      });

      const result = await currencyService.exchangeRates();

      expect(result[0]).toEqual(
        expect.objectContaining({
          from_code: '',
          to_code: '',
          rate: 0,
          date: ''
        })
      );
    });

    it('parses numeric rate from string', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [{ id: '4', attributes: { rate: '1.123456' } }]
      });

      const result = await currencyService.exchangeRates();

      expect(result[0].rate).toBeCloseTo(1.123456);
    });
  });
});
