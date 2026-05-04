<script lang="ts">
  type Severity = 'high' | 'medium' | 'low' | 'none';

  interface Props {
    name: string;
    patientId: string;
    roomName: string;
    severity?: Severity;
    lastEvent?: string;
    lastEventTime?: string;
    selected?: boolean;
    onclick?: () => void;
  }

  let {
    name,
    patientId,
    roomName,
    severity = 'none',
    lastEvent,
    lastEventTime,
    selected = false,
    onclick,
  }: Props = $props();

  const dotColor: Record<Severity, string> = {
    high:   'bg-red-500',
    medium: 'bg-amber-400',
    low:    'bg-yellow-300',
    none:   'bg-emerald-500',
  };
</script>

<li>
  <button
    type="button"
    {onclick}
    class="w-full text-left flex items-center gap-3 px-4 py-3 transition-colors hover:bg-gray-50
      {selected ? 'bg-primary-50 border-l-2 border-primary-600' : 'border-l-2 border-transparent'}"
  >
    <!-- Status dot -->
    <span class="w-2.5 h-2.5 rounded-full shrink-0 {dotColor[severity]}"></span>

    <!-- Info -->
    <div class="flex-1 min-w-0">
      <p class="text-sm font-medium text-gray-900 truncate">{name}</p>
      <p class="text-xs text-gray-500 truncate">
        ID: {patientId}
        {#if lastEvent}
          · <span class="text-gray-400">{lastEvent}</span>
        {/if}
        {#if lastEventTime}
          · <span class="text-gray-400">{lastEventTime}</span>
        {/if}
      </p>
    </div>

    <!-- Room badge -->
    <span class="shrink-0 text-xs font-medium text-gray-500 bg-gray-100 border border-gray-200 rounded px-1.5 py-0.5 whitespace-nowrap">
      {roomName}
    </span>
  </button>
</li>