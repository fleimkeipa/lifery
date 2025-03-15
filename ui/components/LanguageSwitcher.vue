<template>
    <UDropdown :items="[dropdownItems]" :ui="{ item: { base: 'flex items-center gap-x-2 cursor-pointer' } }">
        <UButton color="gray" variant="ghost" class="flex items-center gap-2">
            {{ currentLocale.label }}
            <UIcon name="i-heroicons-chevron-down-20-solid" />
        </UButton>

        <template #item="{ item }">
            <div @click="switchLanguage(item.value)" class="w-full">{{ item.label }}</div>
        </template>
    </UDropdown>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'

interface LocaleItem {
    label: string
    value: string
}

const { locale, setLocale } = useI18n()
const route = useRoute()
const router = useRouter()

const dropdownItems: LocaleItem[] = [
    { label: 'Türkçe', value: 'tr' },
    { label: 'English', value: 'en' }
]

const currentLocale = computed(() => {
    return dropdownItems.find(l => l.value === locale.value) || dropdownItems[0]
})

const switchLanguage = async (code: "tr" | "en") => {
    const query = { ...route.query, locale: code }
    await router.push({ path: route.path, query })
    setLocale(code)
}
</script>