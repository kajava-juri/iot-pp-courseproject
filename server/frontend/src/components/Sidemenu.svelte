<script lang="ts">
    import { page } from "$app/state";
    import {Sidebar, SidebarGroup, SidebarItem, SidebarButton, uiHelpers} from "flowbite-svelte"
    import {CirclePlusOutline, ListOutline, AdjustmentsVerticalOutline, CalendarMonthOutline} from "flowbite-svelte-icons"
    let activeUrl = $state(page.url.pathname);
    const activeClass = "flex items-center p-2 text-base font-normal text-white bg-primary-600 dark:bg-primary-700 rounded-lg dark:text-white hover:bg-primary-800 dark:hover:bg-primary-800";
    const nonActiveClass = "flex items-center p-2 text-base font-normal text-green-900 rounded-lg dark:text-white hover:bg-green-100 dark:hover:bg-green-700";

    $effect(() => {
        activeUrl = page.url.pathname;
    })
</script>

<div class="relative">
    <Sidebar
        activeUrl={activeUrl}
        alwaysOpen={true}
        class="z-50 h-full relative bg-gray-100 *:bg-gray-100"
        classes={{ nonactive: nonActiveClass, active: activeClass, div: "overflow-y-auto px-3 py-4" }}
        >
        <SidebarGroup>
            <SidebarItem label="Add Patient" href="/patient">
                {#snippet icon()}
                    <CirclePlusOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                {/snippet}
            </SidebarItem>
            <SidebarItem label="Patient List" href="/">
                {#snippet icon()}
                    <ListOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                {/snippet}
            </SidebarItem>
            <SidebarItem label="System Status" href="/system_status">
                {#snippet icon()}
                    <AdjustmentsVerticalOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                {/snippet}
            </SidebarItem>
            <SidebarItem label="Events" href="/event">
                {#snippet icon()}
                    <CalendarMonthOutline class="h-5 w-5 text-gray-500 transition duration-75 group-hover:text-gray-900 dark:text-gray-400 dark:group-hover:text-white" />
                {/snippet}

            </SidebarItem>
        </SidebarGroup>
    </Sidebar>
</div>