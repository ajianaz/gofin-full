import { browser } from '$app/environment';
import id from '$lib/i18n/locales/id.json';
import en from '$lib/i18n/locales/en.json';

export type Locale = 'id' | 'en';

const localeCodes: Record<Locale, string> = {
	id: 'id-ID',
	en: 'en-US'
};

const messages: Record<Locale, Record<string, string>> = { id, en };

class I18nStore {
	locale = $state<Locale>(
		browser ? ((localStorage.getItem('gofin_locale') as Locale) ?? 'id') : 'id'
	);

	get localeCode() {
		return localeCodes[this.locale];
	}

	t = (key: string, params?: Record<string, string | number>): string => {
		let text = messages[this.locale][key] ?? messages.id[key] ?? key;
		if (params) {
			for (const [k, v] of Object.entries(params)) {
				text = text.replace(`{${k}}`, String(v));
			}
		}
		return text;
	};

	setLocale = (l: Locale) => {
		this.locale = l;
		if (browser) localStorage.setItem('gofin_locale', l);
	};
}

export const localeStore = new I18nStore();
