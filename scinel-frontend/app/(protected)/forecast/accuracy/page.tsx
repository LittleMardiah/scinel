'use client';

import { useState, useEffect, useCallback } from 'react';
import { AccuracyChart } from '@/components/forecast/AccuracyChart';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Select } from '@/components/common/Select';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { RefreshCw, AlertCircle, TrendingUp, TrendingDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Toast } from '@/components/common/Toast';
import { apiClient } from '@/lib/api';

const warehouseOptions = [
  { value: '', label: 'Semua Gudang' },
  { value: 'wh-jakarta-001', label: 'Jakarta Main' },
  { value: 'wh-surabaya-001', label: 'Surabaya Branch' },
  { value: 'wh-bandung-001', label: 'Bandung Branch' },
];

const algorithmOptions = [
  { value: '', label: 'Semua Algoritma' },
  { value: 'exponential_smoothing', label: 'ETS (Exponential Smoothing)' },
  { value: 'arima', label: 'ARIMA' },
];

export default function ForecastAccuracyPage() {
  // State
  const [data, setData] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters
  const [warehouseId, setWarehouseId] = useState('');
  const [algorithm, setAlgorithm] = useState('');

  // Fetch accuracy data
  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const params: any = {};
      if (warehouseId) params.warehouse_id = warehouseId;
      if (algorithm) params.algorithm = algorithm;

      const response = await apiClient.get('/forecasts/accuracy', { params });
      const accuracyData = response.data.data?.accuracy || [];
      setData(accuracyData);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data accuracy');
      setData([]);
    } finally {
      setIsLoading(false);
    }
  }, [warehouseId, algorithm]);

  // Polling (30s)
  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  // Best algorithm
  const best = data.length > 0
    ? data.reduce((a, b) => a.mape < b.mape ? a : b)
    : null;

  // Format data for chart
  const chartData = data.map((item) => ({
    algorithm: item.algorithm === 'exponential_smoothing' ? 'ETS' : 'ARIMA',
    mape: item.mape || 0,
    rmse: item.rmse || 0,
    mae: item.mae || 0,
  }));

  const hasData = data && data.length > 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Akurasi Forecast
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Perbandingan akurasi algoritma ETS vs ARIMA
          </p>
        </div>
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

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="grid gap-4 sm:grid-cols-2">
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
                Algoritma
              </label>
              <Select
                options={algorithmOptions}
                value={algorithm}
                onChange={(e) => setAlgorithm(e.target.value)}
                placeholder="Pilih algoritma"
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Best Algorithm Summary */}
      {best && (
        <Card className="border-green-200 dark:border-green-800">
          <CardContent className="p-4">
            <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4">
              <div className="flex items-center gap-2">
                <TrendingUp className="h-5 w-5 text-green-600 dark:text-green-400" />
                <span className="font-medium text-gray-900 dark:text-white">
                  Best Algorithm:
                </span>
                <span className="font-bold text-green-600 dark:text-green-400">
                  {best.algorithm === 'exponential_smoothing' ? 'ETS' : 'ARIMA'}
                </span>
              </div>
              <div className="flex flex-wrap gap-4 text-sm">
                <span className="text-gray-500 dark:text-gray-400">
                  MAPE: <strong className="text-gray-900 dark:text-white">{best.mape?.toFixed(2)}%</strong>
                </span>
                <span className="text-gray-500 dark:text-gray-400">
                  RMSE: <strong className="text-gray-900 dark:text-white">{best.rmse?.toFixed(2)}</strong>
                </span>
                <span className="text-gray-500 dark:text-gray-400">
                  MAE: <strong className="text-gray-900 dark:text-white">{best.mae?.toFixed(2)}</strong>
                </span>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Accuracy Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Perbandingan Akurasi</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <Skeleton className="h-48 w-full" />
          ) : error ? (
            <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
              <AlertCircle className="h-5 w-5" />
              <span>{error}</span>
            </div>
          ) : !hasData ? (
            <div className="py-8 text-center text-gray-500 dark:text-gray-400">
              Belum ada data accuracy. Generate forecast terlebih dahulu.
            </div>
          ) : (
            <div className="space-y-6">
              <AccuracyChart
                data={chartData}
                metric="mape"
                height={250}
              />
              <AccuracyChart
                data={chartData}
                metric="rmse"
                height={200}
              />
              <AccuracyChart
                data={chartData}
                metric="mae"
                height={200}
              />
            </div>
          )}
        </CardContent>
      </Card>

      {/* Accuracy Table */}
      <Card>
        <CardHeader>
          <CardTitle>Detail Metrik</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              {[...Array(3)].map((_, i) => (
                <Skeleton key={i} className="h-8 w-full" />
              ))}
            </div>
          ) : !hasData ? (
            <div className="py-8 text-center text-gray-500 dark:text-gray-400">
              Tidak ada data
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-200 dark:border-gray-700">
                    <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Algoritma</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">MAPE (%)</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">RMSE</th>
                    <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">MAE</th>
                    <th className="px-4 py-2 text-center font-medium text-gray-500 dark:text-gray-400">Tanggal Evaluasi</th>
                    <th className="px-4 py-2 text-center font-medium text-gray-500 dark:text-gray-400">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {data.map((item, idx) => {
                    const isBest = best && item.mape === best.mape;
                    const mapeValue = item.mape || 0;

                    return (
                      <tr key={idx} className={cn(
                        'border-b border-gray-100 dark:border-gray-800',
                        isBest && 'bg-green-50 dark:bg-green-900/10'
                      )}>
                        <td className="px-4 py-2 font-medium text-gray-900 dark:text-white">
                          {item.algorithm === 'exponential_smoothing' ? 'ETS' : 'ARIMA'}
                          {isBest && (
                            <span className="ml-2 inline-block text-xs text-green-600 dark:text-green-400">🏆 Best</span>
                          )}
                        </td>
                        <td className={cn(
                          'px-4 py-2 text-right font-medium',
                          mapeValue < 10 ? 'text-green-600 dark:text-green-400' :
                          mapeValue < 20 ? 'text-yellow-600 dark:text-yellow-400' :
                          'text-red-600 dark:text-red-400'
                        )}>
                          {mapeValue.toFixed(2)}%
                        </td>
                        <td className="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                          {item.rmse?.toFixed(2) || '-'}
                        </td>
                        <td className="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                          {item.mae?.toFixed(2) || '-'}
                        </td>
                        <td className="px-4 py-2 text-center text-gray-500 dark:text-gray-400">
                          {item.evaluation_date
                            ? new Date(item.evaluation_date).toLocaleDateString('id-ID')
                            : '-'}
                        </td>
                        <td className="px-4 py-2 text-center">
                          <span className={cn(
                            'inline-block rounded-full px-2 py-0.5 text-xs font-medium',
                            mapeValue < 10
                              ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                              : mapeValue < 20
                              ? 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400'
                              : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                          )}>
                            {mapeValue < 10 ? 'Baik' : mapeValue < 20 ? 'Sedang' : 'Perlu Perbaikan'}
                          </span>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Interpretation Guide */}
      <Card className="border-blue-200 dark:border-blue-800">
        <CardContent className="p-4">
          <p className="text-sm text-gray-600 dark:text-gray-300">
            <span className="font-medium">📊 Interpretasi:</span>
            {' '}
            <span className="text-green-600 dark:text-green-400">MAPE &lt; 10%</span> = Excellent,
            <span className="text-yellow-600 dark:text-yellow-400"> 10-20%</span> = Good,
            <span className="text-red-600 dark:text-red-400"> &gt; 20%</span> = Needs Improvement
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
