<template>
  <div :class="[
    bgClass,
    darkBgClass,
    borderClass,
    darkBorderClass,
    'rounded-xl p-4 shadow flex flex-col items-center'
  ]">
    <span class="text-3xl mb-2">{{ emoji }}</span>
    <div class="font-bold mb-1 text-gray-800 dark:text-gray-100">{{ title }}</div>
    <div class="text-sm text-gray-600 dark:text-gray-200 text-center mb-2">{{ example }}</div>
    <div class="w-full bg-white dark:bg-gray-800 rounded p-2 text-xs text-gray-500 dark:text-gray-300" :class="[borderClass, darkBorderClass]">
      {{ desc }}
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  bgClass: string;
  darkBgClass: string;
  borderClass: string;
  darkBorderClass: string;
  emoji: string;
  title: string;
  example: string;
  desc: string;
}>()
</script> 