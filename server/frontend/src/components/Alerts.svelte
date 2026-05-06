<script lang="ts">
    import type { Snippet } from 'svelte';
    interface Props {
        children: Snippet;
    }
     let { children }: Props = $props();
</script>

<div class="flex flex-col justify-start p-4 space-x-1 space-y-2 bg-gray-100 items-center">
    {@render children()}
</div>