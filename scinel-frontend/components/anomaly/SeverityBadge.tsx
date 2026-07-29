'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';

export const SeverityBadge: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>SeverityBadge</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-sm text-gray-500 dark:text-gray-400">Placeholder for SeverityBadge</div>
      </CardContent>
    </Card>
  );
};
