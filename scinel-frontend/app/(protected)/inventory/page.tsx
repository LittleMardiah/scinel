'use client';
import { cn } from "@/lib/utils";

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { InventoryTable } from '@/components/inventory/InventoryTable';
import { InventoryFilters } from '@/components/inventory/InventoryFilters';
import { Pagination } from '@/components/common/Pagination';
import { Button } from '@/components/common/Button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Plus, RefreshCw } from 'lucide-react';
import { apiClient } from '@/lib/api';

const ITEMS_PER_PAGE = 20;

export default function InventoryPage() {
  const router = useRouter();

  // State
  const [data, setData] = useState<any[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Filters & Pagination
  const [currentPage, setCurrentPage] = useState(1);
  const [warehouseId, setWarehouseId] = useState('');
  const [status, setStatus] = useState('');
  const [sortKey, setSortKey] = useState('sku_code');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  // Fetch data
  const fetchData = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const offset = (currentPage - 1) * ITEMS_PER_PAGE;
      const params: any = {
        limit: ITEMS_PER_PAGE,
        offset,
      };
      if (warehouseId) params.warehouse_id = warehouseId;
      if (status) params.status = status;

      const response = await apiClient.get('/inventory', { params });
      const inventoryData = response.data.data?.inventory || [];
      const total = response.data.data?.total || 0;

      // Client-side sorting (backend may not support sort yet)
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
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data inventaris');
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, warehouseId, status, sortKey, sortOrder]);

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

  const handleResetFilters = () => {
    setWarehouseId('');
    setStatus('');
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Manajemen Inventaris
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Kelola stok barang di semua gudang
          </p>
        </div>
        <div className="flex gap-2">
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
          <Link href="/inventory/transaction/new">
            <Button size="sm" className="gap-2">
              <Plus className="h-4 w-4" />
              Transaksi Baru
            </Button>
          </Link>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <InventoryFilters
            warehouseId={warehouseId}
            status={status}
            onWarehouseChange={setWarehouseId}
            onStatusChange={setStatus}
            onReset={handleResetFilters}
          />
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <InventoryTable
            data={data}
            isLoading={isLoading}
            onSort={handleSort}
            sortKey={sortKey}
            sortOrder={sortOrder}
            onRowClick={(id: string) => router.push(`/inventory/${id}`)}
          />
        </CardContent>
      </Card>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-between items-center text-sm text-gray-500 dark:text-gray-400">
          <span>
            Menampilkan {(currentPage - 1) * ITEMS_PER_PAGE + 1} -{' '}
            {Math.min(currentPage * ITEMS_PER_PAGE, totalItems)} dari {totalItems} item
          </span>
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
          />
        </div>
      )}

      {error && (
        <div className="rounded-md bg-red-50 p-4 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}
