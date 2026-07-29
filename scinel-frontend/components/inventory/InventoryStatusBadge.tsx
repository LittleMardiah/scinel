'use client';

import React from 'react';
import { Badge } from '@/components/ui/Badge';

interface InventoryStatusBadgeProps {
  status: 'normal' | 'low_stock' | 'critical' | string;
}

export const InventoryStatusBadge: React.FC<InventoryStatusBadgeProps> = ({ status }) => {
  const statusMap = {
    normal: { label: 'Normal', variant: 'success' as const },
    low_stock: { label: 'Stok Rendah', variant: 'warning' as const },
    critical: { label: 'Kritis', variant: 'danger' as const },
  };

  const defaultStatus = { label: status || 'Unknown', variant: 'default' as const };
  const config = statusMap[status as keyof typeof statusMap] || defaultStatus;

  return <Badge variant={config.variant}>{config.label}</Badge>;
};
