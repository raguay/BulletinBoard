import { writable } from 'svelte/store';

export const raw = writable({
  html: '',
  width: 0,
  height: 0,
  x: 0,
  y: 0
});

