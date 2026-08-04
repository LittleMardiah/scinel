import { render, screen } from '@testing-library/react';
import { RecentTransactions } from '@/components/dashboard/RecentTransactions';
import { describe, it, expect } from '@jest/globals';

const mockData = [
  { id: '1', sku: 'SKU-001', type: 'in' as const, quantity: 10, warehouse: 'Jakarta', created_at: new Date().toISOString() },
  { id: '2', sku: 'SKU-002', type: 'out' as const, quantity: 5, warehouse: 'Bandung', created_at: new Date().toISOString() },
];

describe('RecentTransactions', () => {
  it('should render transaction list', () => {
    render(<RecentTransactions transactions={mockData} />);
    expect(screen.getByText('SKU-001')).toBeInTheDocument();
    expect(screen.getByText('SKU-002')).toBeInTheDocument();
  });

  it('should show empty state', () => {
    render(<RecentTransactions transactions={[]} />);
    expect(screen.getByText(/belum ada transaksi/i)).toBeInTheDocument();
  });
});
