'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export interface AnomalyPoint {
  date?: string;
  value?: number;
  expected?: number;
}

interface AnomalyChartProps {
  data?: AnomalyPoint[];
}

export const AnomalyChart: React.FC<AnomalyChartProps> = ({ data = [] }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>AnomalyChart</CardTitle>
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <p className="py-6 text-center text-sm text-gray-500 dark:text-gray-400">
            Belum ada data anomali.
          </p>
        ) : (
          <ul className="space-y-2">
            {data.map((point, index) => (
              <li
                key={point.date ?? index}
                className="flex items-center justify-between text-sm text-gray-700 dark:text-gray-300"
              >
                <span>{point.date ?? '-'}</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {point.value ?? '-'}
                </span>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
};
