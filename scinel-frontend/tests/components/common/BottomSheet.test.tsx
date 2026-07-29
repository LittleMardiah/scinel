import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BottomSheet } from '@/components/common/BottomSheet';
import { describe, it, expect, vi } from 'vitest';

describe('BottomSheet', () => {
  it('should render when open', () => {
    render(<BottomSheet isOpen={true} onClose={() => {}} title="Test Sheet"><div>Content</div></BottomSheet>);
    expect(screen.getByText('Test Sheet')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should call onClose when backdrop clicked', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<BottomSheet isOpen={true} onClose={onClose}><div>Content</div></BottomSheet>);
    const backdrop = document.querySelector('.fixed.inset-0.bg-black\\/50');
    if (backdrop) await user.click(backdrop);
    expect(onClose).toHaveBeenCalled();
  });
});
