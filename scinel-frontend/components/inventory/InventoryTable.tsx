'use client';

import React from 'react';
import { Badge } from '@/components/ui/Badge';
import { Skeleton } from '@/components/common/Skeleton';
import { ChevronUp, ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';

type InventoryStatus = 'normal' | 'low_stock' | 'critical';

export interface InventoryItem {
  id: string;
  sku_code: string;
  name: string;
  quantity: number;
  status: InventoryStatus | string;
  warehouse_id: string;
  warehouse_name?: string;
  updated_at: string;
}

interface InventoryTableProps {
  data: InventoryItem[];
  isLoading?: boolean;
  onSort?: (key: string) => void;
  sortKey?: string;
  sortOrder?: 'asc' | 'desc';
  onRowClick?: (id: string) => void;
}

const STATUS_META: Record<string, { label: string; variant: 'success' | 'warning' | 'danger' }> = {
  normal: { label: 'Normal', variant: 'success' },
  low_stock: { label: 'Stok Rendah', variant: 'warning' },
  critical: { label: 'Kritis', variant: 'danger' },
};

export const InventoryTable: React.FC<InventoryTableProps> = ({
  data,
  isLoading,
  onSort,
  sortKey = 'sku_code',
  sortOrder = 'asc',
  onRowClick,
}) => {
  const headers = [
    { key: 'sku_code', label: 'Kode SKU' },
    { key: 'name', label: 'Nama Barang' },
    { key: 'quantity', label: 'Stok' },
    { key: 'status', label: 'Status' },
    { key: 'warehouse_name', label: 'Gudang' },
    { key: 'updated_at', label: 'Diperbarui' },
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
    if (Number.isNaN(date.getTime())) return '-';
    return date.toLocaleDateString('id-ID', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  };

  const renderSkeletonRows = () => (
    <>
      {[...Array(5)].map((_, i) => (
        <tr key={i} className="border-b border-gray-100 dark:border-gray-800">
          <td className="px-4 py-3"><Skeleton className="h-4 w-24" /></td>
          <td className="px-4 py-3"><Skeleton className="h-4 w-40" /></td>
          <td className="px-4 py-3"><Skeleton className="h-4 w-12" /></td>
          <td className="px-4 py-3"><Skeleton className="h-6 w-20 rounded-full" /></td>
          <td className="px-4 py-3"><Skeleton className="h-4 w-28" /></td>
          <td className="px-4 py-3"><Skeleton className="h-4 w-24" /></td>
        </tr>
      ))}
    </>
  );

  const renderHeader = () => (
    <thead>
      <tr className="border-b border-gray-200 dark:border-gray-700">
        {headers.map((h) => (
          <th
            key={h.key}
            className="px-4 py-3 text-left font-medium text-gray-500 dark:text-gray-400 cursor-pointer hover:text-gray-700 dark:hover:text-gray-200"
            onClick={() => handleSort(h.key)}
          >
            <div className="flex items-center gap-1">
              {h.label}
              {getSortIcon(h.key)}
            </div>
          </th>
        ))}
      </tr>
    </thead>
  );

  if (isLoading) {
    return (
      <div className="overflow-x-auto">
        <table className="w-full min-w-[700px] text-sm">
          {renderHeader()}
          <tbody>{renderSkeletonRows()}</tbody>
        </table>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="py-12 text-center text-gray-500 dark:text-gray-400">
        <p className="text-lg">Belum ada data</p>
        <p className="text-sm">Belum ada barang di inventaris.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[700px] text-sm">
        {renderHeader()}
        <tbody>
          {data.map((item) => {
            const status = STATUS_META[item.status] ?? { label: item.status, variant: 'default' as const };
            return (
              <tr
                key={item.id}
                onClick={() => onRowClick && onRowClick(item.id)}
                className={cn(
                  'border-b border-gray-100 dark:border-gray-800 transition-colors',
                  onRowClick && 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800/50'
                )}
              >
                <td className="px-4 py-3 font-medium text-gray-900 dark:text-white whitespace-nowrap">
                  {item.sku_code}
                </td>
                <td className="px-4 py-3 text-gray-700 dark:text-gray-300">{item.name}</td>
                <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">{item.quantity}</td>
                <td className="px-4 py-3">
                  <Badge variant={status.variant}>{status.label}</Badge>
                </td>
                <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                  {item.warehouse_name || item.warehouse_id || '-'}
                </td>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400 whitespace-nowrap">
                  {formatDate(item.updated_at)}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};
