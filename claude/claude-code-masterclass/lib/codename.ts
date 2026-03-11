const ADJECTIVES = [
  "swift",
  "silent",
  "phantom",
  "crimson",
  "shadow",
  "iron",
  "silver",
  "bold",
  "clever",
  "daring",
  "ghost",
  "hollow",
  "jade",
  "keen",
  "lunar",
  "marble",
  "neon",
  "obsidian",
  "polar",
  "rustic",
  "scarlet",
  "twisted",
  "velvet",
  "wicked",
  "zinc",
];

const NOUNS = [
  "fox",
  "raven",
  "wolf",
  "cipher",
  "viper",
  "falcon",
  "hawk",
  "pawn",
  "vault",
  "blade",
  "cobra",
  "dagger",
  "echo",
  "fuse",
  "glitch",
  "hound",
  "jackal",
  "knave",
  "lynx",
  "mole",
  "nexus",
  "oracle",
  "panther",
  "quill",
  "sphinx",
];

const VERBS = [
  "strikes",
  "lurks",
  "prowls",
  "vanishes",
  "hunts",
  "drifts",
  "glides",
  "hacks",
  "intercepts",
  "jams",
  "leaps",
  "masks",
  "navigates",
  "outfoxes",
  "pierces",
  "runs",
  "schemes",
  "stalks",
  "tracks",
  "unlocks",
  "vaults",
  "watches",
  "xrays",
  "yields",
  "zeroes",
];

function capitalize(word: string): string {
  return word.charAt(0).toUpperCase() + word.slice(1).toLowerCase();
}

function pick(arr: string[]): string {
  return arr[Math.floor(Math.random() * arr.length)];
}

export function generateCodename(): string {
  return (
    capitalize(pick(ADJECTIVES)) +
    capitalize(pick(NOUNS)) +
    capitalize(pick(VERBS))
  );
}
