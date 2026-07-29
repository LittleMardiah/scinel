'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { SeverityBadge } from './SeverityBadge';
import { Button } from '@/components/common/Button';
import { Skeleton } from '@/components/common/Skeleton';
import { Badge } from '@/components/ui/Badge';
import { MoreVertical, Eye, CheckCircle, XCircle, Trash2 } from 'lucide-react';
import { cn } from '@/lib/utils';

interface Anomaly {
  id: string;
  sku_code: string;
  type: string;
  severity: 'low' | 'medium' | 'high' | 'critical';
  description?: string;
  z_score: number;
  status: 'open' | 'acknowledged' | 'resolved';
  detected_at: string;
  resolved_at?: string;
}

interface AnomalyListProps {
  data: Anomaly[];
  isLoading?: boolean;
  onSort: (key: string) => void;
  sortKey: string;
  sortOrder: 'asc' | 'desc';
  onAcknowledge: (id: string) => void;
  onResolve: (id: string) => void;
  onDelete?: (id: string) => void;
  isAdmin?: boolean;
}

export const AnomalyList: React.FC<AnomalyListProps> = ({
  data,
  isLoading,
  onSort,
  sortKey,
  sortOrder,
  onAcknowledge,
  onResolve,
  onDelete,
  isAdmin,
}) => {
  const router = useRouter();
  const [actionMenu, setActionMenu] = useState<string | null>(null);

  const headers = [
    { key: 'sku_code', label: 'SKU' },
    { key: 'type', label: 'Tipe' },
    { key: 'severity', label: 'Severity' },
    { key: 'z_score', label: 'Z-Score' },
    { key: 'detected_at', label: 'Detected At' },
    { key: 'status', label: 'Status' },
    { key: 'actions', label: 'Aksi' },
  ];

  const getSortIcon = (key: string) => {
    if (sortKey !== key) return null;
    return sortOrder === 'asc' ? '↑' : '↓';
  };

  const handleSort = (key: string) => {
    if (key !== 'actions') onSort(key);
  };

  const handleRowClick = (id: string) => {
    router.push(`/anomalies/${id}`);
  };

  const statusColors = {
    open: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400',
    acknowledged: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
    resolved: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
  };

  const statusLabels = {
    open: 'Terbuka',
    acknowledged: 'Diteliti',
    resolved: 'Selesai',
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
                <td className="px-4 py-3"><Skeleton className="h-4 w-20" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-24" /></td>
                <td className="px-4 py-3"><Skeleton className="h-6 w-16 rounded-full" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-12" /></td>
                <td className="px-4 py-3"><Skeleton className="h-4 w-28" /></td>
                <td className="px-4 py-3"><Skeleton className="h-6 w-16 rounded-full" /></td>
                <td className="px-4 py-3"><Skeleton className="h-8 w-8 rounded-full" /></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="py-12 text-center text-gray-500 dark:text-gray-400">
        <p className="text-lg">Belum ada anomali</p>
        <p className="text-sm">Sistem akan mendeteksi anomali secara otomatis.</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full min-w-[700px] text-sm">
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
              className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors cursor-pointer"
              onClick={() => handleRowClick(item.id)}
            >
              <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                {item.sku_code}
              </td>
              <td className="px-4 py-3 text-gray-700 dark:text-gray-300">
                {item.type || '-'}
              </td>
              <td className="px-4 py-3">
                <SeverityBadge severity={item.severity} />
              </td>
              <td className="px-4 py-3 font-medium text-gray-900 dark:text-white">
                {item.z_score?.toFixed(2) || '-'}
              </td>
              <td className="px-4 py-3 text-gray-500 dark:text-gray-400">
                {new Date(item.detected_at).toLocaleString('id-ID', {
                  day: '2-digit',
                  month: 'short',
                  hour: '2-digit',
                  minute: '2-digit',
                })}
              </td>
              <td className="px-4 py-3">
                <Badge className={statusColors[item.status as keyof typeof statusColors]}>
                  {statusLabels[item.status as keyof typeof statusLabels] || item.status}
                </Badge>
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
                    <div className="absolute right-0 z-10 mt-1 min-w-[160px] rounded-md border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-700 dark:bg-gray-800">
                      <button
                        className="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                        onClick={(e) => {
                          e.stopPropagation();
                          router.push(`/anomalies/${item.id}`);
                          setActionMenu(null);
                        }}
                      >
                        <Eye className="h-4 w-4" /> Detail
                      </button>
                      {item.status === 'open' && (
                        <button
                          className="flex w-full items-center gap-2 px-4 py-2 text-sm text-yellow-600 hover:bg-gray-100 dark:text-yellow-400 dark:hover:bg-gray-700"
                          onClick={(e) => {
                            e.stopPropagation();
                            onAcknowledge(item.id);
                            setActionMenu(null);
                          }}
                        >
                          <CheckCircle className="h-4 w-4" /> Teliti
                        </button>
                      )}
                      {(item.status === 'open' || item.status === 'acknowledged') && (
                        <button
                          className="flex w-full items-center gap-2 px-4 py-2 text-sm text-green-600 hover:bg-gray-100 dark:text-green-400 dark:hover:bg-gray-700"
                          onClick={(e) => {
                            e.stopPropagation();
                            onResolve(item.id);
                            setActionMenu(null);
                          }}
                        >
                          <CheckCircle className="h-4 w-4" /> Selesai
                        </button>
                      )}
                      {isAdmin && (
                        <button
                          className="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:text-red-400 dark:hover:bg-gray-700"
                          onClick={(e) => {
                            e.stopPropagation();
                            onDelete?.(item.id);
                            setActionMenu(null);
                          }}
                        >
                          <Trash2 className="h-4 w-4" /> Hapus
                        </button>
                      )}
                    </div>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
