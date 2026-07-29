'use client';

import { useState, useEffect, useCallback } from 'react';
import { useAuth } from '@/context/AuthContext';
import { AnomalyList } from '@/components/anomaly/AnomalyList';
import { Pagination } from '@/components/common/Pagination';
import { Button } from '@/components/common/Button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Select } from '@/components/common/Select';
import { RefreshCw, Filter, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

const ITEMS_PER_PAGE = 20;

const statusOptions = [
  { value: '', label: 'Semua Status' },
  { value: 'open', label: 'Terbuka' },
  { value: 'acknowledged', label: 'Diteliti' },
  { value: 'resolved', label: 'Selesai' },
];

const severityOptions = [
  { value: '', label: 'Semua Severity' },
  { value: 'low', label: 'Rendah' },
  { value: 'medium', label: 'Sedang' },
  { value: 'high', label: 'Tinggi' },
  { value: 'critical', label: 'Kritis' },
];

export default function AnomaliesPage() {
  const { user } = useAuth();
  const isAdmin = user?.role === 'admin';

  // State
  const [data, setData] = useState<any[]>([]);
  const [totalItems, setTotalItems] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showFilters, setShowFilters] = useState(false);

  // Filters
  const [statusFilter, setStatusFilter] = useState('');
  const [severityFilter, setSeverityFilter] = useState('');

  // Pagination & Sorting
  const [currentPage, setCurrentPage] = useState(1);
  const [sortKey, setSortKey] = useState('detected_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  // Fetch anomalies
  const fetchData = useCallback(async () => {
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
      if (statusFilter) params.status = statusFilter;
      if (severityFilter) params.severity = severityFilter;

      const response = await apiClient.get('/anomalies', { params });
      const anomalyData = response.data.data?.anomalies || [];
      const total = response.data.data?.total || 0;

      setData(anomalyData);
      setTotalItems(total);
    } catch (err: any) {
      setError(err.message || 'Gagal memuat data anomali');
    } finally {
      setIsLoading(false);
    }
  }, [currentPage, statusFilter, severityFilter, sortKey, sortOrder]);

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
    setStatusFilter('');
    setSeverityFilter('');
    setCurrentPage(1);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  const handleAcknowledge = async (id: string) => {
    try {
      await apiClient.put(`/anomalies/${id}`, { status: 'acknowledged' });
      Toast.success('Anomali sedang diteliti');
      fetchData();
    } catch (err: any) {
      Toast.error(err.message || 'Gagal mengupdate status');
    }
  };

  const handleResolve = async (id: string) => {
    try {
      await apiClient.put(`/anomalies/${id}`, { 
        status: 'resolved',
        resolution_notes: 'Anomali telah diselesaikan'
      });
      Toast.success('Anomali selesai');
      fetchData();
    } catch (err: any) {
      Toast.error(err.message || 'Gagal mengupdate status');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Yakin ingin menghapus anomali ini?')) return;
    try {
      await apiClient.delete(`/anomalies/${id}`);
      Toast.success('Anomali dihapus');
      fetchData();
    } catch (err: any) {
      Toast.error(err.message || 'Gagal menghapus anomali');
    }
  };

  const totalPages = Math.ceil(totalItems / ITEMS_PER_PAGE);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Deteksi Anomali
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Pantau kejadian tidak biasa di supply chain
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
            onClick={fetchData}
            disabled={isLoading}
            className="gap-2"
          >
            <RefreshCw className={cn('h-4 w-4', isLoading && 'animate-spin')} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Filters */}
      {showFilters && (
        <Card>
          <CardContent className="p-4">
            <div className="grid gap-4 sm:grid-cols-2">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Status
                </label>
                <Select
                  options={statusOptions}
                  value={statusFilter}
                  onChange={(e) => setStatusFilter(e.target.value)}
                  placeholder="Pilih status"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Severity
                </label>
                <Select
                  options={severityOptions}
                  value={severityFilter}
                  onChange={(e) => setSeverityFilter(e.target.value)}
                  placeholder="Pilih severity"
                />
              </div>
            </div>
            <div className="mt-4 flex justify-end">
              <Button
                variant="ghost"
                size="sm"
                onClick={handleResetFilters}
                className="gap-2"
              >
                <X className="h-4 w-4" />
                Reset Filter
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <AnomalyList
            data={data}
            isLoading={isLoading}
            onSort={handleSort}
            sortKey={sortKey}
            sortOrder={sortOrder}
            onAcknowledge={handleAcknowledge}
            onResolve={handleResolve}
            onDelete={handleDelete}
            isAdmin={isAdmin}
          />
        </CardContent>
      </Card>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex flex-col sm:flex-row justify-between items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
          <span>
            Menampilkan {(currentPage - 1) * ITEMS_PER_PAGE + 1} -{' '}
            {Math.min(currentPage * ITEMS_PER_PAGE, totalItems)} dari {totalItems} anomali
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
