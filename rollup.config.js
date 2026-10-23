// rollup.config.js
import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import css from 'rollup-plugin-css-only';
import { terser } from 'rollup-plugin-terser';

export default [
  // Saving page bundle
  {
    input: 'src/main.js',
    output: {
      sourcemap: true,
      format: 'iife',
      name: 'SavingApp',
      file: 'public/build/bundle-saving.js'
    },
    plugins: [
      svelte({
        compilerOptions: { dev: true }
      }),
      css({ output: 'bundle-saving.css' }),
      resolve({ browser: true, dedupe: ['svelte'] }),
      commonjs(),
      terser()
    ],
    watch: { clearScreen: false }
  },
  // Analysis page bundle
  {
    input: 'src/main-analysis.js',
    output: {
      sourcemap: true,
      format: 'iife',
      name: 'AnalysisApp',
      file: 'public/build/bundle-analysis.js'
    },
    plugins: [
      svelte({
        compilerOptions: { dev: true }
      }),
      css({ output: 'bundle-analysis.css' }),
      resolve({ browser: true, dedupe: ['svelte'] }),
      commonjs(),
      terser()
    ],
    watch: { clearScreen: false }
  }
];
