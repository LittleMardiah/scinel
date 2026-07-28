import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'
export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    include: ['**/*.test.ts', '**/*.test.tsx'],
    coverage: { provider: 'v8', reporter: ['text', 'html', 'lcov'], statements: 50, branches: 50, functions: 50, lines: 50 },
  },
  resolve: { alias: { '@': path.resolve(__dirname, './') } },
})
