import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { LoginForm } from '@/components/auth/LoginForm';
import { describe, it, expect, jest } from '@jest/globals';

describe('LoginForm', () => {
  it('should render email and password fields', () => {
    render(<LoginForm onSubmit={jest.fn()} />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/kata sandi/i)).toBeInTheDocument();
  });

  it('should validate email format', async () => {
    const handleSubmit = jest.fn();
    render(<LoginForm onSubmit={handleSubmit} />);
    const emailInput = screen.getByLabelText(/email/i);
    const submitBtn = screen.getByRole('button', { name: /masuk/i });
    
    fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
    fireEvent.click(submitBtn);
    
    // Cari error message dengan role="alert"
    const errorMessage = await screen.findByRole('alert');
    expect(errorMessage).toHaveTextContent(/format email tidak valid/i);
    expect(handleSubmit).not.toHaveBeenCalled();
  });

  it('should submit with valid credentials', () => {
    const handleSubmit = jest.fn();
    render(<LoginForm onSubmit={handleSubmit} />);
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } });
    fireEvent.change(screen.getByLabelText(/kata sandi/i), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /masuk/i }));
    expect(handleSubmit).toHaveBeenCalledWith('test@example.com', 'password123');
  });

  it('should show error message', () => {
    render(<LoginForm onSubmit={jest.fn()} error="Login failed" />);
    expect(screen.getByText('Login failed')).toBeInTheDocument();
  });
});
