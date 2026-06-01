<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { authStore } from '$lib/stores/auth.svelte.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { themeStore } from '$lib/stores/theme.svelte.js';
	import { LanguageSwitcher } from '$lib/components/shared/index.js';

	const titleMap: Record<string, string> = {
		'/dashboard': 'layout.sidebar.dashboard',
		'/transactions': 'layout.sidebar.transactions',
		'/wallets': 'layout.sidebar.wallets',
		'/budgets': 'layout.sidebar.budgets',
		'/piggy-banks': 'layout.sidebar.piggyBanks',
		'/bills': 'layout.sidebar.bills',
		'/recurring': 'layout.sidebar.recurring',
		'/rules': 'layout.sidebar.rules',
		'/reports': 'layout.sidebar.reports',
		'/settings': 'layout.sidebar.settings',
		'/currencies': 'layout.sidebar.currencies',
		'/categories': 'layout.sidebar.categories',
		'/tags': 'layout.sidebar.tags'
	};

	const pageTitle = $derived(() => {
		const path = $page.url.pathname;
		if (titleMap[path]) return t(titleMap[path]);
		const base = '/' + path.split('/').filter(Boolean)[0];
		if (titleMap[base]) return t(titleMap[base]);
		return t('app.name');
	});
	import {
		Sidebar,
		SidebarContent,
		SidebarFooter,
		SidebarGroup,
		SidebarGroupContent,
		SidebarGroupLabel,
		SidebarHeader,
		SidebarInset,
		SidebarMenu,
		SidebarMenuButton,
		SidebarMenuItem,
		SidebarProvider,
		SidebarRail,
		SidebarTrigger,
		Separator
	} from '$lib/components/ui/sidebar/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Avatar, AvatarFallback } from '$lib/components/ui/avatar/index.js';
	import {
		LayoutDashboard,
		ArrowLeftRight,
		Wallet,
		PiggyBank,
		Receipt,
		Repeat,
		Shield,
		FileText,
		Settings,
		LogOut,
		ChevronsUpDown,
		Hexagon,
		Sun,
		Moon
	} from '@lucide/svelte';

	let { children } = $props();
	const t = localeStore.t;
	let mounted = $state(false);

	const displayName = $derived(
		authStore.user?.name || authStore.user?.email?.split('@')[0] || ''
	);
	const userEmail = $derived(authStore.user?.email || '');

	// Reactively update sidebar footer DOM elements (shadcn SidebarFooter breaks Svelte reactivity)
	$effect(() => {
		if (!mounted) return;
		const nameEl = document.getElementById('user-name');
		const emailEl = document.getElementById('user-email');
		if (nameEl) nameEl.textContent = displayName;
		if (emailEl) emailEl.textContent = userEmail;
	});

	const menuNav = $derived([
		{ label: t('layout.sidebar.dashboard'), href: '/dashboard', icon: LayoutDashboard },
		{ label: t('layout.sidebar.transactions'), href: '/transactions', icon: ArrowLeftRight },
		{ label: t('layout.sidebar.wallets'), href: '/wallets', icon: Wallet },
		{ label: t('layout.sidebar.budgets'), href: '/budgets', icon: PiggyBank }
	]);

	const financeNav = $derived([
		{ label: t('layout.sidebar.piggyBanks'), href: '/piggy-banks', icon: PiggyBank },
		{ label: t('layout.sidebar.bills'), href: '/bills', icon: Receipt },
		{ label: t('layout.sidebar.recurring'), href: '/recurring', icon: Repeat }
	]);

	const otherNav = $derived([
		{ label: t('layout.sidebar.rules'), href: '/rules', icon: Shield },
		{ label: t('layout.sidebar.reports'), href: '/reports', icon: FileText },
		{ label: t('layout.sidebar.settings'), href: '/settings', icon: Settings }
	]);

	const userInitials = $derived(
		authStore.user?.name
			?.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2) ?? 'U'
	);

	onMount(async () => {
		mounted = true;
		if (authStore.accessToken) {
			await authStore.restore();
		}
		if (!authStore.isAuthenticated || !authStore.user) {
			goto('/login');
		}
	});

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}

	async function handleLogout() {
		await authStore.logout();
		if (typeof localStorage !== 'undefined') {
			localStorage.removeItem('access_token');
			localStorage.removeItem('refresh_token');
		}
		goto('/login');
	}
</script>

<svelte:head>
	<title>{pageTitle()} · Gofin</title>
</svelte:head>

