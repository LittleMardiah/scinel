'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

interface AlgorithmComparisonProps {
  etsMape?: number;
  arimaMape?: number;
}

export const AlgorithmComparison: React.FC<AlgorithmComparisonProps> = ({ etsMape, arimaMape }) => {
  const algorithms = [
    { label: 'ETS', value: etsMape },
    { label: 'ARIMA', value: arimaMape },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>AlgorithmComparison</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4">
          {algorithms.map((algorithm) => (
            <div key={algorithm.label} className="text-center">
              <p className="text-sm text-gray-500 dark:text-gray-400">{algorithm.label}</p>
              <p className="mt-1 text-xl font-semibold text-gray-900 dark:text-white">
                {algorithm.value ?? '-'}
              </p>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};
