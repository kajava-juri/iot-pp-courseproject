import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: [vitePreprocess()],
  kit: {
    adapter: adapter(),
    alias: {
      '$components': 'src/components',
      '$stores': 'src/stores',
      '$types': 'src/types',
      '$utils': 'src/utils'
    }
  },
  vite: {
    ssr: {
      noExternal: ['flowbite-svelte']
    }
  }
};

export default config;
