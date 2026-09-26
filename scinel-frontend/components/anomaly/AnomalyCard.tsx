'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Badge } from '@/components/ui/Badge';

interface AnomalyCardProps {
  title?: string;
  severity?: string;
}

export const AnomalyCard: React.FC<AnomalyCardProps> = ({ title, severity }) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>AnomalyCard</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="font-medium text-gray-900 dark:text-white">{title ?? '-'}</p>
        {severity && <Badge variant="warning">{severity}</Badge>}
      </CardContent>
    </Card>
  );
};
