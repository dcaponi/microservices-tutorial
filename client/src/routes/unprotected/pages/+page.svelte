<script lang="ts">
    export let data;

    console.log(data)
</script>

<h1>This is a page</h1>
<h2>resume svc response</h2>
<p>{data.res.ping}</p>
<h2>candidate svc response</h2>
<p>{data.can.ping}</p>