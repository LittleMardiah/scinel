import { render, screen } from '@testing-library/react';
import { Navbar } from '@/components/layout/Navbar';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { describe, it, expect } from '@jest/globals';

describe('Navbar', () => {
  it('should render navbar', () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <Navbar />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(screen.getByText('ScIntel')).toBeInTheDocument();
  });
});
