<script setup lang="ts">
import { useI18n } from 'vue-i18n';
const { t, locale } = useI18n();

import IntroCard from '@/components/IntroCard.vue';

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});
</script>

<template>
  <div class="min-h-screen w-full bg-white dark:bg-gray-950">
    <div class="flex flex-col items-center justify-center py-8 p-4 z-10">
      <div class="max-w-xl w-full bg-white/80 dark:bg-gray-900/80 rounded-2xl shadow-xl p-10 flex flex-col items-center gap-8">
        <h1 class="text-5xl font-extrabold text-blue-700 dark:text-blue-300 mb-2">Lifery</h1>
        <p class="text-lg text-gray-700 dark:text-gray-100 text-center mb-4">
          {{ t('intro.title') }}
        </p>
        <div class="w-full flex flex-col gap-4 mt-4">
          <div class="flex flex-col gap-4 w-full justify-center">
            <IntroCard
              bgClass="bg-blue-50"
              darkBgClass="dark:bg-blue-900"
              borderClass="border border-blue-100"
              darkBorderClass="dark:border-blue-800"
              emoji="📝"
              :title="t('intro.eventTitle')"
              :example="t('intro.eventExample')"
              :desc="t('intro.eventDesc')"
            />
            <IntroCard
              bgClass="bg-purple-50"
              darkBgClass="dark:bg-purple-900"
              borderClass="border border-purple-100"
              darkBorderClass="dark:border-purple-800"
              emoji="📅"
              :title="t('intro.eraTitle')"
              :example="t('intro.eraExample')"
              :desc="t('intro.eraDesc')"
            />
            <IntroCard
              bgClass="bg-green-50"
              darkBgClass="dark:bg-green-900"
              borderClass="border border-green-100"
              darkBorderClass="dark:border-green-800"
              emoji="🤝"
              :title="t('intro.connectTitle')"
              :example="t('intro.connectExample')"
              :desc="t('intro.connectDesc')"
            />
          </div>
        </div>
        <div class="flex gap-4 mt-6">
          <UButton to="/login" size="xl" color="primary">{{ t('login.title') }}</UButton>
          <UButton to="/register" size="xl" color="gray" variant="outline">{{ t('register.title') }}</UButton>
        </div>
      </div>
    </div>
  </div>
</template> 