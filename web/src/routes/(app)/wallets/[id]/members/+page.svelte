<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { BackButton } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { walletMemberService } from '$lib/services/wallet-members.js';
	import type { WalletMember } from '$lib/types/domain.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
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
		<div class="overflow-x-auto">
			<table class="w-full text-sm">
				<thead>
					<tr class="border-b bg-muted/50">
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.role')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.name')}</th>
					</tr>
				</thead>
				<tbody>
					{#if isLoading}
						<tr>
							<td colspan="2" class="p-8 text-center text-sm text-muted-foreground">{t('common.loading')}</td>
						</tr>
					{:else if errorMsg}
						<tr>
							<td colspan="2" class="p-8 text-center text-sm text-destructive">{errorMsg}</td>
						</tr>
					{:else if members.length === 0}
						<tr>
							<td colspan="2" class="p-8 text-center text-sm text-muted-foreground">{t('common.noData')}</td>
						</tr>
					{:else}
						{#each members as member}
							<tr class="border-b hover:bg-muted/30">
								<td class="p-3">
									<Badge variant={roleVariant(member.role)}>
										{member.role}
									</Badge>
								</td>
								<td class="p-3 font-mono text-muted-foreground">{member.user_id}</td>
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>
