'use client';

import React from 'react';

export const Footer: React.FC = () => {
  return (
    <footer className="border-t border-gray-200 bg-white px-4 py-3 text-center text-sm text-gray-500 dark:border-gray-700 dark:bg-gray-900 dark:text-gray-400">
      &copy; {new Date().getFullYear()} ScIntel. All rights reserved.
    </footer>
  );
};
