<template>
  <span class="font-extrabold text-2xl cursor-pointer select-none text-blue-700 dark:text-blue-300" @click="goHome">
    Lifery
  </span>
</template>

<script setup lang="ts">

function goHome() {
  navigateTo('/');
}
</script>