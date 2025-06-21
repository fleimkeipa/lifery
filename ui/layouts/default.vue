<template>
  <div class="flex h-screen w-full flex-col items-center justify-start">
    <div class="w-2/3">
      <div v-if="$route.path !== '/intro'" class="flex justify-between items-center py-4">
        <navigation-menu />
        <div class="flex items-center gap-2">
          <NotificationBell />
          <UButton
            :to="locale === 'tr' ? '/user/settings' : '/en/user/settings'"
            color="gray"
            variant="ghost"
            icon="i-heroicons-user-circle"
            :aria-label="$t('settings.title')"
            class="p-2"
          />
          <LanguageSwitcher />
          <UButton color="red" variant="ghost" @click="logout">{{ $t('common.logout') }}</UButton>
        </div>
      </div>
      <slot class="mt-4" />
    </div>
  </div>
</template>

<script setup>
const { locale } = useI18n();

const logout = () => {
  localStorage.removeItem('auth_token')
  window.location.href = '/login'
}
</script>
