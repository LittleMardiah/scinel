'use client';
import { Button, ButtonProps } from './Button';
import { forwardRef } from 'react';

export const TouchButton = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className = '', size = 'md', ...props }, ref) => {
    // Force minimum touch target 44x44 on mobile, keep original sizes on desktop
    const touchClass = 'min-h-[44px] min-w-[44px] flex items-center justify-center';
    return <Button ref={ref} className={`${touchClass} ${className}`} size={size} {...props} />;
  }
);
TouchButton.displayName = 'TouchButton';
