<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    children: Snippet;
    title?: string;
    count?: number;
  }

  let { children, title = 'Patients', count }: Props = $props();
</script>

<div class="flex flex-col h-full bg-white border-r border-gray-200">
  <!-- Header -->
  <div class="px-4 py-3 border-b border-gray-200">
    <div class="flex items-center justify-between mb-3">
      <h2 class="text-sm font-semibold text-gray-700 uppercase tracking-wider">{title}</h2>
      {#if count !== undefined}
        <span class="text-xs font-medium bg-gray-100 text-gray-600 rounded-full px-2">{count}</span>
      {/if}
    </div>
    <input
      type="text"
      placeholder="Search patients..."
      class="w-full text-sm border border-gray-200 rounded-lg px-3 py-2 text-gray-700 placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-primary-500 focus:border-primary-500"
    />
  </div>

  <!-- List -->
  <ul class="flex-1 overflow-y-auto divide-y divide-gray-100">
    {@render children()}
  </ul>
</div>