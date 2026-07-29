'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Select } from '@/components/common/Select';
import { TextArea } from '@/components/ui/TextArea';
import { cn } from '@/lib/utils';
import apiClient from '@/lib/api';

interface TransactionFormData {
  warehouse_id: string;
  sku_id: string;
  type: 'in' | 'out';
  quantity: number;
  reason: string;
  details?: string;
}

interface TransactionFormProps {
  onSubmit: (data: TransactionFormData) => Promise<void>;
  isLoading?: boolean;
  error?: string | null;
  onCancel?: () => void;
}

export const TransactionForm: React.FC<TransactionFormProps> = ({
  onSubmit,
  isLoading,
  error,
  onCancel,
}) => {
  // Form state
  const [formData, setFormData] = useState<TransactionFormData>({
    warehouse_id: '',
    sku_id: '',
    type: 'in',
    quantity: 1,
    reason: '',
    details: '',
  });

  // SKU options (for autocomplete)
  const [skuOptions, setSkuOptions] = useState<{ value: string; label: string }[]>([]);
  const [skuSearch, setSkuSearch] = useState('');
  const [isLoadingSku, setIsLoadingSku] = useState(false);
  const [currentStock, setCurrentStock] = useState<number | null>(null);

  // Form errors
  const [formErrors, setFormErrors] = useState<Record<string, string>>({});

  // Fetch SKU options for autocomplete
  const fetchSkuOptions = useCallback(async (search: string) => {
    if (!search || search.length < 2) {
      setSkuOptions([]);
      return;
    }
    setIsLoadingSku(true);
    try {
      const response = await apiClient.get('/inventory', {
        params: { limit: 100, offset: 0, search: search },
      });
      const items = response.data.data?.inventory || [];
      const options = items.map((item: any) => ({
        value: item.sku_id || item.id,
        label: `${item.sku_code} - ${item.name || 'N/A'}`,
      }));
      setSkuOptions(options);
    } catch (err) {
      console.error('Failed to fetch SKU options:', err);
    } finally {
      setIsLoadingSku(false);
    }
  }, []);

  // Fetch current stock when SKU and warehouse change
  const fetchCurrentStock = useCallback(async () => {
    if (!formData.sku_id || !formData.warehouse_id) {
      setCurrentStock(null);
      return;
    }
    try {
      const response = await apiClient.get('/inventory', {
        params: {
          sku_id: formData.sku_id,
          warehouse_id: formData.warehouse_id,
        },
      });
      const items = response.data.data?.inventory || [];
      if (items.length > 0) {
        setCurrentStock(items[0].quantity || 0);
      } else {
        setCurrentStock(0);
      }
    } catch (err) {
      setCurrentStock(null);
    }
  }, [formData.sku_id, formData.warehouse_id]);

  // Debounce SKU search
  useEffect(() => {
    const timer = setTimeout(() => {
      if (skuSearch) fetchSkuOptions(skuSearch);
    }, 300);
    return () => clearTimeout(timer);
  }, [skuSearch, fetchSkuOptions]);

  // Fetch stock when SKU or warehouse changes
  useEffect(() => {
    fetchCurrentStock();
  }, [fetchCurrentStock]);

  // Validate form
  const validateForm = (): boolean => {
    const errors: Record<string, string> = {};

    if (!formData.warehouse_id) {
      errors.warehouse_id = 'Gudang harus dipilih';
    }
    if (!formData.sku_id) {
      errors.sku_id = 'SKU harus dipilih';
    }
    if (!formData.type) {
      errors.type = 'Tipe transaksi harus dipilih';
    }
    if (formData.quantity <= 0) {
      errors.quantity = 'Kuantitas harus lebih dari 0';
    }
    if (formData.type === 'out' && currentStock !== null && formData.quantity > currentStock) {
      errors.quantity = `Stok tidak mencukupi (tersedia: ${currentStock})`;
    }
    if (!formData.reason.trim()) {
      errors.reason = 'Alasan harus diisi';
    } else if (formData.reason.length > 255) {
      errors.reason = 'Alasan maksimal 255 karakter';
    }
    if (formData.details && formData.details.length > 1000) {
      errors.details = 'Detail maksimal 1000 karakter';
    }

    setFormErrors(errors);
    return Object.keys(errors).length === 0;
  };

  // Handle submit
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;
    onSubmit(formData);
  };

  // Handle field change
  const handleChange = (
    field: keyof TransactionFormData,
    value: string | number
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field
    if (formErrors[field]) {
      setFormErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  // Warehouse options
  const warehouseOptions = [
    { value: 'wh-jakarta-001', label: 'Jakarta Main' },
    { value: 'wh-surabaya-001', label: 'Surabaya Branch' },
    { value: 'wh-bandung-001', label: 'Bandung Branch' },
  ];

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Global error */}
      {error && (
        <div className="rounded-md bg-red-50 p-3 text-sm text-red-700 dark:bg-red-900/20 dark:text-red-400">
          {error}
        </div>
      )}

      {/* Warehouse */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Gudang <span className="text-red-500">*</span>
        </label>
        <Select
          options={warehouseOptions}
          value={formData.warehouse_id}
          onChange={(e) => handleChange('warehouse_id', e.target.value)}
          placeholder="Pilih gudang"
          error={formErrors.warehouse_id}
          required
        />
      </div>

      {/* SKU with autocomplete */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          SKU <span className="text-red-500">*</span>
        </label>
        <Input
          type="text"
          placeholder="Cari SKU... (min 2 karakter)"
          value={skuSearch}
          onChange={(e) => setSkuSearch(e.target.value)}
          disabled={isLoading}
          error={formErrors.sku_id}
        />
        {isLoadingSku && (
          <p className="mt-1 text-xs text-gray-400">Mencari SKU...</p>
        )}
        {skuOptions.length > 0 && (
          <div className="mt-1 max-h-48 overflow-y-auto rounded-md border border-gray-200 dark:border-gray-700">
            {skuOptions.map((opt) => (
              <button
                key={opt.value}
                type="button"
                className="w-full px-3 py-2 text-left text-sm hover:bg-gray-100 dark:hover:bg-gray-800"
                onClick={() => {
                  handleChange('sku_id', opt.value);
                  setSkuSearch(opt.label);
                  setSkuOptions([]);
                }}
              >
                {opt.label}
              </button>
            ))}
          </div>
        )}
        {currentStock !== null && formData.sku_id && (
          <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
            Stok saat ini: <span className="font-medium">{currentStock}</span> unit
          </p>
        )}
      </div>

      {/* Type (In/Out) */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Tipe Transaksi <span className="text-red-500">*</span>
        </label>
        <div className="flex gap-4">
          <label className="flex items-center gap-2 text-sm">
            <input
              type="radio"
              name="type"
              value="in"
              checked={formData.type === 'in'}
              onChange={(e) => handleChange('type', e.target.value as 'in')}
              className="h-4 w-4 text-blue-600"
            />
            Masuk
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input
              type="radio"
              name="type"
              value="out"
              checked={formData.type === 'out'}
              onChange={(e) => handleChange('type', e.target.value as 'out')}
              className="h-4 w-4 text-blue-600"
            />
            Keluar
          </label>
        </div>
        {formErrors.type && (
          <p className="mt-1 text-sm text-red-500">{formErrors.type}</p>
        )}
      </div>

      {/* Quantity */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Kuantitas <span className="text-red-500">*</span>
        </label>
        <Input
          type="number"
          min="1"
          step="1"
          value={formData.quantity}
          onChange={(e) => handleChange('quantity', parseInt(e.target.value) || 0)}
          error={formErrors.quantity}
          required
        />
      </div>

      {/* Reason */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Alasan <span className="text-red-500">*</span>
        </label>
        <Input
          type="text"
          placeholder="Contoh: Supplier delivery #PO-2026-001"
          value={formData.reason}
          onChange={(e) => handleChange('reason', e.target.value)}
          error={formErrors.reason}
          maxLength={255}
          required
        />
        <p className="mt-1 text-xs text-gray-400">
          {formData.reason.length}/255 karakter
        </p>
      </div>

      {/* Details (optional) */}
      <div>
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Detail Tambahan <span className="text-gray-400 text-xs">(opsional)</span>
        </label>
        <TextArea
          placeholder="Informasi tambahan tentang transaksi"
          value={formData.details}
          onChange={(e) => handleChange('details', e.target.value)}
          error={formErrors.details}
          rows={3}
          maxLength={1000}
        />
        <p className="mt-1 text-xs text-gray-400">
          {formData.details?.length || 0}/1000 karakter
        </p>
      </div>

      {/* Actions */}
      <div className="flex gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
        <Button
          type="submit"
          isLoading={isLoading}
          disabled={isLoading || Object.keys(formErrors).length > 0}
          className="flex-1"
        >
          {isLoading ? 'Memproses...' : 'Buat Transaksi'}
        </Button>
        {onCancel && (
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            disabled={isLoading}
          >
            Batal
          </Button>
        )}
      </div>
    </form>
  );
};
