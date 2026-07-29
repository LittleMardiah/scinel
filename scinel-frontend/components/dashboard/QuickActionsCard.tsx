'use client';

import React from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { Button } from '@/components/common/Button';
import { Plus, BarChart3, AlertTriangle, Settings } from 'lucide-react';

export const QuickActionsCard: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Aksi Cepat</CardTitle>
      </CardHeader>
      <CardContent className="grid grid-cols-2 gap-2">
        <Link href="/inventory/transaction/new" className="w-full">
          <Button variant="primary" size="sm" className="w-full justify-start gap-2">
            <Plus className="h-4 w-4" /> Transaksi
          </Button>
        </Link>
        <Link href="/forecast" className="w-full">
          <Button variant="secondary" size="sm" className="w-full justify-start gap-2">
            <BarChart3 className="h-4 w-4" /> Forecast
          </Button>
        </Link>
        <Link href="/anomalies" className="w-full">
          <Button variant="danger" size="sm" className="w-full justify-start gap-2">
            <AlertTriangle className="h-4 w-4" /> Anomali
          </Button>
        </Link>
        <Link href="/settings" className="w-full">
          <Button variant="ghost" size="sm" className="w-full justify-start gap-2">
            <Settings className="h-4 w-4" /> Pengaturan
          </Button>
        </Link>
      </CardContent>
    </Card>
  );
};
