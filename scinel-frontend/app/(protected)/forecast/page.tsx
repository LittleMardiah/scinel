'use client';

import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { ForecastChart } from '@/components/forecast/ForecastChart';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Select } from '@/components/common/Select';
import { Input } from '@/components/common/Input';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { RefreshCw, Zap, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

const warehouseOptions = [
  { value: '', label: 'Semua Gudang' },
  { value: 'wh-jakarta-001', label: 'Jakarta Main' },
  { value: 'wh-surabaya-001', label: 'Surabaya Branch' },
  { value: 'wh-bandung-001', label: 'Bandung Branch' },
];

const skuOptions = [
  { value: '', label: 'Semua SKU' },
  { value: 'SKU-NIKE-001', label: 'SKU-NIKE-001 - Nike Shoes' },
  { value: 'SKU-NIKE-002', label: 'SKU-NIKE-002 - Nike Shirt' },
  { value: 'SKU-ADIDAS-001', label: 'SKU-ADIDAS-001 - Adidas Shoes' },
  { value: 'SKU-PUMA-001', label: 'SKU-PUMA-001 - Puma Shoes' },
];

export default function ForecastPage() {
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';

  // State
  const [data, setData] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isGenerating, setIsGenerating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [warehouseId, setWarehouseId] = useState('');
  const [skuId, setSkuId] = useState('');
  const [days, setDays] = useState(30);

  // Fetch forecasts
  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const params: any = { days };
      if (warehouseId) params.warehouse_id = warehouseId;
      if (skuId) params.sku_id = skuId;

      const response = await apiClient.get('/forecasts', { params });
      const forecastData = response.data.data?.forecasts || [];
      setData(forecastData);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data forecast');
      setData([]);
    } finally {
      setIsLoading(false);
    }
  }, [warehouseId, skuId, days]);

  // Handle generate forecast (admin only)
  const handleGenerate = async () => {
    if (!isAdmin) {
      Toast.warning('Hanya admin yang dapat generate forecast');
      return;
    }

    setIsGenerating(true);
    try {
      const response = await apiClient.post('/forecasts/generate');
      Toast.success('Batch forecast sedang diproses!');
      // Refresh data after a delay
      setTimeout(() => fetchData(), 3000);
    } catch (err: any) {
      Toast.error(err.message || 'Gagal generate forecast');
    } finally {
      setIsGenerating(false);
    }
  };

  // Polling (30s)
  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  const hasData = data && data.length > 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Demand Forecasting
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Prediksi permintaan dengan ETS & ARIMA algorithms
          </p>
        </div>
        {isAdmin && (
          <Button
            variant="primary"
            size="sm"
            onClick={handleGenerate}
            disabled={isGenerating}
            className="gap-2"
          >
            <Zap className={cn('h-4 w-4', isGenerating && 'animate-spin')} />
            {isGenerating ? 'Memproses...' : 'Generate Forecast'}
          </Button>
        )}
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Gudang
              </label>
              <Select
                options={warehouseOptions}
                value={warehouseId}
                onChange={(e) => setWarehouseId(e.target.value)}
                placeholder="Pilih gudang"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                SKU
              </label>
              <Select
                options={skuOptions}
                value={skuId}
                onChange={(e) => setSkuId(e.target.value)}
                placeholder="Pilih SKU"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Hari Forecast
              </label>
              <Input
                type="number"
                min="1"
                max="90"
                value={days}
                onChange={(e) => setDays(parseInt(e.target.value) || 30)}
              />
            </div>
            <div className="flex items-end">
              <Button
                variant="secondary"
                size="sm"
                onClick={fetchData}
                disabled={isLoading}
                className="gap-2"
              >
                <RefreshCw className={cn('h-4 w-4', isLoading && 'animate-spin')} />
                Refresh
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Visualisasi Forecast</CardTitle>
        </CardHeader>
        <CardContent>
          <ForecastChart
            data={data}
            isLoading={isLoading}
            algorithm="ETS (Exponential Smoothing)"
          />
        </CardContent>
      </Card>

      {/* Forecast Table */}
      <Card>
        <CardHeader>
          <CardTitle>Data Forecast</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => (
                <Skeleton key={i} className="h-8 w-full" />
              ))}
            </div>
          ) : error ? (
            <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
              <AlertCircle className="h-5 w-5" />
              <span>{error}</span>
            </div>
          ) : !hasData ? (
            <div className="py-8 text-center text-gray-500 dark:text-gray-400">
              {isAdmin
                ? 'Klik tombol Generate Forecast untuk membuat prediksi'
                : 'Belum ada data forecast. Hubungi admin untuk generate.'}
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200 dark:border-gray-700">
                    <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Tanggal</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">Prediksi</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">Lower Bound</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">Upper Bound</th>
                    <th className="px-4 py-2 text-center font-medium text-gray-500 dark:text-gray-400">Confidence</th>
                    <th className="px-4 py-2 text-center font-medium text-gray-500 dark:text-gray-400">Algorithm</th>
                  </tr>
                </thead>
                <tbody>
                  {data.slice(0, 30).map((item, idx) => (
                    <tr key={idx} className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50">
                      <td className="px-4 py-2 text-gray-700 dark:text-gray-300">
                        {new Date(item.forecast_date || item.date).toLocaleDateString('id-ID', {
                          day: '2-digit',
                          month: 'short',
                          year: 'numeric',
                        })}
                      </td>
                      <td className="px-4 py-2 text-right font-medium text-gray-900 dark:text-white">
                        {item.predicted_demand}
                      </td>
                      <td className="px-4 py-2 text-right text-gray-500 dark:text-gray-400">
                        {item.lower_bound ?? '-'}
                      </td>
                      <td className="px-4 py-2 text-right text-gray-500 dark:text-gray-400">
                        {item.upper_bound ?? '-'}
                      </td>
                      <td className="px-4 py-2 text-center">
                        <span className="inline-block rounded-full bg-blue-100 px-2 py-0.5 text-xs font-medium text-blue-700 dark:bg-blue-900/30 dark:text-blue-400">
                          {item.confidence ? `${(item.confidence * 100).toFixed(0)}%` : '95%'}
                        </span>
                      </td>
                      <td className="px-4 py-2 text-center text-xs text-gray-500 dark:text-gray-400">
                        {item.algorithm || 'ETS'}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {data.length > 30 && (
                <p className="mt-2 text-sm text-gray-400">
                  Menampilkan 30 dari {data.length} data forecast
                </p>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
