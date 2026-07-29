'use client';

import React from 'react';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface AccuracyData {
  algorithm: string;
  mape: number;
  rmse: number;
  mae: number;
}

interface AccuracyChartProps {
  data: AccuracyData[];
  height?: number;
  metric?: 'mape' | 'rmse' | 'mae';
}

export const AccuracyChart: React.FC<AccuracyChartProps> = ({
  data,
  height = 250,
  metric = 'mape',
}) => {
  if (!data || data.length === 0) {
    return (
      <div className="flex h-48 items-center justify-center text-gray-500 dark:text-gray-400">
        Belum ada data accuracy
      </div>
    );
  }

  const metricLabels = {
    mape: 'MAPE (%)',
    rmse: 'RMSE',
    mae: 'MAE',
  };

  const metricColors = {
    mape: '#0ea5e9',
    rmse: '#f59e0b',
    mae: '#10b981',
  };

  const chartData = data.map((item) => ({
    name: item.algorithm === 'exponential_smoothing' ? 'ETS' : 'ARIMA',
    value: item[metric],
  }));

  return (
    <div style={{ width: '100%', height }}>
      <ResponsiveContainer>
        <BarChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-gray-200 dark:stroke-gray-700" />
          <XAxis dataKey="name" className="text-sm text-gray-500 dark:text-gray-400" />
          <YAxis className="text-sm text-gray-500 dark:text-gray-400" />
          <Tooltip
            contentStyle={{
              backgroundColor: 'rgba(255, 255, 255, 0.9)',
              borderColor: '#e5e7eb',
              borderRadius: '8px',
            }}
            formatter={(value: any) => {
              if (typeof value === 'number') return value.toFixed(2);
              return '-';
            }}
            labelStyle={{ color: '#1f2937' }}
          />
          <Bar
            dataKey="value"
            name={metricLabels[metric]}
            fill={metricColors[metric]}
            radius={[4, 4, 0, 0]}
          />
        </BarChart>
      </ResponsiveContainer>
      <div className="mt-2 text-center text-xs text-gray-400">
        <span className="font-medium text-gray-600 dark:text-gray-300">Interpretasi:</span>
        {' '}
        Lower {metricLabels[metric]} = better accuracy
      </div>
    </div>
  );
};
