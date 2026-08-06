import { useState, useCallback, useEffect, useRef } from "react";
import { apiClient } from "@/lib/api";

interface UseApiOptions {
  autoFetch?: boolean;
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

export function useApi<T = any>(url: string, options: UseApiOptions = {}) {
  const [data, setData] = useState<T | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<any>(null);

  // Simpan options di ref agar selalu terbaru tanpa trigger re-render
  const optionsRef = useRef(options);
  useEffect(() => {
    optionsRef.current = options;
  }, [options]);

  const fetch = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.get(url);
      setData(response.data.data);
      optionsRef.current.onSuccess?.(response.data.data);
    } catch (err) {
      setError(err);
      optionsRef.current.onError?.(err);
    } finally {
      setIsLoading(false);
    }
  }, [url]); // Hanya bergantung pada url, bukan options

  useEffect(() => {
    if (optionsRef.current.autoFetch !== false) {
      fetch();
    }
  }, [fetch]);

  return { data, isLoading, error, refetch: fetch };
}
