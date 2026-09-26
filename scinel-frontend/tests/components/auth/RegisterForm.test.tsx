import { render, screen, fireEvent } from '@testing-library/react';
import { RegisterForm } from '@/components/auth/RegisterForm';
import { describe, it, expect, jest } from '@jest/globals';

describe('RegisterForm', () => {
  it('should render all form fields', () => {
    render(<RegisterForm onSubmit={jest.fn(async () => {})} />);
    expect(screen.getByPlaceholderText(/nama lengkap/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/nama@perusahaan\.com/i)).toBeInTheDocument();
    // Ada 2 field password (password & confirm), gunakan getAll
    const passwordFields = [
      screen.getByPlaceholderText(/minimal 8 karakter/i),
      screen.getByPlaceholderText(/ulangi kata sandi/i),
    ];
    expect(passwordFields.length).toBe(2);
    expect(passwordFields[0]).toBeInTheDocument();
    expect(passwordFields[1]).toBeInTheDocument();
    expect(screen.getByText(/konfirmasi kata sandi/i)).toBeInTheDocument();
  });

  it('should validate required fields', async () => {
    const handleSubmit = jest.fn(async () => {});
    render(<RegisterForm onSubmit={handleSubmit} />);
    fireEvent.click(screen.getByRole('button', { name: /daftar/i }));
    const errors = await screen.findAllByText(/diperlukan/i);
    expect(errors.length).toBeGreaterThan(0);
    expect(handleSubmit).not.toHaveBeenCalled();
  });
});
