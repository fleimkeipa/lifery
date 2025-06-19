<template>
  <div class="flex h-screen w-full flex-col items-center justify-start">
    <div class="w-full">
      <div class="flex justify-between items-center p-4">
        <LiferyLogo />
        <LanguageSwitcher />
      </div>
      <slot class="mt-4" />
    </div>
  </div>
</template>
