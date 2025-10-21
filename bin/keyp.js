#!/usr/bin/env node

/**
 * keyp CLI entry point
 * Loads the compiled TypeScript CLI from lib/cli/index.js
 */

import('./lib/cli/index.js').catch((err) => {
  console.error('Failed to load CLI:', err);
  process.exit(1);
});
