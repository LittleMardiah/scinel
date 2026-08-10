'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { RefreshCw } from 'lucide-react';
import { cn } from '@/lib/utils';
import { apiClient } from '@/lib/api';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
} from 'recharts';

const COLORS = ['#22c55e', '#f59e0b', '#ef4444', '#dc2626'];

export default function AnomalyAnalyticsPage() {
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [trendData, setTrendData] = useState<any[]>([]);
  const [severityData, setSeverityData] = useState<any[]>([]);
  const [topSkus, setTopSkus] = useState<any[]>([]);
  const [insights, setInsights] = useState<any>({});

  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      // Fetch all anomalies (no filters)
      const response = await apiClient.get('/anomalies', {
        params: { limit: 1000, offset: 0 },
      });
      const anomalies = response.data.data?.anomalies || [];

      // Process data for charts
      processAnomalyData(anomalies);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data analitik');
    } finally {
      setIsLoading(false);
    }
  }, []);

  const processAnomalyData = (anomalies: any[]) => {
    // Severity breakdown
    const severityCount: Record<string, number> = {};
    anomalies.forEach((a) => {
      severityCount[a.severity] = (severityCount[a.severity] || 0) + 1;
    });
    const severityData = Object.entries(severityCount).map(([name, value]) => ({
      name: name.charAt(0).toUpperCase() + name.slice(1),
      value,
    }));
    setSeverityData(severityData);

    // Top affected SKUs
    const skuCount: Record<string, number> = {};
    anomalies.forEach((a) => {
      const key = a.sku_code || 'Unknown';
      skuCount[key] = (skuCount[key] || 0) + 1;
    });
    const topSkus = Object.entries(skuCount)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5)
      .map(([name, count]) => ({ name, count }));
    setTopSkus(topSkus);

    // Trend by severity over time (last 30 days)
    const days = 30;
    const today = new Date();
    const trendMap: Record<string, any> = {};
    for (let i = days - 1; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);
      const key = date.toISOString().slice(0, 10);
      trendMap[key] = { date: key, low: 0, medium: 0, high: 0, critical: 0 };
    }

    anomalies.forEach((a) => {
      const date = new Date(a.detected_at).toISOString().slice(0, 10);
      if (trendMap[date]) {
        const severity = a.severity || 'low';
        if (trendMap[date][severity] !== undefined) {
          trendMap[date][severity] += 1;
        }
      }
    });

    setTrendData(Object.values(trendMap));

    // Insights
    const total = anomalies.length;
    const resolved = anomalies.filter((a) => a.status === 'resolved').length;
    const open = anomalies.filter((a) => a.status === 'open').length;
    const critical = anomalies.filter((a) => a.severity === 'critical').length;

    setInsights({
      total,
      resolved,
      open,
      critical,
      resolutionRate: total > 0 ? (resolved / total) * 100 : 0,
    });
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-24 w-full" />
          ))}
        </div>
        <Skeleton className="h-64 w-full" />
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-64 w-full" />
          <Skeleton className="h-64 w-full" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
          <p className="text-red-600 dark:text-red-400">Gagal memuat data</p>
          <p className="mt-2 text-sm text-gray-500">{error}</p>
          <Button variant="primary" onClick={fetchData} className="mt-4">
            Coba Lagi
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Analitik Anomali
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Tren dan statistik deteksi anomali
          </p>
        </div>
        <Button
          variant="secondary"
          size="sm"
          onClick={fetchData}
          className="gap-2"
        >
          <RefreshCw className="h-4 w-4" />
          Refresh
        </Button>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">Total Anomali</p>
            <p className="text-2xl font-bold text-gray-900 dark:text-white">{insights.total || 0}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">Selesai</p>
            <p className="text-2xl font-bold text-green-600 dark:text-green-400">{insights.resolved || 0}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">Terbuka</p>
            <p className="text-2xl font-bold text-red-600 dark:text-red-400">{insights.open || 0}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-gray-500 dark:text-gray-400">Resolution Rate</p>
            <p className="text-2xl font-bold text-blue-600 dark:text-blue-400">
              {insights.resolutionRate?.toFixed(1) || 0}%
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Trend Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Tren Anomali (30 Hari)</CardTitle>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart data={trendData}>
              <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
              <XAxis dataKey="date" className="text-xs text-gray-500 dark:text-gray-400" />
              <YAxis className="text-xs text-gray-500 dark:text-gray-400" />
              <Tooltip />
              <Legend />
              <Area type="monotone" dataKey="low" name="Rendah" fill="#22c55e" stroke="#22c55e" fillOpacity={0.2} />
              <Area type="monotone" dataKey="medium" name="Sedang" fill="#f59e0b" stroke="#f59e0b" fillOpacity={0.2} />
              <Area type="monotone" dataKey="high" name="Tinggi" fill="#ef4444" stroke="#ef4444" fillOpacity={0.2} />
              <Area type="monotone" dataKey="critical" name="Kritis" fill="#dc2626" stroke="#dc2626" fillOpacity={0.2} />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Severity & Top SKUs */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Distribusi Severity</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={250}>
              <PieChart>
                <Pie
                  data={severityData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  label
                >
                  {severityData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Top 5 SKU Paling Bermasalah</CardTitle>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={250}>
              <BarChart data={topSkus}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
                <XAxis dataKey="name" className="text-xs text-gray-500 dark:text-gray-400" />
                <YAxis className="text-xs text-gray-500 dark:text-gray-400" />
                <Tooltip />
                <Bar dataKey="count" fill="#0ea5e9" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      {/* Insights */}
      <Card>
        <CardHeader>
          <CardTitle>Insights</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid gap-2 text-sm text-gray-600 dark:text-gray-300">
            {insights.total > 0 ? (
              <>
                <p>
                  • Total <strong>{insights.total}</strong> anomali terdeteksi.
                </p>
                <p>
                  • <strong>{insights.resolutionRate?.toFixed(1)}%</strong> anomali telah diselesaikan.
                </p>
                <p>
                  • <strong>{insights.critical || 0}</strong> anomali kritis membutuhkan perhatian segera.
                </p>
                {topSkus.length > 0 && (
                  <p>
                    • SKU paling bermasalah: <strong>{topSkus[0]?.name}</strong> ({topSkus[0]?.count} anomali)
                  </p>
                )}
              </>
            ) : (
              <p>Belum ada data yang cukup untuk insights.</p>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
