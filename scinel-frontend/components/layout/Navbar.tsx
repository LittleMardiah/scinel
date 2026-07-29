'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { useTheme } from '@/context/ThemeContext';
import { ThemeToggle } from '@/components/ui/ThemeToggle';
import { Button } from '@/components/common/Button';
import {
  Menu,
  X,
  User,
  LogOut,
  Settings,
  ChevronDown,
} from 'lucide-react';

export const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    router.push('/login');
  };

  return (
    <nav className="sticky top-0 z-50 border-b border-gray-200 bg-white px-4 py-3 shadow-sm dark:border-gray-700 dark:bg-gray-900">
      <div className="mx-auto flex max-w-7xl items-center justify-between">
        {/* Logo */}
        <Link href="/dashboard" className="flex items-center gap-2">
          <span className="text-2xl font-bold text-blue-600 dark:text-blue-400">
            ScIntel
          </span>
          <span className="hidden text-sm font-medium text-gray-500 dark:text-gray-400 md:inline">
            Supply Chain Intelligence
          </span>
        </Link>

        {/* Desktop Navigation */}
        <div className="hidden items-center gap-4 md:flex">
          <ThemeToggle />
          
          {/* User Menu */}
          <div className="relative">
            <button
              onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
              className="flex items-center gap-2 rounded-full px-3 py-1.5 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              <div className="flex h-8 w-8 items-center justify-center rounded-full bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400">
                <User className="h-4 w-4" />
              </div>
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                {user?.full_name || user?.email || 'User'}
              </span>
              <ChevronDown className="h-4 w-4 text-gray-500" />
            </button>

            {isUserMenuOpen && (
              <div className="absolute right-0 mt-2 w-48 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-gray-700 dark:bg-gray-800">
                <div className="border-b border-gray-200 px-4 py-2 dark:border-gray-700">
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    {user?.full_name || 'User'}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {user?.email}
                  </p>
                </div>
                <Link
                  href="/profile"
                  className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                  onClick={() => setIsUserMenuOpen(false)}
                >
                  <User className="h-4 w-4" /> Profile
                </Link>
                <Link
                  href="/settings"
                  className="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700"
                  onClick={() => setIsUserMenuOpen(false)}
                >
                  <Settings className="h-4 w-4" /> Settings
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-gray-100 dark:text-red-400 dark:hover:bg-gray-700"
                >
                  <LogOut className="h-4 w-4" /> Logout
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Mobile Menu Toggle */}
        <button
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          className="md:hidden"
        >
          {isMobileMenuOpen ? (
            <X className="h-6 w-6 text-gray-700 dark:text-gray-300" />
          ) : (
            <Menu className="h-6 w-6 text-gray-700 dark:text-gray-300" />
          )}
        </button>
      </div>

      {/* Mobile Menu */}
      {isMobileMenuOpen && (
        <div className="mt-3 border-t border-gray-200 pt-3 dark:border-gray-700 md:hidden">
          <div className="space-y-2">
            <Link
              href="/dashboard"
              className="block px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Dashboard
            </Link>
            <Link
              href="/inventory"
              className="block px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Inventory
            </Link>
            <Link
              href="/forecast"
              className="block px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Forecast
            </Link>
            <Link
              href="/anomalies"
              className="block px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-800"
              onClick={() => setIsMobileMenuOpen(false)}
            >
              Anomalies
            </Link>
            <div className="border-t border-gray-200 pt-2 dark:border-gray-700">
              <button
                onClick={handleLogout}
                className="flex w-full items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-gray-100 dark:text-red-400 dark:hover:bg-gray-800"
              >
                <LogOut className="h-4 w-4" /> Logout
              </button>
            </div>
          </div>
        </div>
      )}
    </nav>
  );
};
