import { render, screen } from '@testing-library/react';
import DashboardPage from '@/app/(protected)/dashboard/page';
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from '@/context/ThemeContext';
import { describe, it, expect } from '@jest/globals';

// Mock useDashboardData
jest.mock('@/hooks/useDashboardData', () => ({
  useDashboardData: () => ({
    stats: { totalSku: 100, totalWarehouse: 3, lowStock: 5, anomalies: 2 },
    isLoading: false,
  }),
}));

describe('DashboardPage', () => {
  it('should render dashboard', () => {
    render(
      <ThemeProvider>
        <AuthProvider>
          <DashboardPage />
        </AuthProvider>
      </ThemeProvider>
    );
    expect(screen.getByText(/Dashboard/i)).toBeInTheDocument();
  });
});
