import { defineConfig } from '@playwright/test';

// The E2E tests run against the binary produced by `make build`; e2e/fixtures.ts starts it.
export default defineConfig({
  testDir: './e2e',
});
