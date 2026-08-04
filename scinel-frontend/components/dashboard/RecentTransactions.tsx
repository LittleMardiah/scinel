'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { ResponsiveTable } from '@/components/common/ResponsiveTable';
import { format } from 'date-fns';
import { id } from 'date-fns/locale';

interface Transaction {
  id: string;
  sku: string;
  type: 'in' | 'out';
  quantity: number;
  reason?: string;
  warehouse?: string;
  created_at: string;
}

interface RecentTransactionsProps {
  transactions: Transaction[];
}

export function RecentTransactions({ transactions }: RecentTransactionsProps) {
  if (!transactions || transactions.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Transaksi Terbaru</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-gray-500 dark:text-gray-400 text-center py-8">
            Belum ada transaksi terbaru
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Transaksi Terbaru</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveTable>
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
              {transactions.map((tx) => (
                <tr key={tx.id} className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50">
                  <td className="px-4 py-2 font-medium text-gray-900 dark:text-white">{tx.sku}</td>
                  <td className="px-4 py-2">
                    <span className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                      tx.type === 'in' 
                        ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                        : 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400'
                    }`}>
                      {tx.type === 'in' ? 'Masuk' : 'Keluar'}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-right">{tx.quantity}</td>
                  <td className="px-4 py-2 text-gray-600 dark:text-gray-300">{tx.reason || '-'}</td>
                  <td className="px-4 py-2 text-right text-gray-500 dark:text-gray-400">
                    {format(new Date(tx.created_at), 'd/M/yyyy', { locale: id })}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </ResponsiveTable>
      </CardContent>
    </Card>
  );
}
