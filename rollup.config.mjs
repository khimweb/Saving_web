// rollup.config.mjs
import svelte from 'rollup-plugin-svelte';
import resolve from '@rollup/plugin-node-resolve';
import commonjs from '@rollup/plugin-commonjs';
import css from 'rollup-plugin-css-only';
import terser from '@rollup/plugin-terser';

const plugins = (name) => [
  svelte({ compilerOptions: { dev: false } }),
  css({ output: `bundle-${name}.css` }),
  resolve({ browser: true, dedupe: ['svelte'] }),
  commonjs(),
  terser()
];

export default [
  {
    input: 'src/main.js',
    output: { sourcemap:false, format:'iife', name:'SavingApp', file:'public/build/bundle-saving.js' },
    plugins: plugins('saving')
  },
  {
    input: 'src/main-gold.js',
    output: { sourcemap:false, format:'iife', name:'GoldApp', file:'public/build/bundle-gold.js' },
    plugins: plugins('gold')
  },
  {
    input: 'src/main-analysis.js',
    output: { sourcemap:false, format:'iife', name:'AnalysisApp', file:'public/build/bundle-analysis.js' },
    plugins: plugins('analysis')
  }
];
