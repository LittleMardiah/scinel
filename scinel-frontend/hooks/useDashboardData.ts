import { useState, useEffect, useCallback } from 'react';
import apiClient from '@/lib/api';

interface DashboardData {
  inventory: any[];
  transactions: any[];
  forecasts: any[];
  health: { status: string; timestamp?: string };
  kpi: {
    totalSku: number;
    totalWarehouse: number;
    lowStock: number;
    anomalies: number;
  };
}

const initialData: DashboardData = {
  inventory: [],
  transactions: [],
  forecasts: [],
  health: { status: 'loading' },
  kpi: { totalSku: 0, totalWarehouse: 0, lowStock: 0, anomalies: 0 },
};

export function useDashboardData(pollingInterval: number = 30000) {
  const [data, setData] = useState<DashboardData>(initialData);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);

      // Fetch Inventory (for KPI)
      const invRes = await apiClient.get('/inventory?limit=1');
      const invItems = invRes.data.data?.inventory || [];

      // Fetch Transactions (recent)
      const txnRes = await apiClient.get('/transactions?limit=5');
      const txnItems = txnRes.data.data?.transactions || [];

      // Fetch Forecasts (for chart)
      const forecastRes = await apiClient.get('/forecasts?days=30');
      const forecastItems = forecastRes.data.data?.forecasts || [];

      // Fetch Health
      let healthStatus = 'offline';
      try {
        const healthRes = await apiClient.get('/health');
        if (healthRes.data?.status === 'healthy' || healthRes.status === 200) {
          healthStatus = 'online';
        }
      } catch {
        healthStatus = 'offline';
      }

      // Calculate KPI
      const totalSku = invItems.length > 0 ? 125 : 0; // Placeholder, nanti dihitung dari data real
      const totalWarehouse = 3; // Placeholder
      const lowStock = invItems.filter((i: any) => i.status === 'low_stock').length || 5;
      const anomalies = 2; // Placeholder

      setData({
        inventory: invItems,
        transactions: txnItems,
        forecasts: forecastItems,
        health: { status: healthStatus, timestamp: new Date().toISOString() },
        kpi: {
          totalSku: totalSku || 124,
          totalWarehouse: 3,
          lowStock: lowStock || 4,
          anomalies: anomalies || 2,
        },
      });
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data dashboard');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();

    const interval = setInterval(fetchData, pollingInterval);

    return () => clearInterval(interval);
  }, [fetchData, pollingInterval]);

  return { data, isLoading, error, refetch: fetchData };
}
