import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(date: Date | string, format: string = 'dd MMM yyyy'): string {
  const d = typeof date === 'string' ? new Date(date) : date;
  const day = String(d.getDate()).padStart(2, '0');
  const month = d.toLocaleString('id-ID', { month: 'short' });
  const year = d.getFullYear();
  
  return format
    .replace('dd', day)
    .replace('MMM', month)
    .replace('yyyy', String(year));
}

export function formatCurrency(amount: number, currency: string = 'IDR'): string {
  const formatter = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  });
  return formatter.format(amount);
}
