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
} from 'recharts';

interface ForecastData {
  date: string;
  predicted_demand: number;
  lower_bound?: number;
  upper_bound?: number;
}

interface ForecastChartProps {
  data: ForecastData[];
}

export const ForecastChart: React.FC<ForecastChartProps> = ({ data }) => {
  if (!data || data.length === 0) {
    return (
      <div className="flex h-64 items-center justify-center text-gray-500 dark:text-gray-400">
        Belum ada data forecast
      </div>
    );
  }

  const chartData = data.slice(0, 30).map((item) => ({
    date: new Date(item.date).toLocaleDateString('id-ID', { day: '2-digit', month: 'short' }),
    demand: item.predicted_demand,
    lower: item.lower_bound || item.predicted_demand * 0.8,
    upper: item.upper_bound || item.predicted_demand * 1.2,
  }));

  return (
    <ResponsiveContainer width="100%" height={250}>
      <LineChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
        <XAxis dataKey="date" className="text-xs text-gray-500 dark:text-gray-400" />
        <YAxis className="text-xs text-gray-500 dark:text-gray-400" />
        <Tooltip
          contentStyle={{
            backgroundColor: 'var(--tooltip-bg, #fff)',
            borderColor: 'var(--tooltip-border, #e5e7eb)',
            borderRadius: '8px',
          }}
        />
        <Legend />
        <Line
          type="monotone"
          dataKey="demand"
          name="Prediksi Demand"
          stroke="#0ea5e9"
          strokeWidth={2}
          dot={false}
        />
        <Line
          type="monotone"
          dataKey="upper"
          name="Upper Bound"
          stroke="#93c5fd"
          strokeDasharray="5 5"
          dot={false}
        />
        <Line
          type="monotone"
          dataKey="lower"
          name="Lower Bound"
          stroke="#93c5fd"
          strokeDasharray="5 5"
          dot={false}
        />
      </LineChart>
    </ResponsiveContainer>
  );
};
