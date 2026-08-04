import { render, screen, fireEvent } from '@testing-library/react';
import { Modal } from '@/components/common/Modal';
import { describe, it, expect, jest } from '@jest/globals';

describe('Modal', () => {
  it('should render when isOpen is true', () => {
    render(<Modal isOpen={true} onClose={jest.fn()} title="Test Modal"><div>Content</div></Modal>);
    expect(screen.getByText('Test Modal')).toBeInTheDocument();
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should not render when isOpen is false', () => {
    render(<Modal isOpen={false} onClose={jest.fn()}><div>Content</div></Modal>);
    expect(screen.queryByText('Content')).not.toBeInTheDocument();
  });

  it('should call onClose when close button clicked', () => {
    const onClose = jest.fn();
    render(<Modal isOpen={true} onClose={onClose}><div>Content</div></Modal>);
    const closeBtn = screen.getByRole('button', { name: /close/i });
    fireEvent.click(closeBtn);
    expect(onClose).toHaveBeenCalled();
  });
});
