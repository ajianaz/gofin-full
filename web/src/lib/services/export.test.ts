import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(() => null),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
};

describe('exportService', () => {
  let exportService: typeof import('./export.js').exportService;

  beforeEach(async () => {
    vi.clearAllMocks();
    Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock, writable: true });

    // Mock only URL.createObjectURL and revokeObjectURL, keep URL as a constructor
    const originalCreateObjectURL = vi.fn(() => 'blob:http://localhost/fake-uuid');
    const originalRevokeObjectURL = vi.fn();
    vi.spyOn(URL, 'createObjectURL').mockImplementation(originalCreateObjectURL);
    vi.spyOn(URL, 'revokeObjectURL').mockImplementation(originalRevokeObjectURL);

    // Mock document.createElement
    const clickFn = vi.fn();
    const mockAnchor = { href: '', download: '', click: clickFn };
    vi.spyOn(document, 'createElement').mockReturnValue(mockAnchor as any);

    // Re-import to get fresh module with mocks
    vi.resetModules();
    const mod = await import('./export.js');
    exportService = mod.exportService;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('downloadCSV', () => {
    it('downloads CSV blob with correct URL', async () => {
      const mockBlob = new Blob(['csv,data'], { type: 'text/csv' });
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: true,
        blob: () => Promise.resolve(mockBlob)
      }));

      await exportService.downloadCSV('2026-01-01', '2026-01-31', 'w1');

      expect(fetch).toHaveBeenCalledWith(
        '/api/v1/export/csv?start=2026-01-01&end=2026-01-31&wallet_id=w1',
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json' }
        })
      );
      expect(URL.createObjectURL).toHaveBeenCalledWith(mockBlob);
    });

    it('includes auth header when token exists', async () => {
      localStorageMock.getItem.mockReturnValue('my-token' as any);
      const mockBlob = new Blob(['csv'], { type: 'text/csv' });
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: true,
        blob: () => Promise.resolve(mockBlob)
      }));

      await exportService.downloadCSV();

      expect(fetch).toHaveBeenCalledWith(
        '/api/v1/export/csv',
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json', Authorization: 'Bearer my-token' }
        })
      );
    });

    it('builds URL without params when none provided', async () => {
      const mockBlob = new Blob(['csv']);
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: true,
        blob: () => Promise.resolve(mockBlob)
      }));

      await exportService.downloadCSV();

      expect(fetch).toHaveBeenCalledWith('/api/v1/export/csv', expect.any(Object));
    });

    it('throws when response is not ok', async () => {
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: false,
        statusText: 'Not Found'
      }));

      await expect(exportService.downloadCSV()).rejects.toThrow('Export failed: Not Found');
    });
  });

  describe('downloadOFX', () => {
    it('downloads OFX blob with correct URL', async () => {
      const mockBlob = new Blob(['ofx,data'], { type: 'application/x-ofx' });
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: true,
        blob: () => Promise.resolve(mockBlob)
      }));

      await exportService.downloadOFX('2026-01-01', '2026-01-31', 'w1');

      expect(fetch).toHaveBeenCalledWith(
        '/api/v1/export/ofx?start=2026-01-01&end=2026-01-31&wallet_id=w1',
        expect.any(Object)
      );
    });

    it('throws when response is not ok', async () => {
      vi.stubGlobal('fetch', vi.fn().mockResolvedValue({
        ok: false,
        statusText: 'Internal Server Error'
      }));

      await expect(exportService.downloadOFX()).rejects.toThrow('Export failed: Internal Server Error');
    });
  });
});
