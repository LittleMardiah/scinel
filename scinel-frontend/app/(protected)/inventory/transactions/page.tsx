'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { InventoryHistory } from '@/components/inventory/InventoryHistory';
import { Pagination } from '@/components/common/Pagination';
import { Button } from '@/components/common/Button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Select } from '@/components/common/Select';
import { Input } from '@/components/common/Input';
import { Plus, RefreshCw, Filter, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

const ITEMS_PER_PAGE = 20;

export default function TransactionsPage() {
  // State
  const [transactions, setTransactions] = useState<any[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showFilters, setShowFilters] = useState(false);

  // Filters
  const [warehouseId, setWarehouseId] = useState('');
  const [type, setType] = useState('');
  const [dateFrom, setDateFrom] = useState('');
  const [dateTo, setDateTo] = useState('');
  const [searchSku, setSearchSku] = useState('');

  // Pagination & Sorting
  const [currentPage, setCurrentPage] = useState(1);
  const [sortKey, setSortKey] = useState('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Fetch transactions
  const fetchTransactions = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const offset = (currentPage - 1) * ITEMS_PER_PAGE;
      const params: any = {
        limit: ITEMS_PER_PAGE,
        offset,
        sort_by: sortKey,
        sort_order: sortOrder,
      };
      if (warehouseId) params.warehouse_id = warehouseId;
      if (type) params.type = type;
      if (dateFrom) params.date_from = dateFrom;
      if (dateTo) params.date_to = dateTo;
      if (searchSku) params.sku = searchSku;

      const response = await apiClient.get('/transactions', { params });
      const data = response.data.data?.transactions || [];
      const total = response.data.data?.total || 0;

      setTransactions(data);
      setTotalItems(total);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat riwayat transaksi');
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, warehouseId, type, dateFrom, dateTo, searchSku, sortKey, sortOrder]);

  // Polling (30s)
  useEffect(() => {
    fetchTransactions();
    const interval = setInterval(fetchTransactions, 30000);
    return () => clearInterval(interval);
  }, [fetchTransactions]);

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
    setType('');
    setDateFrom('');
    setDateTo('');
    setSearchSku('');
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleExportCSV = () => {
    // Simple CSV export
    if (transactions.length === 0) {
      Toast.info('Tidak ada data untuk diexport');
      return;
    }

    const headers = ['Tanggal', 'SKU', 'Tipe', 'Kuantitas', 'Gudang', 'Alasan'];
    const rows = transactions.map((t) => [
      new Date(t.created_at).toLocaleString('id-ID'),
      t.sku_code,
      t.type === 'in' ? 'Masuk' : 'Keluar',
      t.quantity,
      t.warehouse_name || t.warehouse_id,
      t.reason || '',
    ]);

    const csvContent = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n');
    const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `transactions_${new Date().toISOString().slice(0, 10)}.csv`;
    link.click();
    URL.revokeObjectURL(link.href);
    Toast.success('CSV berhasil diexport');
  };

  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Riwayat Transaksi
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Lihat semua pergerakan stok di semua gudang
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowFilters(!showFilters)}
            className="gap-2"
          >
            <Filter className="h-4 w-4" />
            Filter
          </Button>
          <Button
            variant="secondary"
            size="sm"
            onClick={fetchTransactions}
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
      {showFilters && (
        <Card>
          <CardContent className="p-4">
            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Gudang
                </label>
                <Select
                  options={[
                    { value: '', label: 'Semua Gudang' },
                    { value: 'wh-jakarta-001', label: 'Jakarta Main' },
                    { value: 'wh-surabaya-001', label: 'Surabaya Branch' },
                    { value: 'wh-bandung-001', label: 'Bandung Branch' },
                  ]}
                  value={warehouseId}
                  onChange={(e) => setWarehouseId(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Tipe
                </label>
                <Select
                  options={[
                    { value: '', label: 'Semua Tipe' },
                    { value: 'in', label: 'Masuk' },
                    { value: 'out', label: 'Keluar' },
                  ]}
                  value={type}
                  onChange={(e) => setType(e.target.value)}
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  SKU
                </label>
                <Input
                  type="text"
                  placeholder="Cari SKU..."
                  value={searchSku}
                  onChange={(e) => setSearchSku(e.target.value)}
                />
              </div>
              <div className="flex items-end gap-2">
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Tanggal Mulai
                  </label>
                  <Input
                    type="date"
                    value={dateFrom}
                    onChange={(e) => setDateFrom(e.target.value)}
                  />
                </div>
                <div className="flex-1">
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Tanggal Akhir
                  </label>
                  <Input
                    type="date"
                    value={dateTo}
                    onChange={(e) => setDateTo(e.target.value)}
                  />
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleResetFilters}
                  className="mb-0.5"
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <InventoryHistory
            transactions={transactions}
            isLoading={isLoading}
            onSort={handleSort}
            sortKey={sortKey}
            sortOrder={sortOrder}
            onExportCSV={handleExportCSV}
          />
        </CardContent>
      </Card>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
          <span>
            Menampilkan {(currentPage - 1) * ITEMS_PER_PAGE + 1} -{' '}
            {Math.min(currentPage * ITEMS_PER_PAGE, totalItems)} dari {totalItems} transaksi
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
