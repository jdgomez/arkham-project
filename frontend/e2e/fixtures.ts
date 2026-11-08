import { test as base, expect } from '@playwright/test';
import { spawn } from 'node:child_process';
import { createServer } from 'node:net';
import { once } from 'node:events';

// The demo game lives in the server's memory, so each test gets its own server
// (the binary produced by `make build`) and always starts from turn 1, Upkeep.
export const test = base.extend<{ baseURL: string }>({
  // Playwright requires the fixtures argument to be an object pattern, even when empty.
  // eslint-disable-next-line no-empty-pattern
  baseURL: async ({}, use) => {
    const port = await freePort();
    const server = spawn('../bin/arkham', ['-addr', `127.0.0.1:${port}`], { stdio: 'inherit' });
    const url = `http://127.0.0.1:${port}`;
    await waitUntilUp(url);
    await use(url);
    server.kill();
    await once(server, 'exit');
  },
});

export { expect };

async function freePort(): Promise<number> {
  const probe = createServer();
  probe.listen(0);
  await once(probe, 'listening');
  const address = probe.address();
  probe.close();
  if (address === null || typeof address === 'string') throw new Error('no free port');
  return address.port;
}

async function waitUntilUp(url: string): Promise<void> {
  for (let attempt = 0; attempt < 50; attempt++) {
    try {
      await fetch(url);
      return;
    } catch {
      await new Promise((resolve) => setTimeout(resolve, 100));
    }
  }
  throw new Error(`server at ${url} did not start`);
}
