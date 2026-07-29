'use client';

import React, { useState } from 'react';
import { Badge } from '@/components/ui/Badge';
import { Skeleton } from '@/components/common/Skeleton';
import { Button } from '@/components/common/Button';
import { Download, ChevronUp, ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';

interface Transaction {
  id: string;
  sku_code: string;
  type: 'in' | 'out';
  quantity: number;
  reason?: string;
  warehouse_id: string;
  warehouse_name?: string;
  created_by?: string;
  created_at: string;
  previous_quantity?: number;
  new_quantity?: number;
}

interface InventoryHistoryProps {
  transactions: Transaction[];
  isLoading?: boolean;
  onSort?: (key: string) => void;
  sortKey?: string;
  sortOrder?: 'asc' | 'desc';
  onExportCSV?: () => void;
}

export const InventoryHistory: React.FC<InventoryHistoryProps> = ({
  transactions,
  isLoading,
  onSort,
  sortKey = 'created_at',
  sortOrder = 'desc',
  onExportCSV,
}) => {
  const [expandedRow, setExpandedRow] = useState<string | null>(null);

  const headers = [
    { key: 'created_at', label: 'Tanggal & Waktu' },
    { key: 'sku_code', label: 'Kode SKU' },
    { key: 'type', label: 'Tipe' },
    { key: 'quantity', label: 'Kuantitas' },
    { key: 'warehouse_name', label: 'Gudang' },
    { key: 'reason', label: 'Alasan' },
  ];

  const getSortIcon = (key: string) => {
    if (sortKey !== key) return null;
    return sortOrder === 'asc' ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />;
  };

  const handleSort = (key: string) => {
    if (onSort) onSort(key);
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDateShort = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  };

  if (isLoading) {
    return (
      <div className="overflow-x-auto">
        <table className="w-full min-w-[700px] text-sm">
          <thead>
            <tr className="border-b border-gray-200 dark:border-gray-700">
              {headers.map((h) => (
                <th key={h.key} className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400">
                  {h.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {[...Array(5)].map((_, i) => (
              <tr key={i} className="border-b border-gray-100 dark:border-gray-800">
                <td className="px-4 py-3"><Skeleton className="h-4 w-28" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-20" /></td>
                <td className="px-4 py-3"><Skeleton className="h-6 w-16 rounded-full" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-12" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-24" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-32" /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  if (!transactions || transactions.length === 0) {
    return (
      <div className="py-12 text-center text-gray-500 dark:text-gray-400">
        <p className="text-lg">Belum ada transaksi</p>
        <p className="text-sm">Mulai catat transaksi untuk melihat riwayat.</p>
        {onExportCSV && (
          <Button variant="secondary" size="sm" onClick={onExportCSV} className="mt-4 gap-2">
            <Download className="h-4 w-4" />
            Export CSV
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Menampilkan {transactions.length} transaksi
        </p>
        {onExportCSV && (
          <Button variant="secondary" size="sm" onClick={onExportCSV} className="gap-2">
            <Download className="h-4 w-4" />
            Export CSV
          </Button>
        )}
      </div>

      <div className="overflow-x-auto">
        <table className="w-full min-w-[700px] text-sm">
          <thead>
            <tr className="border-b border-gray-200 dark:border-gray-700">
              {headers.map((h) => (
                <th
                  key={h.key}
                  className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 dark:hover:text-gray-200"
                  onClick={() => h.key !== 'actions' && handleSort(h.key)}
                >
                  <div className="flex items-center gap-1">
                    {h.label}
                    {getSortIcon(h.key)}
                  </div>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {transactions.map((txn) => (
              <React.Fragment key={txn.id}>
                <tr
                  className={cn(
                    'border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors cursor-pointer',
                    expandedRow === txn.id && 'bg-gray-50 dark:bg-gray-800/50'
                  )}
                  onClick={() => setExpandedRow(expandedRow === txn.id ? null : txn.id)}
                >
                  <td className="px-4 py-3 text-gray-600 dark:text-gray-300 whitespace-nowrap">
                    {formatDate(txn.created_at)}
                  </td>
                  <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                    {txn.sku_code}
                  </td>
                  <td className="px-4 py-3">
                    <Badge variant={txn.type === 'in' ? 'success' : 'danger'}>
                      {txn.type === 'in' ? 'Masuk' : 'Keluar'}
                    </Badge>
                  </td>
                  <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                    {txn.quantity}
                  </td>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                    {txn.warehouse_name || txn.warehouse_id || '-'}
                  </td>
                  <td className="px-4 py-3 text-gray-500 dark:text-gray-400 truncate max-w-[150px]">
                    {txn.reason || '-'}
                  </td>
                </tr>
                {expandedRow === txn.id && (
                  <tr className="bg-gray-50 dark:bg-gray-800/30">
                    <td colSpan={6} className="px-4 py-3">
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <span className="text-gray-500 dark:text-gray-400">ID Transaksi:</span>
                          <span className="ml-2 text-gray-700 dark:text-gray-300 font-mono text-xs">
                            {txn.id}
                          </span>
                        </div>
                        {txn.previous_quantity !== undefined && (
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Stok Sebelum:</span>
                            <span className="ml-2 text-gray-700 dark:text-gray-300">{txn.previous_quantity}</span>
                          </div>
                        )}
                        {txn.new_quantity !== undefined && (
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Stok Sesudah:</span>
                            <span className="ml-2 text-gray-700 dark:text-gray-300">{txn.new_quantity}</span>
                          </div>
                        )}
                        {txn.created_by && (
                          <div>
                            <span className="text-gray-500 dark:text-gray-400">Dibuat Oleh:</span>
                            <span className="ml-2 text-gray-700 dark:text-gray-300">{txn.created_by}</span>
                          </div>
                        )}
                        <div>
                          <span className="text-gray-500 dark:text-gray-400">Tanggal:</span>
                          <span className="ml-2 text-gray-700 dark:text-gray-300">
                            {formatDateShort(txn.created_at)}
                          </span>
                        </div>
                      </div>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
