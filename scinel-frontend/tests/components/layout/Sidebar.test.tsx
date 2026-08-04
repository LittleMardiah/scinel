import { render, screen } from '@testing-library/react';
import { Sidebar } from '@/components/layout/Sidebar';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { describe, it, expect } from '@jest/globals';

describe('Sidebar', () => {
  it('should render sidebar', () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <Sidebar />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(screen.getByText(/Dashboard/i)).toBeInTheDocument();
  });
});
