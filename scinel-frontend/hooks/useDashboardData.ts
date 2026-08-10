import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api';

const MOCK_INVENTORY = {
  total_sku: 5,
  total_warehouse: 3,
  low_stock_count: 2,
  anomaly_count: 1,
  health: { status: 'online' },
  recent_transactions: [
    { id: 'txn-1', sku: 'SKU-001', type: 'in', quantity: 50, created_at: new Date().toISOString() }
  ],
  forecast: []
};

export function useDashboardData() {
  const [data, setData] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await apiClient.get('/inventory?limit=1');
        setData(res.data.data);
      } catch (err: any) {
        if (err.response?.status === 404) {
          console.warn('⚠️ [DEBUG] Endpoint /inventory 404, using mock data.');
          setData(MOCK_INVENTORY);
        } else {
          setError(err);
        }
      } finally {
        setIsLoading(false);
      }
    };
    fetchData();
  }, []);

  return { data, isLoading, error };
}
