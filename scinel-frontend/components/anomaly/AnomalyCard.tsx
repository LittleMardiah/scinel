'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export const AnomalyCard: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>AnomalyCard</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-gray-500 dark:text-gray-400">Placeholder for AnomalyCard</div>
      </CardContent>
    </Card>
  );
};
