'use client';

import React from 'react';
import { cn } from '@/lib/utils';

interface RadioProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Radio: React.FC<RadioProps> = ({ label, error, className, ...props }) => {
  return (
    <div className="flex items-center gap-2">
      <input type="radio" className={cn('h-4 w-4 border-gray-300 text-blue-600 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800', className)} {...props} />
      {label && <label className="text-sm text-gray-700 dark:text-gray-300">{label}</label>}
      {error && <p className="text-sm text-red-500">{error}</p>}
    </div>
  );
};
