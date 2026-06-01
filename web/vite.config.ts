import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:8080',
				changeOrigin: true
			},
			'/health': {
				target: 'http://localhost:8080',
				changeOrigin: true
			}
		}
	},
	build: {
		rollupOptions: {
			output: {
				manualChunks(id: string) {
					// Separate chunk for UI component library (shadcn-svelte)
					if (id.includes('node_modules/bits-ui') || id.includes('node_modules/clsx') || id.includes('node_modules/tailwind-merge')) {
						return 'ui-vendor';
					}
					// Separate chunk for data/table library
					if (id.includes('node_modules/@tanstack')) {
						return 'data-vendor';
					}
				}
			}
		}
	}
});
