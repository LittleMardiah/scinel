import { render, screen } from '@testing-library/react';
import { TouchButton } from '@/components/common/TouchButton';
describe('TouchButton', () => {
  it('should render with min touch classes', () => {
    render(<TouchButton>Tap me</TouchButton>);
    const btn = screen.getByText('Tap me');
    expect(btn).toHaveClass('min-h-[44px] min-w-[44px]');
  });
});
