import { render, screen } from '@testing-library/react';
import { AnomalyList } from '@/components/anomaly/AnomalyList';
import { describe, it, expect } from '@jest/globals';

const mockData = [
  { id: '1', sku: 'SKU-001', severity: 'high', description: 'Unusual drop', status: 'open', detected_at: '2026-07-29' },
  { id: '2', sku: 'SKU-002', severity: 'low', description: 'Minor spike', status: 'resolved', detected_at: '2026-07-28' },
];

describe('AnomalyList', () => {
  it('should render anomaly items', () => {
    render(
      <AnomalyList
        data={mockData}
        onSort={jest.fn()}
        sortKey="status"
        sortOrder="asc"
        onAcknowledge={jest.fn()}
        onResolve={jest.fn()}
      />
    );
    expect(screen.getByText(/SKU-001/i)).toBeInTheDocument();
    expect(screen.getByText(/SKU-002/i)).toBeInTheDocument();
    expect(screen.getByText(/open/i)).toBeInTheDocument();
    expect(screen.getByText(/resolved/i)).toBeInTheDocument();
  });
});
