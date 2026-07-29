'use client';
import { ReactNode, useRef, useState, useEffect } from 'react';
import { cn } from '@/lib/utils';

interface ResponsiveTableProps {
  children: ReactNode;
  className?: string;
  scrollIndicator?: boolean;
}

export function ResponsiveTable({ children, className = '', scrollIndicator = true }: ResponsiveTableProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [showLeftShadow, setShowLeftShadow] = useState(false);
  const [showRightShadow, setShowRightShadow] = useState(false);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const updateShadows = () => {
      const { scrollLeft, scrollWidth, clientWidth } = el;
      setShowLeftShadow(scrollLeft > 0);
      setShowRightShadow(scrollLeft + clientWidth < scrollWidth - 1);
    };

    updateShadows();
    el.addEventListener('scroll', updateShadows);
    window.addEventListener('resize', updateShadows);

    return () => {
      el.removeEventListener('scroll', updateShadows);
      window.removeEventListener('resize', updateShadows);
    };
  }, []);

  return (
    <div className="relative w-full">
      <div
        ref={containerRef}
        className={cn('overflow-x-auto overscroll-x-contain scroll-smooth', className)}
        style={{ WebkitOverflowScrolling: 'touch' }}
      >
        {children}
      </div>

      {scrollIndicator && (
        <>
          {/* Left fade indicator */}
          <div
            className={cn(
              'absolute top-0 left-0 h-full w-8 pointer-events-none transition-opacity duration-200',
              'bg-gradient-to-r from-white dark:from-gray-900 to-transparent',
              showLeftShadow ? 'opacity-100' : 'opacity-0'
            )}
          />
          {/* Right fade indicator */}
          <div
            className={cn(
              'absolute top-0 right-0 h-full w-8 pointer-events-none transition-opacity duration-200',
              'bg-gradient-to-l from-white dark:from-gray-900 to-transparent',
              showRightShadow ? 'opacity-100' : 'opacity-0'
            )}
          />
        </>
      )}
    </div>
  );
}
