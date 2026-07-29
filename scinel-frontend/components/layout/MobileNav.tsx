'use client';
import { useState, useEffect } from 'react';
import { Menu, X, LogOut, User, Settings } from 'lucide-react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { TouchButton } from '@/components/common/TouchButton';
import { ThemeToggle } from '@/components/common/ThemeToggle';
import { cn } from '@/lib/utils';

interface NavItem {
  name: string;
  href: string;
  icon?: React.ReactNode;
  roles?: ('admin' | 'manager' | 'analyst' | 'viewer')[];
}

export function MobileNav({ items }: { items: NavItem[] }) {
  const [isOpen, setIsOpen] = useState(false);
  const pathname = usePathname();
  const router = useRouter();
  const { user, logout } = useAuth();

  // Close on route change
  useEffect(() => {
    setIsOpen(false);
  }, [pathname]);

  // Prevent body scroll when drawer open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  // Filter items based on user role
  const filteredItems = items.filter((item) => {
    if (!item.roles) return true;
    if (!user) return false;
    return item.roles.includes(user.role as any);
  });

  const handleLogout = async () => {
    await logout();
    router.push('/login');
    setIsOpen(false);
  };

  return (
    <>
      {/* Hamburger Button */}
      <TouchButton
        variant="ghost"
        size="sm"
        onClick={() => setIsOpen(true)}
        className="lg:hidden relative z-40 rounded-full p-2 hover:bg-gray-100 dark:hover:bg-gray-800"
        aria-label="Buka menu navigasi"
      >
        <Menu className="w-6 h-6" />
      </TouchButton>

      {/* Overlay */}
      <div
        className={cn(
          'fixed inset-0 z-50 bg-black/50 backdrop-blur-sm transition-all duration-300 lg:hidden',
          isOpen ? 'opacity-100 pointer-events-auto' : 'opacity-0 pointer-events-none'
        )}
        onClick={() => setIsOpen(false)}
        aria-hidden="true"
      />

      {/* Drawer */}
      <div
        className={cn(
          'fixed top-0 right-0 z-50 h-full w-[280px] max-w-[80vw] bg-white dark:bg-gray-900 shadow-2xl transition-transform duration-300 ease-out lg:hidden',
          'flex flex-col',
          isOpen ? 'translate-x-0' : 'translate-x-full'
        )}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-800">
          <span className="text-lg font-bold text-blue-600 dark:text-blue-400">ScIntel</span>
          <TouchButton
            variant="ghost"
            size="sm"
            onClick={() => setIsOpen(false)}
            className="rounded-full p-2 hover:bg-gray-100 dark:hover:bg-gray-800"
            aria-label="Tutup menu"
          >
            <X className="w-6 h-6" />
          </TouchButton>
        </div>

        {/* User Info */}
        {user && (
          <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-800">
            <p className="font-medium text-gray-900 dark:text-white truncate">{user.full_name || user.email}</p>
            <p className="text-sm text-gray-500 dark:text-gray-400 capitalize">{user.role}</p>
          </div>
        )}

        {/* Navigation Links */}
        <nav className="flex-1 overflow-y-auto p-4 space-y-1">
          {filteredItems.map((item) => {
            const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'flex items-center gap-3 px-3 py-3 rounded-lg text-base font-medium transition-colors touch-manipulation',
                  isActive
                    ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 active:bg-gray-200 dark:active:bg-gray-700'
                )}
                onClick={() => setIsOpen(false)}
              >
                {item.icon}
                {item.name}
              </Link>
            );
          })}
        </nav>

        {/* Footer Actions */}
        <div className="border-t border-gray-200 dark:border-gray-800 p-4 space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600 dark:text-gray-400">Tema</span>
            <ThemeToggle />
          </div>
          <TouchButton
            variant="ghost"
            onClick={handleLogout}
            className="w-full justify-start gap-3 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20"
          >
            <LogOut className="w-5 h-5" />
            Keluar
          </TouchButton>
        </div>
      </div>
    </>
  );
}
