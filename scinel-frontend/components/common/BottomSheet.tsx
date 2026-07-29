'use client';
import { ReactNode, useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import { TouchButton } from './TouchButton';
import { cn } from '@/lib/utils';

interface BottomSheetProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  className?: string;
  closeOnBackdropClick?: boolean;
}

export function BottomSheet({
  isOpen,
  onClose,
  title,
  children,
  className = '',
  closeOnBackdropClick = true,
}: BottomSheetProps) {
  const sheetRef = useRef<HTMLDivElement>(null);

  // Prevent body scroll when open
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

  // Close on Escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) onClose();
    };
    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return createPortal(
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/50 backdrop-blur-sm transition-all duration-300 animate-in fade-in">
      {/* Backdrop */}
      <div
        className="absolute inset-0"
        onClick={closeOnBackdropClick ? onClose : undefined}
        aria-hidden="true"
      />

      {/* Sheet Content */}
      <div
        ref={sheetRef}
        className={cn(
          'relative w-full max-w-md mx-auto bg-white dark:bg-gray-900 rounded-t-2xl shadow-2xl',
          'max-h-[85vh] overflow-y-auto animate-in slide-in-from-bottom duration-300',
          className
        )}
        style={{ touchAction: 'pan-y' }}
      >
        {/* Drag handle (visual) */}
        <div className="sticky top-0 z-10 flex items-center justify-between px-4 py-3 bg-white/95 dark:bg-gray-900/95 backdrop-blur-sm border-b border-gray-200 dark:border-gray-800">
          <div className="flex-1">
            {title && <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h3>}
          </div>
          <TouchButton
            variant="ghost"
            size="sm"
            onClick={onClose}
            className="rounded-full p-1 hover:bg-gray-100 dark:hover:bg-gray-800"
            aria-label="Tutup"
          >
            <X className="w-5 h-5" />
          </TouchButton>
        </div>

        <div className="p-4 pb-6">{children}</div>
      </div>
    </div>,
    document.body
  );
}
