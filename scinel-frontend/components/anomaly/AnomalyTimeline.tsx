'use client';

import React from 'react';
import { cn } from '@/lib/utils';
import { CheckCircle, AlertCircle, Clock, Check } from 'lucide-react';

interface TimelineEvent {
  timestamp: string;
  event: 'detected' | 'acknowledged' | 'resolved';
  notes?: string;
  user_id?: string;
}

interface AnomalyTimelineProps {
  events: TimelineEvent[];
  className?: string;
}

const eventConfig = {
  detected: {
    icon: AlertCircle,
    color: 'text-red-500 dark:text-red-400',
    bgColor: 'bg-red-100 dark:bg-red-900/30',
    label: 'Terdeteksi',
  },
  acknowledged: {
    icon: Clock,
    color: 'text-yellow-500 dark:text-yellow-400',
    bgColor: 'bg-yellow-100 dark:bg-yellow-900/30',
    label: 'Diteliti',
  },
  resolved: {
    icon: Check,
    color: 'text-green-500 dark:text-green-400',
    bgColor: 'bg-green-100 dark:bg-green-900/30',
    label: 'Selesai',
  },
};

export const AnomalyTimeline: React.FC<AnomalyTimelineProps> = ({ events, className }) => {
  if (!events || events.length === 0) {
    return (
      <div className="py-4 text-center text-gray-500 dark:text-gray-400">
        Belum ada timeline
      </div>
    );
  }

  return (
    <div className={cn('space-y-4', className)}>
      {events.map((event, index) => {
        const config = eventConfig[event.event];
        const Icon = config.icon;
        const isLast = index === events.length - 1;

        return (
          <div key={index} className="relative flex gap-4">
            {/* Timeline line */}
            {!isLast && (
              <div className="absolute left-3 top-8 h-full w-0.5 bg-gray-200 dark:bg-gray-700" />
            )}

            {/* Icon */}
            <div className={cn('z-10 flex h-6 w-6 items-center justify-center rounded-full', config.bgColor)}>
              <Icon className={cn('h-3.5 w-3.5', config.color)} />
            </div>

            {/* Content */}
            <div className="flex-1 pb-4">
              <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-3">
                <span className="font-medium text-gray-900 dark:text-white">
                  {config.label}
                </span>
                <span className="text-sm text-gray-500 dark:text-gray-400">
                  {new Date(event.timestamp).toLocaleString('id-ID', {
                    day: '2-digit',
                    month: 'short',
                    year: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })}
                </span>
              </div>
              {event.notes && (
                <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
                  {event.notes}
                </p>
              )}
              {event.user_id && (
                <p className="mt-1 text-xs text-gray-400">
                  Oleh: {event.user_id}
                </p>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
};
