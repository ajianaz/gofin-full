<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '$lib/components/ui/table/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { walletMemberService } from '$lib/services/wallet-members.js';
	import type { WalletMember } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import EmptyState from '$lib/components/shared/EmptyState.svelte';
	const t = localeStore.t;

	let walletId = $derived($page.params.id!);
	let members = $state<WalletMember[]>([]);
	let isLoading = $state(true);
	let errorMsg = $state('');

	function roleVariant(role: string) {
		if (role === 'owner') return 'default';
		if (role === 'editor') return 'secondary';
		return 'outline';
	}

	onMount(async () => {
		try {
			members = await walletMemberService.list(walletId);
		} catch (e) {
			errorMsg = t('common.error');
			console.error(e);
		} finally {
			isLoading = false;
		}
	});
</script>

<BackButton href="/wallets" label={t('wallets.title')} />

<div class="mb-6">
	<h1 class="text-2xl font-bold text-foreground">{t('wallets.members.title')}</h1>
	<p class="text-sm text-muted-foreground mt-0.5">{t('wallets.members.description')}</p>
</div>

<Card>
	<CardContent class="p-0">
		<Table>
				<TableHeader>
					<TableRow>
						<TableHead>{t('wallets.members.role')}</TableHead>
						<TableHead>{t('wallets.members.name')}</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					{#if isLoading}
				{#each Array(5) as _}
<TableRow>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
					<TableCell><Skeleton class="h-4 w-full" /></TableCell>
				</TableRow>
				{/each}
	{:else if errorMsg}
						<TableRow>
							<TableCell colspan="2" class="p-8 text-center text-sm text-destructive">{errorMsg}</TableCell>
						</TableRow>
					{:else if members.length === 0}
						<TableRow>
							<TableCell colspan="2"><EmptyState /></TableCell>
						</TableRow>
					{:else}
						{#each members as member}
							<TableRow>
								<TableCell>
									<Badge variant={roleVariant(member.role)}>
										{member.role}
									</Badge>
								</TableCell>
								<TableCell class="font-mono text-muted-foreground">{member.user_id}</TableCell>
							</TableRow>
						{/each}
					{/if}
				</TableBody>
			</Table>
	</CardContent>
</Card>
