import { render, screen } from '@testing-library/react';
import { ThemeSelector } from '@/components/settings/ThemeSelector';
import { ThemeProvider } from '@/context/ThemeContext';
import { describe, it, expect } from '@jest/globals';

describe('ThemeSelector', () => {
  it('should render theme options', () => {
    render(
      <ThemeProvider>
        <ThemeSelector />
      </ThemeProvider>
    );
    expect(screen.getByLabelText('☀️ Terang')).toBeInTheDocument();
    expect(screen.getByLabelText('🌙 Gelap')).toBeInTheDocument();
    expect(screen.getByLabelText('🖥️ Sistem')).toBeInTheDocument();
  });
});
