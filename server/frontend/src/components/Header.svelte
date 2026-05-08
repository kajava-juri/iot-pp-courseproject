<script lang="ts">
    import { alertsCount } from '$stores/alerts';
</script>

<header class="flex p-4 bg-gray-200 text-white justify-between ">
    <h1 class="text-2xl font-bold text-gray-800">IoT Dashboard</h1>
    <div class="flex items-center">
        <button class="rounded-full bg-primary-500 hover:bg-primary-600 text-white px-2 font-semibold tracking-wide">
            Alerts ({$alertsCount})
        </button>
    </div>
</header>