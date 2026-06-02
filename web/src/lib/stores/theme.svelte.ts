import { browser } from '$app/environment';

type Theme = 'light' | 'dark' | 'system';
type ThemePreset = 'blue' | 'violet' | 'emerald' | 'rose' | 'orange' | 'default';

const THEME_PRESETS: ThemePreset[] = ['blue', 'violet', 'emerald', 'rose', 'orange'];

class ThemeStore {
	theme = $state<Theme>(
		browser ? (localStorage.getItem('gofin_theme') as Theme) ?? 'system' : 'system'
	);

	preset = $state<ThemePreset>(
		browser ? ((localStorage.getItem('gofin_theme_preset') as ThemePreset) ?? 'default') : 'default'
	);

	get presets(): ThemePreset[] {
		return THEME_PRESETS;
	}

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

	setPreset = (p: ThemePreset) => {
		this.preset = p;
		if (browser) {
			localStorage.setItem('gofin_theme_preset', p);
			this.applyPreset();
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

	applyPreset = () => {
		if (!browser) return;
		const el = document.documentElement;
		if (this.preset && this.preset !== 'default') {
			el.setAttribute('data-theme-preset', this.preset);
		} else {
			el.removeAttribute('data-theme-preset');
		}
	};
}

export const themeStore = new ThemeStore();

// Auto-apply on load
if (browser) {
	themeStore.apply();
	themeStore.applyPreset();
	// Listen for system preference changes
	window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
		if (themeStore.theme === 'system') themeStore.apply();
	});
}
