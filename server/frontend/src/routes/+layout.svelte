<script>
  let { children } = $props();
  import "../app.css";

  import Header from "$components/Header.svelte";
  import Sidemenu from "$components/Sidemenu.svelte";
    import Alerts from "$components/Alerts.svelte";
    import AlertItem from "$components/AlertItem.svelte";
</script>

<div class="h-screen flex flex-col bg-white text-gray-900">
    <Header />

    <div class="flex flex-1">
        <Sidemenu />
        {@render children()}
    </div>
</div>