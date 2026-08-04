import { render, screen, fireEvent } from '@testing-library/react';
import { Button } from '@/components/common/Button';
import { describe, it, expect, jest } from '@jest/globals';

describe('Button', () => {
  it('should render with text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('should handle click event', () => {
    const handleClick = jest.fn();
    render(<Button onClick={handleClick}>Click me</Button>);
    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('should be disabled when isLoading', () => {
    render(<Button isLoading>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('should apply variant styles', () => {
    const { container } = render(<Button variant="danger">Delete</Button>);
    expect(container.querySelector('button')).toHaveClass('bg-red-600');
  });
});
