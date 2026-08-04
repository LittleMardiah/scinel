import { renderHook } from '@testing-library/react';
import { AuthProvider } from '@/context/AuthContext';
import { useAuth } from '@/context/AuthContext';
import { describe, it, expect } from '@jest/globals';

const wrapper = ({ children }: { children: React.ReactNode }) => (
  <AuthProvider>{children}</AuthProvider>
);

describe('useAuth', () => {
  it('should return auth state', () => {
    const { result } = renderHook(() => useAuth(), { wrapper });
    expect(result.current.isAuthenticated).toBeDefined();
  });

  it('should have login and logout functions', () => {
    const { result } = renderHook(() => useAuth(), { wrapper });
    expect(typeof result.current.login).toBe('function');
    expect(typeof result.current.logout).toBe('function');
  });
});
