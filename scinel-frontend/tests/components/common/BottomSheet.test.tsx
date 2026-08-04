import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BottomSheet } from '@/components/common/BottomSheet';

describe('BottomSheet', () => {
  it('should render when open', () => {
    render(
      <BottomSheet isOpen={true} onClose={jest.fn()} title="Test Sheet">
        <div>Content</div>
      </BottomSheet>
    );
    expect(screen.getByText('Test Sheet')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should call onClose when backdrop clicked', async () => {
    const onClose = jest.fn();
    const user = userEvent.setup();
    render(
      <BottomSheet isOpen={true} onClose={onClose} closeOnBackdropClick={true}>
        <div>Content</div>
      </BottomSheet>
    );
    const backdrop = document.querySelector('div[aria-hidden="true"]');
    if (backdrop) await user.click(backdrop);
    expect(onClose).toHaveBeenCalled();
  });

  it('should call onClose when close button clicked', async () => {
    const onClose = jest.fn();
    const user = userEvent.setup();
    render(
      <BottomSheet isOpen={true} onClose={onClose} title="Test">
        <div>Content</div>
      </BottomSheet>
    );
    const closeBtn = screen.getByLabelText(/tutup/i);
    await user.click(closeBtn);
    expect(onClose).toHaveBeenCalled();
  });
});
