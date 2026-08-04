import { renderHook } from '@testing-library/react';
import { ThemeProvider } from '@/context/ThemeContext';
import { useTheme } from '@/context/ThemeContext';
import { describe, it, expect } from '@jest/globals';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <ThemeProvider>{children}</ThemeProvider>
);

describe('useTheme', () => {
  it('should return theme state', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });
    expect(result.current.theme).toBeDefined();
    expect(result.current.isDark).toBeDefined();
  });

  it('should have setTheme function', () => {
    const { result } = renderHook(() => useTheme(), { wrapper });
    expect(typeof result.current.setTheme).toBe('function');
  });
});
