<script lang="ts">
	import { BackButton, StatusBadge } from '$lib/components/shared/index.js';
	import { Card, CardContent } from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { localeStore } from '$lib/stores/i18n.svelte.js';
	const t = localeStore.t;

	const mockMembers = [
		{ name: 'Anaz S. Aji', email: 'anaz@gofin.app', role: 'owner', joined: '2026-01-01' },
		{ name: 'Budi Santoso', email: 'budi@email.com', role: 'member', joined: '2026-02-15' },
		{ name: 'Citra Dewi', email: 'citra@email.com', role: 'member', joined: '2026-03-01' },
		{ name: 'Dimas Pratama', email: 'dimas@email.com', role: 'viewer', joined: '2026-03-10' }
	];
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
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.name')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.email')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.role')}</th>
						<th class="text-left p-3 font-medium text-muted-foreground">{t('wallets.members.joined')}</th>
					</tr>
				</thead>
				<tbody>
					{#each mockMembers as member}
						<tr class="border-b hover:bg-muted/30">
							<td class="p-3 font-medium text-foreground">{member.name}</td>
							<td class="p-3 text-muted-foreground">{member.email}</td>
							<td class="p-3">
								<Badge variant={member.role === 'owner' ? 'default' : member.role === 'member' ? 'secondary' : 'outline'}>
									{member.role}
								</Badge>
							</td>
							<td class="p-3 text-muted-foreground">{member.joined}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</CardContent>
</Card>
