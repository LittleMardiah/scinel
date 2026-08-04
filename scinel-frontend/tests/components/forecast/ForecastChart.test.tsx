import { render, screen } from '@testing-library/react';
import { ForecastChart } from '@/components/forecast/ForecastChart';
import { describe, it, expect } from '@jest/globals';

const mockForecastData = [
  { date: '2024-01-01', predicted_demand: 100, lower_bound: 80, upper_bound: 120 },
  { date: '2024-01-02', predicted_demand: 110, lower_bound: 88, upper_bound: 132 },
  { date: '2024-01-03', predicted_demand: 105, lower_bound: 84, upper_bound: 126 },
];

describe('ForecastChart', () => {
  it('should render loading state', () => {
    render(<ForecastChart data={[]} isLoading={true} />);
    expect(screen.getByText(/memuat data forecast/i)).toBeInTheDocument();
  });

  it('should render empty state', () => {
    render(<ForecastChart data={[]} isLoading={false} />);
    expect(screen.getByText(/belum ada data forecast/i)).toBeInTheDocument();
  });

  it('should render chart with data', () => {
    const { container } = render(
      <ForecastChart data={mockForecastData} isLoading={false} />
    );
    // Mock recharts di setup render div
    const chartElement = container.querySelector('div');
    expect(chartElement).toBeInTheDocument();
  });

  it('should render algorithm info', () => {
    render(
      <ForecastChart data={mockForecastData} algorithm="ARIMA" isLoading={false} />
    );
    expect(screen.getByText(/ARIMA/i)).toBeInTheDocument();
  });
});
