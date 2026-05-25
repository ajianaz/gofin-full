import { browser } from '$app/environment';

type Theme = 'light' | 'dark' | 'system';

class ThemeStore {
	theme = $state<Theme>(
		browser ? (localStorage.getItem('gofin_theme') as Theme) ?? 'system' : 'system'
	);

	constructor() {
		if (browser && !localStorage.getItem('gofin_theme')) {
			localStorage.setItem('gofin_theme', 'system');
		}
	}

	get resolved(): 'light' | 'dark' {
		if (this.theme === 'system') {
			if (browser) {
				return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
			}
			return 'light';
		}
		return this.theme;
	}

	get isDark() {
		return this.resolved === 'dark';
	}

	setTheme = (t: Theme) => {
		this.theme = t;
		if (browser) {
			localStorage.setItem('gofin_theme', t);
			this.apply();
		}
	};

	toggle = () => {
		this.setTheme(this.resolved === 'dark' ? 'light' : 'dark');
	};

	apply = () => {
		if (!browser) return;
		const el = document.documentElement;
		if (this.resolved === 'dark') {
			el.classList.add('dark');
		} else {
			el.classList.remove('dark');
		}
	};
}

export const themeStore = new ThemeStore();

// Auto-apply on load
if (browser) {
	themeStore.apply();
	// Listen for system preference changes
	window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
		if (themeStore.theme === 'system') themeStore.apply();
	});
}
