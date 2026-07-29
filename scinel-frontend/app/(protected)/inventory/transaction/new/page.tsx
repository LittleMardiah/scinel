'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { TransactionForm } from '@/components/inventory/TransactionForm';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { ArrowLeft } from 'lucide-react';
import { Toast } from '@/components/common/Toast';
import apiClient from '@/lib/api';

interface TransactionFormData {
  warehouse_id: string;
  sku_id: string;
  type: 'in' | 'out';
  quantity: number;
  reason: string;
  details?: string;
}

export default function NewTransactionPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (data: TransactionFormData) => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await apiClient.post('/transactions', {
        warehouse_id: data.warehouse_id,
        sku_id: data.sku_id,
        type: data.type,
        quantity: data.quantity,
        reason: data.reason,
        details: data.details,
      });

      // Success
      Toast.success('Transaksi berhasil dibuat!');
      router.push('/inventory/transactions');
    } catch (err: any) {
      const status = err.response?.status;
      const message = err.response?.data?.message || err.message || 'Gagal membuat transaksi.';

      if (status === 409) {
        setError('Stok tidak mencukupi untuk transaksi keluar. Periksa stok yang tersedia.');
        Toast.error('Stok tidak mencukupi');
      } else if (status === 400) {
        const details = err.response?.data?.details || message;
        setError(`Validasi gagal: ${details}`);
        Toast.error('Validasi gagal');
      } else {
        setError(message);
        Toast.error(message);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    router.back();
  };

  return (
    <div className="max-w-2xl mx-auto py-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Link href="/inventory/transactions">
          <Button variant="ghost" size="sm" className="gap-2">
            <ArrowLeft className="h-4 w-4" />
            Kembali
          </Button>
        </Link>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Transaksi Baru
          </h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Catat pergerakan stok masuk atau keluar
          </p>
        </div>
      </div>

      {/* Form */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Detail Transaksi</CardTitle>
        </CardHeader>
        <CardContent>
          <TransactionForm
            onSubmit={handleSubmit}
            isLoading={isLoading}
            error={error}
            onCancel={handleCancel}
          />
        </CardContent>
      </Card>

      {/* Info */}
      <div className="rounded-md bg-blue-50 p-4 text-sm text-blue-700 dark:bg-blue-900/20 dark:text-blue-400">
        <p>
          <span className="font-medium">Tips:</span> Untuk transaksi <strong>Masuk</strong>,
          pastikan Anda memiliki data supplier atau PO. Untuk transaksi <strong>Keluar</strong>,
          pastikan stok mencukupi.
        </p>
      </div>
    </div>
  );
}
