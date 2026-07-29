'use client';

import { useState, useEffect, useCallback } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { InventoryTable } from '@/components/inventory/InventoryTable';
import { Breadcrumb } from '@/components/common/Breadcrumb';
import { Card, CardContent } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { Package, AlertTriangle, Clock, Plus } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

// Warehouse name mapping (would come from API in production)
const WAREHOUSE_NAMES: Record<string, string> = {
  'wh-jakarta-001': 'Jakarta Main',
  'wh-surabaya-001': 'Surabaya Branch',
  'wh-bandung-001': 'Bandung Branch',
};

const ITEMS_PER_PAGE = 20;

export default function WarehouseInventoryPage() {
  const params = useParams();
  const router = useRouter();
  const warehouseId = params?.warehouse as string;

  // State
  const [data, setData] = useState<any[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<string | null>(null);

  // Pagination & Sorting
  const [currentPage, setCurrentPage] = useState(1);
  const [sortKey, setSortKey] = useState('sku_code');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  // Stats
  const [totalSku, setTotalSku] = useState(0);
  const [lowStockCount, setLowStockCount] = useState(0);

  const warehouseName = WAREHOUSE_NAMES[warehouseId] || warehouseId || 'Gudang';

  // Fetch inventory for this warehouse
  const fetchData = useCallback(async () => {
    if (!warehouseId) return;

    setIsLoading(true);
    setError(null);
    try {
      const offset = (currentPage - 1) * ITEMS_PER_PAGE;
      const params: any = {
        warehouse_id: warehouseId,
        limit: ITEMS_PER_PAGE,
        offset,
      };

      const response = await apiClient.get('/inventory', { params });
      const inventoryData = response.data.data?.inventory || [];
      const total = response.data.data?.total || 0;

      // Client-side sorting
      const sorted = [...inventoryData].sort((a, b) => {
        let valA = a[sortKey] ?? '';
        let valB = b[sortKey] ?? '';
        if (typeof valA === 'string') valA = valA.toLowerCase();
        if (typeof valB === 'string') valB = valB.toLowerCase();
        if (valA < valB) return sortOrder === 'asc' ? -1 : 1;
        if (valA > valB) return sortOrder === 'asc' ? 1 : -1;
        return 0;
      });

      setData(sorted);
      setTotalItems(total);

      // Calculate stats
      const skuCount = inventoryData.length;
      const lowStock = inventoryData.filter((item: any) => 
        item.status === 'low_stock' || item.status === 'critical'
      ).length;
      setTotalSku(skuCount);
      setLowStockCount(lowStock);

      // Last updated
      if (inventoryData.length > 0) {
        const latest = inventoryData.reduce((a: any, b: any) => 
          new Date(a.updated_at) > new Date(b.updated_at) ? a : b
        );
        setLastUpdated(latest.updated_at);
      } else {
        setLastUpdated(null);
      }
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data inventaris');
    } finally {
      setIsLoading(false);
    }
  }, [warehouseId, currentPage, sortKey, sortOrder]);

  // Polling (30s)
  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  // Handlers
  const handleSort = (key: string) => {
    if (sortKey === key) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortKey(key);
      setSortOrder('asc');
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleRowClick = (id: string) => {
    router.push(`/inventory/${id}`);
  };

  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

  // Breadcrumb items
  const breadcrumbItems = [
    { label: 'Dashboard', href: '/dashboard' },
    { label: 'Inventory', href: '/inventory' },
    { label: warehouseName },
  ];

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-8 w-48" />
        </div>
        <div className="grid gap-4 md:grid-cols-3">
          {[...Array(3)].map((_, i) => (
            <Card key={i}>
              <CardContent className="p-6">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="mt-2 h-8 w-16" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
          <p className="text-red-600 dark:text-red-400">Gagal memuat data</p>
          <p className="mt-2 text-sm text-gray-500">{error}</p>
          <Button
            variant="primary"
            size="sm"
            onClick={fetchData}
            className="mt-4"
          >
            Coba Lagi
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Breadcrumb */}
      <Breadcrumb items={breadcrumbItems} />

      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {warehouseName}
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Manajemen stok untuk gudang ini
          </p>
        </div>
        <Link
          href={`/inventory/transaction/new?warehouse=${warehouseId}`}
        >
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            Transaksi Baru
          </Button>
        </Link>
      </div>

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Total SKU
                </p>
                <p className="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
                  {totalSku}
                </p>
              </div>
              <div className="rounded-full bg-blue-100 p-2 dark:bg-blue-900/30">
                <Package className="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Stok Rendah / Kritis
                </p>
                <p className="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
                  {lowStockCount}
                </p>
              </div>
              <div className="rounded-full bg-yellow-100 p-2 dark:bg-yellow-900/30">
                <AlertTriangle className="h-5 w-5 text-yellow-600 dark:text-yellow-400" />
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  Terakhir Diperbarui
                </p>
                <p className="mt-2 text-sm text-gray-900 dark:text-white">
                  {lastUpdated
                    ? new Date(lastUpdated).toLocaleString('id-ID', {
                        day: '2-digit',
                        month: 'short',
                        hour: '2-digit',
                        minute: '2-digit',
                      })
                    : 'Belum ada data'
                  }
                </p>
              </div>
              <div className="rounded-full bg-gray-100 p-2 dark:bg-gray-700">
                <Clock className="h-5 w-5 text-gray-600 dark:text-gray-400" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <InventoryTable
            data={data}
            isLoading={isLoading}
            onSort={handleSort}
            sortKey={sortKey}
            sortOrder={sortOrder}
            onRowClick={handleRowClick}
          />
        </CardContent>
      </Card>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
          <span>
            Menampilkan {(currentPage - 1) * ITEMS_PER_PAGE + 1} -{' '}
            {Math.min(currentPage * ITEMS_PER_PAGE, totalItems)} dari {totalItems} item
          </span>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handlePageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="h-8 px-3"
            >
              &lt;
            </Button>
            {[...Array(Math.min(totalPages, 5))].map((_, i) => {
              const page = i + 1;
              return (
                <Button
                  key={page}
                  variant={page === currentPage ? 'primary' : 'ghost'}
                  size="sm"
                  onClick={() => handlePageChange(page)}
                  className="h-8 w-8 p-0"
                >
                  {page}
                </Button>
              );
            })}
            {totalPages > 5 && (
              <span className="px-2 text-gray-400">...</span>
            )}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handlePageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="h-8 px-3"
            >
              &gt;
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
