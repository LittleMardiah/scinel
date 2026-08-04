import { render, screen, fireEvent } from '@testing-library/react';
import { RegisterForm } from '@/components/auth/RegisterForm';
import { describe, it, expect, jest } from '@jest/globals';

describe('RegisterForm', () => {
  it('should render all form fields', () => {
    render(<RegisterForm onSubmit={jest.fn()} />);
    expect(screen.getByLabelText(/nama lengkap/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    // Ada 2 field dengan label "Kata Sandi" (password & confirm), gunakan getAll
    const passwordFields = screen.getAllByLabelText(/kata sandi/i);
    expect(passwordFields.length).toBe(2);
    expect(passwordFields[0]).toBeInTheDocument();
    expect(passwordFields[1]).toBeInTheDocument();
    expect(screen.getByLabelText(/konfirmasi/i)).toBeInTheDocument();
  });

  it('should validate required fields', async () => {
    const handleSubmit = jest.fn();
    render(<RegisterForm onSubmit={handleSubmit} />);
    fireEvent.click(screen.getByRole('button', { name: /daftar/i }));
    const errors = await screen.findAllByText(/diperlukan/i);
    expect(errors.length).toBeGreaterThan(0);
    expect(handleSubmit).not.toHaveBeenCalled();
  });
});
