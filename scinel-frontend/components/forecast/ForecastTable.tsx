'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export interface ForecastRow {
  date?: string;
  predicted_demand?: number;
  lower_bound?: number;
  upper_bound?: number;
  confidence?: number;
}

interface ForecastTableProps {
  data?: ForecastRow[];
}

export const ForecastTable: React.FC<ForecastTableProps> = ({ data = [] }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>ForecastTable</CardTitle>
      </CardHeader>
      <CardContent>
        {data.length === 0 ? (
          <p className="py-6 text-center text-sm text-gray-500 dark:text-gray-400">
            Belum ada data forecast.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full min-w-[500px] text-sm">
              <thead>
                <tr className="border-b border-gray-200 dark:border-gray-700">
                  <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Tanggal</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Prediksi</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Batas Bawah</th>
                  <th className="px-4 py-2 text-left font-medium text-gray-500 dark:text-gray-400">Batas Atas</th>
                </tr>
              </thead>
              <tbody>
                {data.map((row, index) => (
                  <tr key={row.date ?? index} className="border-b border-gray-100 dark:border-gray-800">
                    <td className="px-4 py-2 text-gray-700 dark:text-gray-300">{row.date ?? '-'}</td>
                    <td className="px-4 py-2 text-gray-900 dark:text-white">{row.predicted_demand ?? '-'}</td>
                    <td className="px-4 py-2 text-gray-500 dark:text-gray-400">{row.lower_bound ?? '-'}</td>
                    <td className="px-4 py-2 text-gray-500 dark:text-gray-400">{row.upper_bound ?? '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
