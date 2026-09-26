'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

interface AccuracyMetricsProps {
  mape?: number;
  rmse?: number;
  mae?: number;
}

export const AccuracyMetrics: React.FC<AccuracyMetricsProps> = ({ mape, rmse, mae }) => {
  const metrics = [
    { label: 'MAPE', value: mape },
    { label: 'RMSE', value: rmse },
    { label: 'MAE', value: mae },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>AccuracyMetrics</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-3 gap-4">
          {metrics.map((metric) => (
            <div key={metric.label} className="text-center">
              <p className="text-sm text-gray-500 dark:text-gray-400">{metric.label}</p>
              <p className="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                {metric.value ?? '-'}
              </p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};
