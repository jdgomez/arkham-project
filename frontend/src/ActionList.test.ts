import { cleanup, render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, test, vi } from 'vitest';
import ActionList from './ActionList.svelte';

afterEach(cleanup);

test('renders one button per Available Action', () => {
  render(ActionList, { actions: [{ type: 'End phase' }, { type: 'Dance' }], onaction: vi.fn() });

  expect(screen.getAllByRole('button').map((button) => button.textContent)).toEqual([
    'End phase',
    'Dance',
  ]);
});

test('clicking an Action emits it', async () => {
  const onaction = vi.fn();
  render(ActionList, { actions: [{ type: 'End phase' }], onaction });

  await userEvent.click(screen.getByRole('button', { name: 'End phase' }));

  expect(onaction).toHaveBeenCalledExactlyOnceWith({ type: 'End phase' });
});
