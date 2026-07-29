'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export const StatsCard: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>StatsCard</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-gray-500 dark:text-gray-400">Placeholder for StatsCard</div>
      </CardContent>
    </Card>
  );
};
