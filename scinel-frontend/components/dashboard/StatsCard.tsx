'use client';

import React from 'react';
import { Card, CardContent } from '@/components/common/Card';
import { cn } from '@/lib/utils';
import { LucideIcon } from 'lucide-react';

interface StatsCardProps {
  title: string;
  value: number | string;
  icon: LucideIcon;
  trend?: { value: number; isPositive: boolean };
  variant?: 'default' | 'warning' | 'danger' | 'success';
}

export const StatsCard: React.FC<StatsCardProps> = ({
  title,
  value,
  icon: Icon,
  trend,
  variant = 'default',
}) => {
  const variantColors = {
    default: 'text-blue-600 bg-blue-100 dark:bg-blue-900/30 dark:text-blue-400',
    warning: 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900/30 dark:text-yellow-400',
    danger: 'text-red-600 bg-red-100 dark:bg-red-900/30 dark:text-red-400',
    success: 'text-green-600 bg-green-100 dark:bg-green-900/30 dark:text-green-400',
  };

  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-start justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500 dark:text-gray-400">
              {title}
            </p>
            <p className="mt-2 text-2xl font-bold text-gray-900 dark:text-white">
              {value}
            </p>
            {trend && (
              <p className={cn(
                'mt-1 text-xs font-medium',
                trend.isPositive ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'
              )}>
                {trend.isPositive ? '↑' : '↓'} {Math.abs(trend.value)}%
              </p>
            )}
          </div>
          <div className={cn('rounded-full p-2', variantColors[variant])}>
            <Icon className="h-5 w-5" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
