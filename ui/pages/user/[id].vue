<script setup lang="ts">
const route = useRoute();
const id = route.params.id

const { data } = useApi(`/events?user_id=${id}`)
</script>

<template>
    {{ data }}
</template>