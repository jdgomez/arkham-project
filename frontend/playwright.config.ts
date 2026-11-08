import { defineConfig } from '@playwright/test';

const port = 8099;

// The E2E tests run against the binary produced by `make build`.
export default defineConfig({
  testDir: './e2e',
  use: { baseURL: `http://localhost:${port}` },
  webServer: {
    command: `../bin/arkham -addr :${port}`,
    url: `http://localhost:${port}`,
  },
});
