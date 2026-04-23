import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ruleService } from './rules.js';

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

describe('ruleService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('unwraps and maps rule group fields with defaults', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: [
          { id: 'rg1', attributes: { title: 'Auto-categorize', active: true } }
        ]
      });

      const result = await ruleService.list();

      expect(result).toHaveLength(1);
      expect(result[0]).toEqual(
        expect.objectContaining({
          id: 'rg1',
          title: 'Auto-categorize',
          stop_processing: false,
          rule_count: 0
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/rules');
    });

    it('returns empty array', async () => {
      mockedApi.get.mockResolvedValueOnce({ data: [] });

      const result = await ruleService.list();

      expect(result).toEqual([]);
    });
  });

  describe('createGroup', () => {
    it('creates rule group and maps response', async () => {
      mockedApi.post.mockResolvedValueOnce({
        data: { id: 'rg2', attributes: { title: 'New Group' } }
      });

      const result = await ruleService.createGroup({ title: 'New Group' });

      expect(result).toEqual(
        expect.objectContaining({
          id: 'rg2',
          title: 'New Group',
          stop_processing: false,
          rule_count: 0
        })
      );
      expect(mockedApi.post).toHaveBeenCalledWith('/rule-groups', { title: 'New Group' });
    });
  });

  describe('get', () => {
    it('gets single rule by id', async () => {
      mockedApi.get.mockResolvedValueOnce({
        data: {
          id: 'rule1',
          attributes: {
            title: 'Categorize Groceries',
            rule_group_id: 'rg1',
            active: true
          }
        }
      });

      const result = await ruleService.get('rule1');

      expect(result).toEqual(
        expect.objectContaining({
          id: 'rule1',
          title: 'Categorize Groceries',
          rule_group_id: 'rg1',
          active: true
        })
      );
      expect(mockedApi.get).toHaveBeenCalledWith('/rules/rule1');
    });
  });

  describe('updateGroup', () => {
    it('calls put with correct path', async () => {
      mockedApi.put.mockResolvedValueOnce({});

      await ruleService.updateGroup('rg1', { title: 'Updated', active: false });

      expect(mockedApi.put).toHaveBeenCalledWith('/rule-groups/rg1', { title: 'Updated', active: false });
    });
  });

  describe('deleteGroup', () => {
    it('calls delete with correct path', async () => {
      mockedApi.delete.mockResolvedValueOnce({});

      await ruleService.deleteGroup('rg1');

      expect(mockedApi.delete).toHaveBeenCalledWith('/rule-groups/rg1');
    });
  });
});
