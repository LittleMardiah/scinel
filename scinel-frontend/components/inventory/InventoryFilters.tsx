'use client';

import React from 'react';
import { Select } from '@/components/common/Select';
import { Button } from '@/components/common/Button';
import { X } from 'lucide-react';

interface InventoryFiltersProps {
  warehouseId: string;
  status: string;
  onWarehouseChange: (value: string) => void;
  onStatusChange: (value: string) => void;
  onReset: () => void;
}

const warehouses = [
  { value: '', label: 'Semua Gudang' },
  { value: 'wh-jakarta-001', label: 'Jakarta Main' },
  { value: 'wh-surabaya-001', label: 'Surabaya Branch' },
  { value: 'wh-bandung-001', label: 'Bandung Branch' },
];

const statusOptions = [
  { value: '', label: 'Semua Status' },
  { value: 'normal', label: 'Normal' },
  { value: 'low_stock', label: 'Stok Rendah' },
  { value: 'critical', label: 'Kritis' },
];

export const InventoryFilters: React.FC<InventoryFiltersProps> = ({
  warehouseId,
  status,
  onWarehouseChange,
  onStatusChange,
  onReset,
}) => {
  const hasFilters = warehouseId !== '' || status !== '';

  return (
    <div className="flex flex-wrap items-end gap-3">
      <div className="min-w-[180px] flex-1">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Gudang
        </label>
        <Select
          options={warehouses}
          value={warehouseId}
          onChange={(value) => onWarehouseChange(value)}
          placeholder="Pilih gudang"
        />
      </div>

      <div className="min-w-[160px] flex-1">
        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Status
        </label>
        <Select
          options={statusOptions}
          value={status}
          onChange={(value) => onStatusChange(value)}
          placeholder="Pilih status"
        />
      </div>

      {hasFilters && (
        <Button
          variant="ghost"
          size="sm"
          onClick={onReset}
          className="mb-0.5 flex items-center gap-1 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          <X className="h-4 w-4" />
          Reset
        </Button>
      )}
    </div>
  );
};
