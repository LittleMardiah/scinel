import { render, screen, fireEvent } from '@testing-library/react';
import { Select } from '@/components/common/Select';
import { describe, it, expect, jest } from '@jest/globals';

describe('Select', () => {
  const options = [{ value: '1', label: 'Option 1' }, { value: '2', label: 'Option 2' }];

  it('should render options', () => {
    render(<Select options={options} onChange={jest.fn()} />);
    expect(screen.getByText('Option 1')).toBeInTheDocument();
    expect(screen.getByText('Option 2')).toBeInTheDocument();
  });

  it('should handle onChange with value string', () => {
    const handleChange = jest.fn();
    render(<Select options={options} onChange={handleChange} />);
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '2' } });
    expect(handleChange).toHaveBeenCalledWith('2');
  });
});
