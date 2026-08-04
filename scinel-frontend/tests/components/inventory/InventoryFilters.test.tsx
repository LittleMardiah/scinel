import { render, screen } from '@testing-library/react';
import { InventoryFilters } from '@/components/inventory/InventoryFilters';
import { describe, it, expect } from '@jest/globals';

describe('InventoryFilters', () => {
  it('should render filters', () => {
    render(<InventoryFilters onFilterChange={jest.fn()} />);
    expect(screen.getByLabelText(/gudang/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/status/i)).toBeInTheDocument();
  });
});
