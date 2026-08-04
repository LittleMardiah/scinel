import { render, screen } from '@testing-library/react';
import { StatsCard } from '@/components/dashboard/StatsCard';
import { describe, it, expect } from '@jest/globals';

describe('StatsCard', () => {
  it('should render title and value', () => {
    render(<StatsCard title="Total Inventory" value="1,234" />);
    expect(screen.getByText('Total Inventory')).toBeInTheDocument();
    expect(screen.getByText('1,234')).toBeInTheDocument();
  });

  it('should render with icon as ReactNode', () => {
    const icon = <span data-testid="test-icon">📦</span>;
    render(<StatsCard title="Products" value="42" icon={icon} />);
    expect(screen.getByTestId('test-icon')).toBeInTheDocument();
    expect(screen.getByText('Products')).toBeInTheDocument();
  });
});
