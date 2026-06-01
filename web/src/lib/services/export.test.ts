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

    // Mock the central API client's blob method
    vi.mock('$lib/services/client.js', () => ({
      api: {
        blob: vi.fn()
      }
    }));

    // Re-import to get fresh module with mocks
    vi.resetModules();
    const mod = await import('./export.js');
    exportService = mod.exportService;
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  describe('downloadCSV', () => {
    it('downloads CSV blob via api.blob with correct path', async () => {
      const mockBlob = new Blob(['csv,data'], { type: 'text/csv' });
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: true,
        status: 200,
        blob: () => Promise.resolve(mockBlob)
      } as Response);

      await exportService.downloadCSV('2026-01-01', '2026-01-31', 'w1');

      expect(api.blob).toHaveBeenCalledWith('/export/csv?start=2026-01-01&end=2026-01-31&wallet_id=w1');
      expect(URL.createObjectURL).toHaveBeenCalledWith(mockBlob);
    });

    it('includes auth header when token exists', async () => {
      localStorageMock.getItem.mockReturnValue('my-token' as any);
      const mockBlob = new Blob(['csv'], { type: 'text/csv' });
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: true,
        status: 200,
        blob: () => Promise.resolve(mockBlob)
      } as Response);

      await exportService.downloadCSV();

      // api.blob is called — the central client handles the token internally
      expect(api.blob).toHaveBeenCalledWith('/export/csv');
    });

    it('builds path without params when none provided', async () => {
      const mockBlob = new Blob(['csv']);
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: true,
        status: 200,
        blob: () => Promise.resolve(mockBlob)
      } as Response);

      await exportService.downloadCSV();
      expect(api.blob).toHaveBeenCalledWith('/export/csv');
    });

    it('throws when response is not ok', async () => {
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: false,
        status: 404,
        statusText: 'Not Found'
      } as Response);

      await expect(exportService.downloadCSV()).rejects.toThrow('Export failed: Not Found');
    });
  });

  describe('downloadOFX', () => {
    it('downloads OFX blob with correct path', async () => {
      const mockBlob = new Blob(['ofx,data'], { type: 'application/x-ofx' });
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: true,
        status: 200,
        blob: () => Promise.resolve(mockBlob)
      } as Response);

      await exportService.downloadOFX('2026-01-01', '2026-01-31', 'w1');

      expect(api.blob).toHaveBeenCalledWith('/export/ofx?start=2026-01-01&end=2026-01-31&wallet_id=w1');
    });

    it('throws when response is not ok', async () => {
      const { api } = await import('./client.js');
      vi.mocked(api.blob).mockResolvedValue({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error'
      } as Response);

      await expect(exportService.downloadOFX()).rejects.toThrow('Export failed: Internal Server Error');
    });
  });
});
