import { renderHook, waitFor } from '@testing-library/react';
import { useApi } from '@/hooks/useApi';
import { apiClient } from '@/lib/api';
import { describe, it, expect, jest, beforeEach, afterEach } from '@jest/globals';

describe('useApi', () => {
  let getSpy: jest.SpyInstance;

  beforeEach(() => {
    getSpy = jest.spyOn(apiClient, 'get').mockResolvedValue({ data: { data: {} } });
  });

  afterEach(() => {
    getSpy.mockRestore();
  });

  it('should fetch data successfully', async () => {
    const mockData = { data: 'test data' };
    getSpy.mockResolvedValueOnce({ data: { data: mockData } });

    const { result } = renderHook(() => useApi('/test', { autoFetch: true }));
    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.data).toEqual(mockData);
    });
  });

  it('should handle error', async () => {
    const mockError = new Error('Network error');
    getSpy.mockRejectedValueOnce(mockError);

    const { result } = renderHook(() => useApi('/test', { autoFetch: true }));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeDefined();
    });
  });

  it('should call onSuccess callback', async () => {
    const onSuccess = jest.fn();
    const mockData = { data: 'test' };
    getSpy.mockResolvedValueOnce({ data: { data: mockData } });

    renderHook(() => useApi('/test', { autoFetch: true, onSuccess }));

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalledWith(mockData);
    });
  });

  it('should support refetch', async () => {
    const mockData = { data: 'test' };
    getSpy.mockResolvedValue({ data: { data: mockData } });

    const { result } = renderHook(() => useApi('/test', { autoFetch: false }));

    await result.current.refetch();

    await waitFor(() => {
      expect(result.current.data).toEqual(mockData);
    });
  });
});
