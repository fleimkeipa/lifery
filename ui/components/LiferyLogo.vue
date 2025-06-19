<template>
  <span
    class="font-extrabold text-2xl cursor-pointer select-none text-blue-700 dark:text-blue-300"
    @click="goHome"
  >
    Lifery
  </span>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
const router = useRouter();

function goHome() {
  if (process.client) {
    const token = localStorage.getItem('auth_token');
    if (token) {
      router.push('/');
    } else {
      router.push('/intro');
    }
  }
}
</script> 