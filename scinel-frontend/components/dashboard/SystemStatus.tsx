'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { CheckCircle, XCircle, Clock } from 'lucide-react';

interface SystemStatusProps {
  status: 'online' | 'offline' | 'loading';
  timestamp?: string;
}

export const SystemStatus: React.FC<SystemStatusProps> = ({ status, timestamp }) => {
  const statusConfig = {
    online: { icon: CheckCircle, color: 'text-green-600', label: 'Online' },
    offline: { icon: XCircle, color: 'text-red-600', label: 'Offline' },
    loading: { icon: Clock, color: 'text-yellow-600', label: 'Loading...' },
  };

  const config = statusConfig[status] || statusConfig.loading;
  const Icon = config.icon;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Status Sistem</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center gap-3">
          <Icon className={cn('h-6 w-6', config.color)} />
          <span className="font-medium text-gray-900 dark:text-white">{config.label}</span>
        </div>
        {timestamp && (
          <p className="mt-2 text-xs text-gray-500 dark:text-gray-400">
            Terakhir diperbarui: {new Date(timestamp).toLocaleTimeString('id-ID')}
          </p>
        )}
      </CardContent>
    </Card>
  );
};
