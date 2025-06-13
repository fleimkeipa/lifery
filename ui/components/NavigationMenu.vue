<script setup lang="ts">
const { t, locale } = useI18n();
const route = useRoute();

const links = computed(() => [
  {
    label: t('common.home'),
    icon: "i-heroicons-home",
    to: locale.value === 'tr' ? "/" : "/en",
  },
  {
    label: t('common.events'),
    icon: "i-heroicons-circle-stack",
    to: locale.value === 'tr' ? "/events" : "/en/events",
  },
  {
    label: t('common.eras'),
    icon: "i-heroicons-cube",
    to: locale.value === 'tr' ? "/eras" : "/en/eras",
  },
  {
    label: t('common.connects'),
    icon: "i-heroicons-folder-open",
    to: locale.value === 'tr' ? "/connects" : "/en/connects",
  },
  {
    label: t('common.explore'),
    icon: "i-heroicons-magnifying-glass",
    to: locale.value === 'tr' ? "/explore" : "/en/explore",
  },
]);
</script>

<template>
  <div>
    <UHorizontalNavigation
      :links="links"
      class="border-b border-gray-200 dark:border-gray-800"
    />
  </div>
</template>
