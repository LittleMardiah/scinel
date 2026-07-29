'use client';

import React from 'react';

export const Navbar: React.FC = () => {
  return (
    <nav className="sticky top-0 z-40 border-b border-gray-200 bg-white px-4 py-3 shadow-sm dark:border-gray-700 dark:bg-gray-900">
      <div className="flex items-center justify-between">
        <div className="text-xl font-bold text-blue-600 dark:text-blue-400">ScIntel</div>
        <div>Navbar (placeholder)</div>
      </div>
    </nav>
  );
};
