'use client';

import { useState, FormEvent } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Select } from '@/components/common/Select';

interface TransactionFormProps {
  onSubmit: (data: TransactionData) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
}

interface TransactionData {
  sku_code: string;
  type: 'in' | 'out';
  quantity: number;
  warehouse_id: string;
  notes?: string;
}

export function TransactionForm({ onSubmit, isLoading, error }: TransactionFormProps) {
  const [formData, setFormData] = useState<TransactionData>({
    sku_code: '',
    type: 'in',
    quantity: 0,
    warehouse_id: '',
  });
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const errors: Record<string, string> = {};

    if (!formData.sku_code) {
      errors.sku_code = 'Kode SKU diperlukan';
    }

    if (formData.quantity <= 0) {
      errors.quantity = 'Kuantitas harus lebih dari 0';
    }

    if (!formData.warehouse_id) {
      errors.warehouse_id = 'Gudang diperlukan';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;
    await onSubmit(formData);
  };

  const warehouseOptions = [
    { value: 'jakarta', label: 'Jakarta' },
    { value: 'bandung', label: 'Bandung' },
    { value: 'surabaya', label: 'Surabaya' },
  ];

  const typeOptions = [
    { value: 'in', label: 'Masuk' },
    { value: 'out', label: 'Keluar' },
  ];

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-200 rounded-md text-sm border border-red-200 dark:border-red-800">
          {error}
        </div>
      )}

      <div>
        <label htmlFor="sku" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kode SKU
        </label>
        <Input
          id="sku"
          type="text"
          placeholder="e.g., SKU-001"
          value={formData.sku_code}
          onChange={(e) => setFormData({ ...formData, sku_code: e.target.value })}
          error={formErrors.sku_code}
        />
      </div>

      <div>
        <label htmlFor="type" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Tipe Transaksi
        </label>
        <Select
          id="type"
          options={typeOptions}
          value={formData.type}
          onChange={(value) => setFormData({ ...formData, type: value as 'in' | 'out' })}
        />
      </div>

      <div>
        <label htmlFor="quantity" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kuantitas
        </label>
        <Input
          id="quantity"
          type="number"
          placeholder="0"
          value={formData.quantity || ''}
          onChange={(e) => setFormData({ ...formData, quantity: parseInt(e.target.value) || 0 })}
          error={formErrors.quantity}
        />
      </div>

      <div>
        <label htmlFor="warehouse" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Gudang
        </label>
        <Select
          id="warehouse"
          options={warehouseOptions}
          value={formData.warehouse_id}
          onChange={(value) => setFormData({ ...formData, warehouse_id: value })}
          error={formErrors.warehouse_id}
        />
      </div>

      <div>
        <label htmlFor="notes" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Catatan (opsional)
        </label>
        <textarea
          id="notes"
          placeholder="Tambahkan catatan transaksi"
          value={formData.notes || ''}
          onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
          className="w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm transition-colors focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-gray-600 dark:bg-gray-800 dark:text-white dark:focus:border-blue-400"
          rows={3}
        />
      </div>

      <Button type="submit" isLoading={isLoading} className="w-full">
        Buat Transaksi
      </Button>
    </form>
  );
}
