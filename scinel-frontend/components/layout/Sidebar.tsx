'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { cn } from '@/lib/utils';
import {
  LayoutDashboard,
  Package,
  BarChart3,
  AlertTriangle,
  Settings,
  Webhook,
  Key,
  LogOut,
} from 'lucide-react';

interface NavItem {
  href: string;
  label: string;
  icon: React.ElementType;
  roles?: ('admin' | 'manager' | 'analyst' | 'viewer')[];
}

export const Sidebar: React.FC = () => {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const userRole = user?.role || 'viewer';

  const navItems: NavItem[] = [
    {
      href: '/dashboard',
      label: 'Dashboard',
      icon: LayoutDashboard,
      roles: ['admin', 'manager', 'analyst', 'viewer'],
    },
    {
      href: '/inventory',
      label: 'Inventory',
      icon: Package,
      roles: ['admin', 'manager', 'analyst', 'viewer'],
    },
    {
      href: '/forecast',
      label: 'Forecast',
      icon: BarChart3,
      roles: ['admin', 'manager', 'analyst', 'viewer'],
    },
    {
      href: '/anomalies',
      label: 'Anomalies',
      icon: AlertTriangle,
      roles: ['admin', 'manager', 'analyst', 'viewer'],
    },
    {
      href: '/webhooks',
      label: 'Webhooks',
      icon: Webhook,
      roles: ['admin', 'manager'],
    },
    {
      href: '/settings',
      label: 'Settings',
      icon: Settings,
      roles: ['admin', 'manager', 'analyst', 'viewer'],
    },
  ];

  const filteredItems = navItems.filter(
    (item) => !item.roles || item.roles.includes(userRole as any)
  );

  return (
    <aside className="hidden w-64 border-r border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-900 md:block">
      <nav className="space-y-1">
        {filteredItems.map((item) => {
          const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
          const Icon = item.icon;

          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                isActive
                  ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
                  : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800'
              )}
            >
              <Icon className="h-5 w-5" />
              {item.label}
            </Link>
          );
        })}

        <div className="border-t border-gray-200 pt-4 dark:border-gray-700">
          <button
            onClick={logout}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-red-600 transition-colors hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
          >
            <LogOut className="h-5 w-5" />
            Logout
          </button>
        </div>
      </nav>
    </aside>
  );
};
