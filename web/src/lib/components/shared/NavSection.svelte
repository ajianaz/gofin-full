<script lang="ts">
	import {
		SidebarGroup,
		SidebarGroupContent,
		SidebarGroupLabel,
		SidebarMenu,
		SidebarMenuItem,
		SidebarMenuButton
	} from '$lib/components/ui/sidebar/index.js';
	import type { Snippet } from 'svelte';
	import { page } from '$app/stores';

	interface NavItem {
		label: string;
		href: string;
		icon: any;
	}

	let {
		label,
		items
	}: {
		label: string;
		items: NavItem[];
	} = $props();

	function isActive(href: string): boolean {
		return $page.url.pathname.startsWith(href);
	}
</script>

<SidebarGroup>
	<SidebarGroupLabel>{label}</SidebarGroupLabel>
	<SidebarGroupContent>
		<SidebarMenu>
			{#each items as item}
				{@const Icon = item.icon}
				<SidebarMenuItem>
					<SidebarMenuButton isActive={isActive(item.href)} tooltipContent={item.label}>
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
