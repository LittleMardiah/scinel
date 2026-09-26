'use client';

import { useState, FormEvent } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Select } from '@/components/common/Select';

interface TransactionFormProps {
  onSubmit: (data: TransactionData) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
  onCancel?: () => void;
}

interface TransactionData {
  sku_code: string;
  type: 'in' | 'out';
  quantity: number;
  warehouse_id: string;
  reason: string; // ← required
}

export function TransactionForm({ onSubmit, isLoading, error, onCancel }: TransactionFormProps) {
  const [formData, setFormData] = useState<TransactionData>({
    sku_code: '',
    type: 'in',
    quantity: 0,
    warehouse_id: '',
    reason: '',
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
    if (!formData.reason) {
      errors.reason = 'Alasan diperlukan';
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
    { value: '62c35799-3209-4f64-85fc-4103557f5295', label: 'Jakarta' },
    { value: '4c8eb0e8-7552-4bf8-be26-d6471dd0d84c', label: 'Bandung' },
    { value: '95ed25d3-c64b-4307-af3c-ef2cb93dd459', label: 'Surabaya' },
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
        <label htmlFor="reason" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Alasan
        </label>
        <Input
          id="reason"
          type="text"
          placeholder="Masukkan alasan transaksi"
          value={formData.reason}
          onChange={(e) => setFormData({ ...formData, reason: e.target.value })}
          error={formErrors.reason}
        />
      </div>

      <div className="flex gap-2">
        {onCancel && (
          <Button type="button" variant="secondary" onClick={onCancel} className="flex-1">
            Batal
          </Button>
        )}
        <Button type="submit" isLoading={isLoading} className="flex-1">
          Buat Transaksi
        </Button>
      </div>
    </form>
  );
}