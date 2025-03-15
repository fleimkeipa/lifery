<template>
  <div class="flex h-screen w-full flex-col items-center justify-start">
    <div class="w-full">
      <div class="flex justify-end p-4">

        <LanguageSwitcher />
      </div>
      <slot class="mt-4" />
    </div>
  </div>
</template>
