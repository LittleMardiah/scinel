'use client';

import React from 'react';
import { Input } from '@/components/common/Input';

interface DatePickerProps {
  value?: string;
  onChange?: (value: string) => void;
  label?: string;
  error?: string;
  disabled?: boolean;
}

export const DatePicker: React.FC<DatePickerProps> = ({ value, onChange, label, error, disabled }) => {
  return (
    <div>
      {label && <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">{label}</label>}
      <Input type="date" value={value} onChange={(e) => onChange?.(e.target.value)} error={error} disabled={disabled} />
    </div>
  );
};
