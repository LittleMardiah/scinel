import { render, screen, fireEvent } from '@testing-library/react';
import { InventoryTable } from '@/components/inventory/InventoryTable';
import { describe, it, expect, jest } from '@jest/globals';

const mockItems = [
  {
    id: '1',
    sku_code: 'SKU-001',
    name: 'Product A',
    quantity: 100,
    status: 'normal' as const,
    warehouse_id: 'jakarta',
    warehouse_name: 'Jakarta',
    updated_at: new Date().toISOString(),
  },
  {
    id: '2',
    sku_code: 'SKU-002',
    name: 'Product B',
    quantity: 5,
    status: 'low_stock' as const,
    warehouse_id: 'bandung',
    warehouse_name: 'Bandung',
    updated_at: new Date().toISOString(),
  },
];

describe('InventoryTable', () => {
  it('should render inventory rows with data prop', () => {
    render(
      <InventoryTable
        data={mockItems}
        onSort={jest.fn()}
        sortKey="sku_code"
        sortOrder="asc"
      />
    );
    expect(screen.getByText('Product A')).toBeInTheDocument();
    expect(screen.getByText('Product B')).toBeInTheDocument();
  });

  it('should call onSort when header clicked', () => {
    const onSort = jest.fn();
    render(
      <InventoryTable
        data={mockItems}
        onSort={onSort}
        sortKey="sku_code"
        sortOrder="asc"
      />
    );
    const header = screen.getByText(/kode sku/i);
    fireEvent.click(header);
    expect(onSort).toHaveBeenCalledWith('sku_code');
  });

  it('should show loading skeleton', () => {
    const { container } = render(
      <InventoryTable
        data={[]}
        isLoading={true}
        onSort={jest.fn()}
        sortKey="sku_code"
        sortOrder="asc"
      />
    );
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('should show empty state', () => {
    render(
      <InventoryTable
        data={[]}
        onSort={jest.fn()}
        sortKey="sku_code"
        sortOrder="asc"
      />
    );
    expect(screen.getByText(/belum ada data/i)).toBeInTheDocument();
  });
});
