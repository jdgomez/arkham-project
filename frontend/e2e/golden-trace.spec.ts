import { test, expect } from './fixtures';

test('the player opens the app and sees the turn and the phase', async ({ page }) => {
  await page.goto('/');

  await expect(page.getByRole('heading', { name: 'Turn 1 - Upkeep' })).toBeVisible();
});
