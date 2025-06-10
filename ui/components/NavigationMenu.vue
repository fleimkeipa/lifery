<script setup lang="ts">
const { t } = useI18n();

const links = [
  {
    label: t('common.home'),
    icon: "i-heroicons-home",
    to: "/",
  },
  {
    label: t('common.events'),
    icon: "i-heroicons-circle-stack",
    to: "/events",
  },
  {
    label: t('common.eras'),
    icon: "i-heroicons-cube",
    to: "/eras",
  },
  {
    label: t('common.connects'),
    icon: "i-heroicons-folder-open",
    to: "/connects",
  },
  {
    label: t('common.explore'),
    icon: "i-heroicons-magnifying-glass",
    to: "/explore",
  },
];
</script>

<template>
  <div>
    <UHorizontalNavigation
      :links="links"
      class="border-b border-gray-200 dark:border-gray-800"
    />
  </div>
</template>
