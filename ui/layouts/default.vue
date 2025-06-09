<template>
  <div class="flex h-screen w-full flex-col items-center justify-start">
    <div class="w-2/3">
      <div class="flex justify-between items-center py-4">
        <navigation-menu />
        <div class="flex items-center gap-2">
          <LanguageSwitcher />
          <UButton color="red" variant="ghost" @click="logout">{{ $t('common.logout') }}</UButton>
        </div>
      </div>
      <slot class="mt-4" />
    </div>
  </div>
</template>

<script setup>
const logout = () => {
  localStorage.removeItem('auth_token')
  window.location.href = '/login'
}
</script>
