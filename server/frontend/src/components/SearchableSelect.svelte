<script lang="ts">
    // Credit to https://github.com/Rolleander
    // ref: https://github.com/themesberg/flowbite/issues/1130
	import { Input } from 'flowbite-svelte';
	import { onDestroy, onMount } from 'svelte';
    import { browser } from '$app/environment';

	type Option = {
		label: string;
		value: string;
	};

	let {
		options = [],
		value = $bindable<string>(''),
		placeholder = '',
		disabled = false,
		clearable = false,
		size = 'sm'
	}: {
		size?: 'sm' | 'md' | 'lg';
		options: Option[];
		value: string;
		placeholder?: string;
		disabled?: boolean;
		clearable?: boolean;
	} = $props();

	let open = $state(false);
	let query = $state('');
	let rect = $state<DOMRect | null>(null);
	let maxHeight = $state<number>(0);
	let openUpwards = $state(false);

	let container: HTMLDivElement;

	const selected = () => options.find((o) => o.value === value) ?? null;

	const filtered = () => options.filter((o) => o.label.toLowerCase().includes(query.toLowerCase()));

	function select(option: Option) {
		value = option.value;
		query = option.label;
		open = false;
	}

	function handleFocus() {
		applyPopupSize();
		open = true;
		query = '';
	}

	function applyPopupSize() {
		rect = container.getBoundingClientRect();
		const spaceBelow = window.innerHeight - rect.bottom - 8;
		const spaceAbove = rect.top - 8;
		if (spaceBelow <= 150 && spaceAbove > 150) {
			openUpwards = true;
			maxHeight = spaceAbove;
		} else {
			openUpwards = false;
			maxHeight = spaceBelow;
		}
	}

	function handleClickOutside(e: MouseEvent) {
		if (
			!container.contains(e.target as Node) &&
			!(e.target as HTMLElement).closest('[data-searchable-select-dropdown]')
		) {
			open = false;
		}
	}

	$effect(() => {
		if (!open && selected()) {
			query = selected()!.label;
		}
	});

	onMount(() => {
        if (!browser) return;
		document.addEventListener('mousedown', handleClickOutside);
		window.addEventListener('resize', applyPopupSize);
	});

	onDestroy(() => {
        if (!browser) return;
		document.removeEventListener('mousedown', handleClickOutside);
		window.removeEventListener('resize', applyPopupSize);
	});
</script>

<div class="relative w-full" bind:this={container}>
	<Input
		{size}
		bind:value={query}
		{placeholder}
		{disabled}
		onfocus={handleFocus}
		oninput={() => (open = true)}
		{clearable}
	>
		{#snippet right()}
			<div class="text-sm font-semibold">
				<i class="fa-solid fa-angle-down"></i>
			</div>
		{/snippet}
	</Input>
</div>

{#if open && rect}
	<div
		data-searchable-select-dropdown
		class="fixed z-9999 overflow-auto rounded-lg
           border border-gray-200 bg-white shadow-md"
		style={`
      left: ${rect.left}px;
      width: ${rect.width}px;
      max-height: ${maxHeight}px;
         ${openUpwards ? `bottom: ${window.innerHeight - rect.top + 4}px;` : `top: ${rect.bottom + 4}px;`}
    `}
	>
		{#if filtered().length === 0}
			<div class="px-3 py-2 text-sm text-gray-500">-</div>
		{:else}
			{#each filtered() as option, i (option.value)}
				<button
					type="button"
					class={`block w-full px-3 ${size == 'sm' ? 'text-xs py-1' : 'text-sm py-2'} text-left 
      ${i % 2 === 0 ? 'bg-gray-50' : 'bg-white'}
      hover:bg-gray-100 focus:bg-gray-100`}
					onclick={() => select(option)}
				>
					{option.label}
				</button>
			{/each}
		{/if}
	</div>
{/if}

