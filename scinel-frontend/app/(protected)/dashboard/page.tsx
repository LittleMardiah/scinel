'use client';

import { useDashboardData } from '@/hooks/useDashboardData';
import { KPIGrid } from '@/components/dashboard/KPIGrid';
import { RecentTransactions } from '@/components/dashboard/RecentTransactions';
import { ForecastChart } from '@/components/dashboard/ForecastChart';
import { SystemStatus } from '@/components/dashboard/SystemStatus';
import { QuickActionsCard } from '@/components/dashboard/QuickActionsCard';
import { Loading } from '@/components/common/Loading';

export default function DashboardPage() {
  const { data, isLoading, error } = useDashboardData();

  if (isLoading || !data) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <Loading size="lg" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center text-red-600">
        <p>Gagal memuat data: {error.message}</p>
      </div>
    );
  }

  const { total_sku, total_warehouse, low_stock_count, anomaly_count, health, recent_transactions, forecast } = data;
  const healthStatus = health?.status === 'online' ? 'online' : 'offline';

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
      <KPIGrid totalSku={total_sku || 0} totalWarehouse={total_warehouse || 0} lowStockCount={low_stock_count || 0} anomalyCount={anomaly_count || 0} />
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2"><ForecastChart data={forecast || []} /></div>
        <div className="space-y-6"><SystemStatus status={healthStatus} /><QuickActionsCard /></div>
      </div>
      <RecentTransactions transactions={recent_transactions || []} />
    </div>
  );
}
