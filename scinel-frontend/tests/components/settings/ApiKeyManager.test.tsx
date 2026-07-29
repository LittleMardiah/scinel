import { render, screen } from '@testing-library/react';
import { ApiKeyManager } from '@/components/settings/ApiKeyManager';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { describe, it, expect, vi } from 'vitest';
vi.mock('@/lib/api', () => ({ getApiKeys: vi.fn().mockResolvedValue([{ id: '1', key_prefix: 'sk_test', name: 'Test', permissions: ['read'], is_active: true, last_used_at: null, created_at: new Date().toISOString() }]), createApiKey: vi.fn(), revokeApiKey: vi.fn(), rotateApiKey: vi.fn() }));
describe('ApiKeyManager', () => {
  it('should render', async () => {
    render(<AuthProvider><ThemeProvider><ApiKeyManager /></ThemeProvider></AuthProvider>);
    expect(await screen.findByText('Test')).toBeInTheDocument();
  });
});
