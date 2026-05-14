import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { resolve } from 'path';

export default defineConfig({
	plugins: [svelte()],
	test: {
		globals: true,
		environment: 'jsdom',
		include: ['src/**/*.{test,spec}.{ts,js}'],
		passWithNoTests: true,
		alias: {
			$lib: resolve('./src/lib'),
			$components: resolve('./src/lib/components'),
			$app: resolve('./src')
		}
	}
});
