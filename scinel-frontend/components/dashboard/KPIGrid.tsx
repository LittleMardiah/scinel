'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

interface KPIGridProps {
  totalSku?: number;
  totalWarehouse?: number;
  lowStockCount?: number;
  anomalyCount?: number;
}

export const KPIGrid: React.FC<KPIGridProps> = ({
  totalSku = 0,
  totalWarehouse = 0,
  lowStockCount = 0,
  anomalyCount = 0,
}) => {
  const items = [
    { label: 'Total SKU', value: totalSku },
    { label: 'Total Gudang', value: totalWarehouse },
    { label: 'Stok Rendah', value: lowStockCount },
    { label: 'Anomali', value: anomalyCount },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>KPIGrid</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
          {items.map((item) => (
            <div key={item.label} className="rounded-lg border border-gray-200 p-4 dark:border-gray-700">
              <p className="text-sm text-gray-500 dark:text-gray-400">{item.label}</p>
              <p className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{item.value}</p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};
