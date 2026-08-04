import { render, screen } from '@testing-library/react';
import RootLayout from '@/app/layout';
import { describe, it, expect } from '@jest/globals';

describe('RootLayout', () => {
  it('should render children', () => {
    render(
      <RootLayout>
        <div>Test Child</div>
      </RootLayout>
    );
    expect(screen.getByText('Test Child')).toBeInTheDocument();
  });
});