<SidebarProvider>
	<Sidebar>
		<SidebarHeader>
			<div class="flex items-center justify-between rounded-md p-2">
				<div class="flex items-center gap-2">
					<div class="flex size-8 items-center justify-center rounded-[10px] bg-sidebar-accent">
						<Hexagon class="size-4 text-sidebar-accent-foreground" />
					</div>
					<div class="flex flex-col justify-center">
						<span class="text-sm font-medium text-sidebar-foreground">{t('layout.sidebar.brand')}</span>
						<span class="text-xs text-sidebar-foreground">{t('layout.sidebar.subtitle')}</span>
					</div>
				</div>
				<ChevronsUpDown class="size-4 text-muted-foreground" />
			</div>
		</SidebarHeader>

		<SidebarContent>
			<SidebarGroup>
				<SidebarGroupLabel>{t('layout.sidebar.menuGroup')}</SidebarGroupLabel>
				<SidebarGroupContent>
					<SidebarMenu>
						{#each menuNav as item}
							{@const Icon = item.icon}
							<SidebarMenuItem>
								<SidebarMenuButton
									isActive={isActive(item.href)}
									tooltipContent={item.label}
								>
									{#snippet child({ props })}
										<a href={item.href} {...props}>
											<Icon />
											<span>{item.label}</span>
										</a>
									{/snippet}
								</SidebarMenuButton>
							</SidebarMenuItem>
						{/each}
					</SidebarMenu>
				</SidebarGroupContent>
			</SidebarGroup>

			<SidebarGroup>
				<SidebarGroupLabel>{t('layout.sidebar.financeGroup')}</SidebarGroupLabel>
				<SidebarGroupContent>
					<SidebarMenu>
						{#each financeNav as item}
							{@const Icon = item.icon}
							<SidebarMenuItem>
								<SidebarMenuButton
									isActive={isActive(item.href)}
									tooltipContent={item.label}
								>
									{#snippet child({ props })}
										<a href={item.href} {...props}>
											<Icon />
											<span>{item.label}</span>
										</a>
									{/snippet}
								</SidebarMenuButton>
							</SidebarMenuItem>
						{/each}
					</SidebarMenu>
				</SidebarGroupContent>
			</SidebarGroup>

			<SidebarGroup>
				<SidebarGroupLabel>{t('layout.sidebar.otherGroup')}</SidebarGroupLabel>
				<SidebarGroupContent>
					<SidebarMenu>
						{#each otherNav as item}
							{@const Icon = item.icon}
							<SidebarMenuItem>
								<SidebarMenuButton
									isActive={isActive(item.href)}
									tooltipContent={item.label}
								>
									{#snippet child({ props })}
										<a href={item.href} {...props}>
											<Icon />
											<span>{item.label}</span>
										</a>
									{/snippet}
								</SidebarMenuButton>
							</SidebarMenuItem>
						{/each}
					</SidebarMenu>
				</SidebarGroupContent>
			</SidebarGroup>
		</SidebarContent>

		<SidebarFooter>
			<div class="flex flex-col gap-2">
				<div class="flex items-center gap-2 rounded-md bg-sidebar p-2">
					<Avatar class="size-8">
						<AvatarFallback class="bg-sidebar-accent text-sidebar-accent-foreground text-xs">
							{userInitials}
						</AvatarFallback>
					</Avatar>
					<div class="min-w-0 flex-1">
						<span class="block truncate text-sm font-medium text-sidebar-foreground" id="user-name"></span>
						<span class="block truncate text-xs text-muted-foreground" id="user-email"></span>
					</div>
				</div>
				<div class="flex items-center justify-between px-2">
					<div class="flex items-center gap-1">
						<LanguageSwitcher />
						<Button
							variant="ghost"
							size="sm"
							onclick={() => themeStore.toggle()}
							class="p-1.5 text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
							title={themeStore.isDark ? t('layout.sidebar.lightMode') : t('layout.sidebar.darkMode')}
						>
							{#if themeStore.isDark}
								<Sun class="size-3.5" />
							{:else}
								<Moon class="size-3.5" />
							{/if}
						</Button>
					</div>
					<Button
						variant="ghost"
						size="sm"
						onclick={handleLogout}
						class="gap-1.5 text-xs text-muted-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
						title={t('layout.sidebar.logout')}
					>
						<LogOut class="size-3.5" />
						<span>{t('layout.sidebar.logout')}</span>
					</Button>
				</div>
			</div>
		</SidebarFooter>

		<SidebarRail />
	</Sidebar>

	<SidebarInset>
		<header class="flex h-14 shrink-0 items-center gap-2 border-b px-4">
			<SidebarTrigger class="-ml-1" />
			<Separator orientation="vertical" class="mr-2 h-4" />
			<h2 class="text-sm font-semibold text-foreground">{t('app.name')}</h2>
		</header>

		<main class="flex-1 overflow-y-auto p-6">
			{@render children()}
		</main>
	</SidebarInset>
</SidebarProvider>
