import { render, screen } from '@testing-library/react';
import { SeverityBadge } from '@/components/anomaly/SeverityBadge';
import { describe, it, expect } from '@jest/globals';

describe('SeverityBadge', () => {
  it('should render with severity labels', () => {
    render(<SeverityBadge severity="critical" />);
    expect(screen.getByText('Kritis')).toBeInTheDocument();
  });

  it('should apply correct color class for severity', () => {
    const { container } = render(<SeverityBadge severity="high" />);
    expect(container.firstChild).toHaveClass('bg-red-100');
  });
});
