'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { ThemeToggle } from '@/components/ui/ThemeToggle';
import { MobileNav } from './MobileNav';
import { TouchButton } from '@/components/common/TouchButton';
import { cn } from '@/lib/utils';

const navItems = [
  { name: 'Dashboard', href: '/dashboard', roles: ['admin', 'manager', 'analyst', 'viewer'] },
  { name: 'Inventory', href: '/inventory', roles: ['admin', 'manager', 'analyst', 'viewer'] },
  { name: 'Forecast', href: '/forecast', roles: ['admin', 'manager', 'analyst', 'viewer'] },
  { name: 'Anomalies', href: '/anomalies', roles: ['admin', 'manager', 'analyst'] },
  { name: 'Settings', href: '/settings', roles: ['admin', 'manager', 'analyst', 'viewer'] },
];

export function Navbar() {
  const { user } = useAuth();
  const pathname = usePathname();

  const filteredItems = navItems.filter((item) => {
    if (!item.roles) return true;
    if (!user) return false;
    return item.roles.includes(user.role as any);
  });

  return (
    <header className="sticky top-0 z-40 border-b border-gray-200 bg-white/95 backdrop-blur-sm dark:border-gray-800 dark:bg-gray-900/95">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Logo */}
        <Link href="/dashboard" className="text-xl font-bold text-blue-600 dark:text-blue-400">
          ScIntel
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden lg:flex items-center gap-1">
          {filteredItems.map((item) => {
            const isActive = pathname === item.href || pathname?.startsWith(item.href + '/');
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  'px-3 py-2 rounded-md text-sm font-medium transition-colors touch-manipulation min-h-[44px] min-w-[44px] flex items-center',
                  isActive
                    ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-800 active:bg-gray-200 dark:active:bg-gray-700'
                )}
              >
                {item.name}
              </Link>
            );
          })}
        </nav>

        {/* Right Side */}
        <div className="flex items-center gap-2">
          <ThemeToggle />
          <div className="hidden lg:block">
            <TouchButton
              variant="ghost"
              size="sm"
              onClick={() => window.location.href = '/profile'}
              className="rounded-full p-2 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              <span className="sr-only">Profil</span>
              <span className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center text-blue-700 dark:text-blue-300 font-medium">
                {user?.full_name?.charAt(0) || user?.email?.charAt(0) || 'U'}
              </span>
            </TouchButton>
          </div>
          {/* Mobile Hamburger */}
          <MobileNav items={filteredItems} />
        </div>
      </div>
    </header>
  );
}
