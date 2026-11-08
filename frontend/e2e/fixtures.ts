import { test as base, expect } from '@playwright/test';
import { spawn, type ChildProcess } from 'node:child_process';
import { once } from 'node:events';
import { fileURLToPath } from 'node:url';

const binary = fileURLToPath(new URL('../../bin/arkham', import.meta.url));

// The demo game lives in the server's memory, so each test gets its own server
// (the binary produced by `make build`) and always starts from turn 1, Upkeep.
export const test = base.extend<{ baseURL: string }>({
  // Playwright requires the fixtures argument to be an object pattern, even when empty.
  // eslint-disable-next-line no-empty-pattern
  baseURL: async ({}, use) => {
    const server = spawn(binary, ['-addr', '127.0.0.1:0'], {
      stdio: ['ignore', 'inherit', 'pipe'],
    });
    try {
      await use(await listeningURL(server));
    } finally {
      await stop(server);
    }
  },
});

export { expect };

// The server binds port 0 itself, so the port is never taken by someone else, and reports where it listens.
// stderr keeps flowing to ours afterwards, so the server never writes into a closed pipe.
function listeningURL(server: ChildProcess): Promise<string> {
  return new Promise((resolve, reject) => {
    const timer = setTimeout(
      () => reject(new Error('server did not start listening in 10s')),
      10_000,
    );
    let log = '';
    server.stderr!.setEncoding('utf8').on('data', (chunk: string) => {
      process.stderr.write(chunk);
      log += chunk;
      const address = /listening on (\S+)/.exec(log)?.[1];
      if (address) {
        clearTimeout(timer);
        resolve(`http://${address}`);
      }
    });
    server.once('exit', () => {
      clearTimeout(timer);
      reject(new Error(`server exited before listening: ${log}`));
    });
  });
}

async function stop(server: ChildProcess): Promise<void> {
  if (server.exitCode !== null || server.signalCode !== null) return;
  server.kill();
  await once(server, 'exit');
}
