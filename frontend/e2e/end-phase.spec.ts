import { test, expect } from './fixtures';

test('the player ends phases until the next turn starts', async ({ page }) => {
  await page.goto('/');

  await page.getByRole('button', { name: 'End phase' }).click();
  await expect(page.getByRole('heading', { name: 'Turn 1 - Movement' })).toBeVisible();

  for (let click = 0; click < 4; click++) {
    await page.getByRole('button', { name: 'End phase' }).click();
  }
  await expect(page.getByRole('heading', { name: 'Turn 2 - Upkeep' })).toBeVisible();
});
