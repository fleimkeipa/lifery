<template>
  <div class="flex h-screen w-full flex-col items-center justify-start">
    <div class="w-2/3">
      <div class="flex justify-between items-center py-4">
        <navigation-menu />
        <LanguageSwitcher />
      </div>
      <slot class="mt-4" />
    </div>
  </div>
</template>
