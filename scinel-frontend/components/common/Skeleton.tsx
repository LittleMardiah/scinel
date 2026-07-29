'use client';

import React from 'react';
import { cn } from '@/lib/utils';

interface SkeletonProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'text' | 'circular' | 'rectangular';
}

export const Skeleton: React.FC<SkeletonProps> = ({ variant = 'text', className, ...props }) => {
  const variantClasses = {
    text: 'h-4 w-full',
    circular: 'h-12 w-12 rounded-full',
    rectangular: 'h-20 w-full rounded-md',
  };
  return <div className={cn('animate-pulse bg-gray-200 dark:bg-gray-700', variantClasses[variant], className)} {...props} />;
};
