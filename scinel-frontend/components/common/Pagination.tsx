'use client';

import React from 'react';
import { cn } from '@/lib/utils';
import { Button } from './Button';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
  className,
}) => {
  const getPageNumbers = () => {
    const pages = [];
    const maxVisible = 5;
    let start = Math.max(1, currentPage - Math.floor(maxVisible / 2));
    let end = Math.min(totalPages, start + maxVisible - 1);
    if (end - start + 1 < maxVisible) {
      start = Math.max(1, end - maxVisible + 1);
    }
    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    return pages;
  };

  if (totalPages <= 1) return null;

  return (
    <div className={cn('flex items-center gap-1', className)}>
      <Button
        variant="ghost"
        size="sm"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="h-8 px-3"
      >
        &lt;
      </Button>

      {getPageNumbers().map((page) => (
        <Button
          key={page}
          variant={page === currentPage ? 'primary' : 'ghost'}
          size="sm"
          onClick={() => onPageChange(page)}
          className={cn(
            'h-8 w-8 p-0',
            page === currentPage && 'pointer-events-none'
          )}
        >
          {page}
        </Button>
      ))}

      <Button
        variant="ghost"
        size="sm"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="h-8 px-3"
      >
        &gt;
      </Button>
    </div>
  );
};
