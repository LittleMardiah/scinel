'use client';

import React from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Area,
  ComposedChart,
} from 'recharts';

interface ForecastData {
  date: string;
  predicted_demand: number;
  lower_bound?: number;
  upper_bound?: number;
  confidence?: number;
}

interface ForecastChartProps {
  data: ForecastData[];
  height?: number;
  isLoading?: boolean;
  algorithm?: string;
}

export const ForecastChart: React.FC<ForecastChartProps> = ({
  data,
  height = 300,
  isLoading,
  algorithm = 'ETS',
}) => {
  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center text-gray-500 dark:text-gray-400">
        Memuat data forecast...
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="flex h-64 items-center justify-center text-gray-500 dark:text-gray-400">
        Belum ada data forecast. Generate forecast terlebih dahulu.
      </div>
    );
  }

  const chartData = data.slice(0, 30).map((item) => ({
    date: new Date(item.date).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }),
    demand: item.predicted_demand,
    lower: item.lower_bound ?? item.predicted_demand * 0.8,
    upper: item.upper_bound ?? item.predicted_demand * 1.2,
  }));

  const formatYAxis = (value: number) => {
    if (value >= 1000) return `${value / 1000}K`;
    return value.toString();
  };

  return (
    <div style={{ width: '100%', height }}>
      <ResponsiveContainer>
        <ComposedChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
          <XAxis
            dataKey="date"
            className="text-xs text-gray-500 dark:text-gray-400"
            tick={{ fontSize: 11 }}
            interval="preserveStartEnd"
          />
          <YAxis
            className="text-xs text-gray-500 dark:text-gray-400"
            tickFormatter={formatYAxis}
            tick={{ fontSize: 11 }}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: 'rgba(255, 255, 255, 0.9)',
              borderColor: '#e5e7eb',
              borderRadius: '8px',
            }}
            labelStyle={{ color: '#1f2937' }}
          />
          <Legend wrapperStyle={{ fontSize: '12px' }} />

          <Area
            type="monotone"
            dataKey="upper"
            name="Upper Bound"
            stroke="none"
            fill="#93c5fd"
            fillOpacity={0.3}
          />
          <Area
            type="monotone"
            dataKey="lower"
            name="Lower Bound"
            stroke="none"
            fill="#93c5fd"
            fillOpacity={0.3}
          />

          <Line
            type="monotone"
            dataKey="demand"
            name="Prediksi Demand"
            stroke="#0ea5e9"
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 6 }}
          />
        </ComposedChart>
      </ResponsiveContainer>

      <div className="mt-2 text-center text-xs text-gray-400">
        Menggunakan algoritma: <span className="font-medium text-gray-600 dark:text-gray-300">{algorithm}</span>
        {' • '}
        Interval kepercayaan 95%
      </div>
    </div>
  );
};
