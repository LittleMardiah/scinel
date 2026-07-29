'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export const RecentTransactions: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>RecentTransactions</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-gray-500 dark:text-gray-400">Placeholder for RecentTransactions</div>
      </CardContent>
    </Card>
  );
};
