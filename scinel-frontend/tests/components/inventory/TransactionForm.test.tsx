import { render, screen, fireEvent } from '@testing-library/react';
import { TransactionForm } from '@/components/inventory/TransactionForm';
import { describe, it, expect, jest } from '@jest/globals';

describe('TransactionForm', () => {
  it('should render all form fields', () => {
    render(<TransactionForm onSubmit={jest.fn()} />);
    expect(screen.getByLabelText(/kode sku/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/tipe transaksi/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kuantitas/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/gudang/i)).toBeInTheDocument();
  });

  it('should validate required fields', async () => {
    const handleSubmit = jest.fn();
    render(<TransactionForm onSubmit={handleSubmit} />);
    fireEvent.click(screen.getByRole('button', { name: /buat transaksi/i }));
    const errors = await screen.findAllByText(/diperlukan/i);
    expect(errors.length).toBeGreaterThan(0);
    expect(handleSubmit).not.toHaveBeenCalled();
  });
});
