'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export const SystemStatus: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>SystemStatus</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-gray-500 dark:text-gray-400">Placeholder for SystemStatus</div>
      </CardContent>
    </Card>
  );
};
