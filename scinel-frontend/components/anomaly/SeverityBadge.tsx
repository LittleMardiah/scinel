'use client';

import React from 'react';
import { Badge } from '@/components/ui/Badge';
import { cn } from '@/lib/utils';

interface SeverityBadgeProps {
  severity: 'low' | 'medium' | 'high' | 'critical' | string;
  className?: string;
}

const severityMap = {
  low: { label: 'Rendah', variant: 'default' as const, color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400' },
  medium: { label: 'Sedang', variant: 'warning' as const, color: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400' },
  high: { label: 'Tinggi', variant: 'danger' as const, color: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400' },
  critical: { label: 'Kritis', variant: 'danger' as const, color: 'bg-red-700 text-white dark:bg-red-800 dark:text-white' },
};

export const SeverityBadge: React.FC<SeverityBadgeProps> = ({ severity, className }) => {
  const config = severityMap[severity as keyof typeof severityMap] || severityMap.low;

  return (
    <Badge variant={config.variant} className={cn(config.color, className)}>
      {config.label}
    </Badge>
  );
};
