import { test, expect } from './fixtures';

test('the player opens the app and sees the map, Mandy Thompson and Azathoth', async ({ page }) => {
  await page.goto('/');

  const areas = {
    'Unvisited Isle': ['Merchant District Streets'],
    'River Docks': ['Merchant District Streets'],
    'The Unnamable': ['Merchant District Streets'],
    Administration: ['Miskatonic U. Streets'],
    Library: ['Miskatonic U. Streets'],
    'Science Building': ['Miskatonic U. Streets'],
    'Merchant District Streets': [
      'Unvisited Isle',
      'River Docks',
      'The Unnamable',
      'Miskatonic U. Streets',
    ],
    'Miskatonic U. Streets': [
      'Administration',
      'Library',
      'Science Building',
      'Merchant District Streets',
    ],
  };
  await expect(page.getByRole('heading', { level: 2 })).toHaveCount(8);
  for (const [name, connections] of Object.entries(areas)) {
    const area = page.getByRole('region', { name, exact: true });
    await expect(area.getByRole('heading', { name })).toBeVisible();
    await expect(area.getByLabel('Connections').getByRole('listitem')).toHaveText(connections);
  }

  await expect(
    page.getByRole('region', { name: 'Library', exact: true }).getByText('Mandy Thompson'),
  ).toBeVisible();
  await expect(page.getByText('Ancient One: Azathoth')).toBeVisible();
});
