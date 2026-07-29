'use client';

import React from 'react';
import { Badge } from '@/components/ui/Badge';

interface Transaction {
  id: string;
  sku_code: string;
  type: 'in' | 'out';
  quantity: number;
  reason?: string;
  created_at: string;
}

interface RecentTransactionsProps {
  transactions: Transaction[];
}

export const RecentTransactions: React.FC<RecentTransactionsProps> = ({ transactions }) => {
  if (!transactions || transactions.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500 dark:text-gray-400">
        Belum ada transaksi terbaru
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 dark:border-gray-700">
            <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">SKU</th>
            <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Tipe</th>
            <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">Qty</th>
            <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Alasan</th>
            <th className="px-4 py-2 text-right font-medium text-gray-500 dark:text-gray-400">Waktu</th>
          </tr>
        </thead>
        <tbody>
          {transactions.slice(0, 5).map((txn) => (
            <tr key={txn.id} className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50">
              <td className="px-4 py-2 font-medium text-gray-900 dark:text-white">{txn.sku_code}</td>
              <td className="px-4 py-2">
                <Badge variant={txn.type === 'in' ? 'success' : 'danger'}>
                  {txn.type === 'in' ? 'Masuk' : 'Keluar'}
                </Badge>
              </td>
              <td className="px-4 py-2 text-right">{txn.quantity}</td>
              <td className="px-4 py-2 text-gray-600 dark:text-gray-300">{txn.reason || '-'}</td>
              <td className="px-4 py-2 text-right text-gray-500 dark:text-gray-400">
                {new Date(txn.created_at).toLocaleDateString('id-ID')}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
