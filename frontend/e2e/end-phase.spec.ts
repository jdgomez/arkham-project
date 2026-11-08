import { test, expect } from './fixtures';

test('the player ends phases until the next turn starts', async ({ page }) => {
  await page.goto('/');

  const phases = ['Movement', 'Arkham Encounters', 'Other World Encounters', 'Mythos'];
  const headings = [...phases.map((phase) => `Turn 1 - ${phase}`), 'Turn 2 - Upkeep'];
  for (const heading of headings) {
    await page.getByRole('button', { name: 'End phase' }).click();
    await expect(page.getByRole('heading', { name: heading })).toBeVisible();
  }
});
