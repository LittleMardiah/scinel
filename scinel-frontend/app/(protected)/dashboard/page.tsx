'use client';

import { useDashboardData } from '@/hooks/useDashboardData';
import { StatsCard } from '@/components/dashboard/StatsCard';
import { RecentTransactions } from '@/components/dashboard/RecentTransactions';
import { ForecastChart } from '@/components/dashboard/ForecastChart';
import { SystemStatus } from '@/components/dashboard/SystemStatus';
import { QuickActionsCard } from '@/components/dashboard/QuickActionsCard';
import { Skeleton } from '@/components/common/Skeleton';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { 
  Package, 
  Warehouse, 
  AlertTriangle, 
  Activity, 
  RefreshCw 
} from 'lucide-react';

type HealthStatus = 'online' | 'offline' | 'loading';

export default function DashboardPage() {
  const { data, isLoading, error, refetch } = useDashboardData();

  // Ensure health.status is one of the allowed types
  const healthStatus: HealthStatus = 
    data.health.status === 'online' || data.health.status === 'offline' 
      ? data.health.status 
      : 'loading';

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i}>
              <CardContent className="p-6">
                <Skeleton className="h-4 w-1/2" />
                <Skeleton className="mt-2 h-8 w-1/3" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Skeleton className="h-64 w-full" />
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-48 w-full" />
          <Skeleton className="h-48 w-full" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
          <p className="text-red-600 dark:text-red-400">Gagal memuat dashboard</p>
          <p className="mt-2 text-sm text-gray-500">{error}</p>
          <button
            onClick={() => refetch()}
            className="mt-4 inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >
            <RefreshCw className="h-4 w-4" />
            Coba Lagi
          </button>
        </div>
      </div>
    );
  }

  const { kpi, transactions, forecasts } = data;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
          Dashboard
        </h1>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Ringkasan supply chain Anda secara real-time
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatsCard
          title="Total SKU"
          value={kpi.totalSku}
          icon={Package}
          trend={{ value: 12, isPositive: true }}
        />
        <StatsCard
          title="Total Gudang"
          value={kpi.totalWarehouse}
          icon={Warehouse}
        />
        <StatsCard
          title="Stok Rendah"
          value={kpi.lowStock}
          icon={AlertTriangle}
          trend={{ value: 5, isPositive: false }}
          variant="warning"
        />
        <StatsCard
          title="Anomali Aktif"
          value={kpi.anomalies}
          icon={Activity}
          variant="danger"
        />
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <div className="md:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle>Forecast Permintaan</CardTitle>
            </CardHeader>
            <CardContent>
              <ForecastChart data={forecasts} />
            </CardContent>
          </Card>
        </div>
        <div>
          <SystemStatus status={healthStatus} timestamp={data.health.timestamp} />
          <div className="mt-4">
            <QuickActionsCard />
          </div>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Transaksi Terbaru</CardTitle>
        </CardHeader>
        <CardContent>
          <RecentTransactions transactions={transactions} />
        </CardContent>
      </Card>
    </div>
  );
}
