import { ResponsiveTable } from "@/components/common/ResponsiveTable";
'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { InventoryStatusBadge } from './InventoryStatusBadge';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { Badge } from '@/components/ui/Badge';
import { MoreVertical, Eye, Pencil, Trash2 } from 'lucide-react';
import { cn } from '@/lib/utils';

interface InventoryItem {
  id: string;
  sku_code: string;
  name?: string;
  quantity: number;
  status: 'normal' | 'low_stock' | 'critical';
  warehouse_id: string;
  warehouse_name?: string;
  last_restock_date?: string;
  updated_at: string;
}

interface InventoryTableProps {
  data: InventoryItem[];
  isLoading?: boolean;
  onSort: (key: string) => void;
  sortKey: string;
  sortOrder: 'asc' | 'desc';
  onRowClick?: (id: string) => void;
}

export const InventoryTable: React.FC<InventoryTableProps> = ({
  data,
  isLoading,
  onSort,
  sortKey,
  sortOrder,
  onRowClick,
}) => {
  const [actionMenu, setActionMenu] = useState<string | null>(null);

  const headers = [
    { key: 'sku_code', label: 'Kode SKU' },
    { key: 'name', label: 'Nama Produk' },
    { key: 'quantity', label: 'Kuantitas' },
    { key: 'status', label: 'Status' },
    { key: 'warehouse_name', label: 'Gudang' },
    { key: 'last_restock_date', label: 'Restock Terakhir' },
    { key: 'actions', label: 'Aksi' },
  ];

  const getSortIcon = (key: string) => {
    if (sortKey !== key) return '↕';
    return sortOrder === 'asc' ? '↑' : '↓';
  };

  const handleSort = (key: string) => {
    onSort(key);
  };

  if (isLoading) {
    return (
      <div className="overflow-x-auto">
        <ResponsiveTable><table className="w-full min-w-[700px] text-sm">
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
                <td className="px-4 py-3"><Skeleton className="h-4 w-20" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-32" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-12" /></td>
                <td className="px-4 py-3"><Skeleton className="h-6 w-16 rounded-full" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-24" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-20" /></td>
                <td className="px-4 py-3"><Skeleton className="h-8 w-8 rounded-full" /></td>
              </tr>
            ))}
          </tbody>
        </table></ResponsiveTable>
      </div>
    );
  }

  if (data.length === 0) {
    return (
      <div className="py-12 text-center text-gray-500 dark:text-gray-400">
        <p className="text-lg">Belum ada data inventaris</p>
        <p className="text-sm">Mulai tambahkan transaksi untuk melihat stok.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <ResponsiveTable><table className="w-full min-w-[700px] text-sm">
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
                  {h.key !== 'actions' && (
                    <span className="text-xs opacity-50">{getSortIcon(h.key)}</span>
                  )}
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((item) => (
            <tr
              key={item.id}
              className={cn(
                'border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors',
                onRowClick && 'cursor-pointer'
              )}
              onClick={() => onRowClick?.(item.id)}
            >
              <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                {item.sku_code}
              </td>
              <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                {item.name || '-'}
              </td>
              <td className="px-4 py-3 text-gray-900 dark:text-white">
                {item.quantity}
              </td>
              <td className="px-4 py-3">
                <InventoryStatusBadge status={item.status} />
              </td>
              <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                {item.warehouse_name || item.warehouse_id || '-'}
              </td>
              <td className="px-4 py-3 text-gray-500 dark:text-gray-400">
                {item.last_restock_date
                  ? new Date(item.last_restock_date).toLocaleDateString('id-ID')
                  : '-'}
              </td>
              <td className="px-4 py-3">
                <div className="relative">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation();
                      setActionMenu(actionMenu === item.id ? null : item.id);
                    }}
                    className="h-8 w-8 p-0"
                  >
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                  {actionMenu === item.id && (
                    <div className="absolute right-0 z-10 mt-1 min-w-[140px] rounded-md border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-700 dark:bg-gray-800">
                      <Link
                        href={`/inventory/${item.id}`}
                        className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                        onClick={() => setActionMenu(null)}
                      >
                        <Eye className="h-4 w-4" /> Detail
                      </Link>
                      <button
                        className="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                        onClick={(e) => {
                          e.stopPropagation();
                          // TODO: Edit action
                          setActionMenu(null);
                        }}
                      >
                        <Pencil className="h-4 w-4" /> Edit
                      </button>
                      <button
                        className="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:text-red-400 dark:hover:bg-gray-700"
                        onClick={(e) => {
                          e.stopPropagation();
                          // TODO: Delete action
                          setActionMenu(null);
                        }}
                      >
                        <Trash2 className="h-4 w-4" /> Hapus
                      </button>
                    </div>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table></ResponsiveTable>
    </div>
  );
};
