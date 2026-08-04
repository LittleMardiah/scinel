import { ReactNode } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/common/Card';
import { cn } from '@/lib/utils';

interface StatsCardProps {
  title: string;
  value: string | number;
  icon?: ReactNode;
  className?: string;
}

export function StatsCard({ title, value, icon, className }: StatsCardProps) {
  return (
    <Card className={cn('', className)}>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-gray-500 dark:text-gray-400">
          {title}
        </CardTitle>
        {icon && <div className="text-gray-400">{icon}</div>}
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold text-gray-900 dark:text-white">{value}</div>
      </CardContent>
    </Card>
  );
}

export default StatsCard;
