import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Gofin',
  description: 'Self-hosted personal finance tracker — Go API + SvelteKit',
  lang: 'en',
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],
    ['meta', { name: 'theme-color', content: '#0f172a' }],
    ['meta', { name: 'og:type', content: 'website' }],
    ['meta', { name: 'og:title', content: 'Gofin — Self-Hosted Finance Tracker' }],
    ['meta', { name: 'og:description', content: 'Track wallets, transactions, budgets, and more. Self-hosted with Docker.' }],
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Docs', link: '/getting-started', activeMatch: '/getting-started|/features|/architecture|/configuration|/rbac|/security' },
      { text: 'Deployment', link: '/deployment' },
      { text: 'API', link: '/api/' },
      { text: 'Development', link: '/development' },
    ],
    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'Getting Started', link: '/getting-started' },
          { text: 'Features', link: '/features' },
        ],
      },
      {
        text: 'Design',
        items: [
          { text: 'Architecture', link: '/architecture' },
          { text: 'RBAC & Permissions', link: '/rbac' },
          { text: 'Security', link: '/security' },
        ],
      },
      {
        text: 'Operations',
        items: [
          { text: 'Configuration', link: '/configuration' },
          { text: 'Deployment', link: '/deployment' },
        ],
      },
      {
        text: 'API Reference',
        items: [
          { text: 'Overview', link: '/api/' },
        ],
      },
      {
        text: 'Development',
        items: [
          { text: 'Developer Setup', link: '/development' },
        ],
      },
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/ajianaz/gofin-full' },
    ],
    footer: {
      message: 'Released under the Apache-2.0 License.',
      copyright: '© 2026 Gofin Contributors',
    },
    search: {
      provider: 'local',
    },
  },
  cleanUrls: true,
  ignoreDeadLinks: true,
  srcExclude: ['**/research/**', '**/RUNBOOK.md'],
})
