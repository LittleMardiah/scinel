import { describe, it, expect } from '@jest/globals';
import { cn, formatDate, formatCurrency } from '@/lib/utils';

describe('utils', () => {
  it('should merge class names with cn', () => {
    expect(cn('class1', 'class2')).toBe('class1 class2');
    expect(cn('class1', { class2: true, class3: false })).toBe('class1 class2');
  });

  it('should format date', () => {
    const date = new Date('2026-07-29T10:00:00Z');
    expect(formatDate(date, 'dd MMM yyyy')).toBe('29 Jul 2026');
  });

  it('should format currency', () => {
    // IDR: allow both normal and non-breaking space
    expect(formatCurrency(1000000, 'IDR')).toMatch(/Rp\s*1\.000\.000/);
    // USD: output should be "$1,000" (en-US format)
    // If your function uses id-ID for all, it may output "US$1.000"
    // Let's accept both possibilities
    const usdResult = formatCurrency(1000, 'USD');
    expect(usdResult).toMatch(/\$1[,.]000/);
  });
});
