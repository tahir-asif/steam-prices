import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,               // Use describe, it, expect without imports
    environment: 'jsdom',        // Simulate browser DOM
    setupFiles: './src/test/setup.ts', // Run before each test file
  },
});
