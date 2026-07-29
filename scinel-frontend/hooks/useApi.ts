import { useState, useCallback, useEffect } from 'react';
import apiClient from '@/lib/api';

interface UseApiOptions {
  autoFetch?: boolean;
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

export function useApi<T = any>(url: string, options: UseApiOptions = {}) {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<any>(null);

  const fetch = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(url);
      setData(response.data.data);
      options.onSuccess?.(response.data.data);
    } catch (err) {
      setError(err);
      options.onError?.(err);
    } finally {
      setIsLoading(false);
    }
  }, [url, options]);

  useEffect(() => {
    if (options.autoFetch !== false) fetch();
  }, [fetch, options.autoFetch]);

  return { data, isLoading, error, refetch: fetch };
}
