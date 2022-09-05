import { writable } from 'svelte/store';

let defaultTheme = {
  name: "Default",
  font: "Fira Code, Menlo",
  fontSize: "16pt",
  textAreaColor: '#454158',
  backgroundColor: '#22212C',
  textColor: '#80ffea',
  borderColor: '#1B1A23',
  Cyan: "#80FFEA",
  Green: "#8AFF80",
  Orange: "#FFCA80",
  Pink: "#FF80BF",
  Purple: "#9580FF",
  Red: "#FF9580",
  Yellow: "#FFFF80",
};

export const theme = writable(defaultTheme);
