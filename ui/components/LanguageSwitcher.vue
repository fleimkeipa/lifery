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
import { onMounted, computed, watch } from 'vue'

interface LocaleItem {
    label: string
    value: string
}

const { locale, setLocale } = useI18n()
const route = useRoute()
const router = useRouter()

const dropdownItems: LocaleItem[] = [
    { label: '🇹🇷', value: 'tr' },
    { label: '🇬🇧', value: 'en' }
]

const currentLocale = computed(() => {
    return dropdownItems.find(l => l.value === locale.value) || dropdownItems[0]
})

const switchLanguage = async (code: "tr" | "en") => {
    if (process.client) {
        // Save the preference
        localStorage.setItem('preferred_locale', code)
    }
    
    // Get the current path without the locale prefix
    const path = route.path.replace(/^\/[a-z]{2}(?=\/|$)/, '')
    
    // Set the new locale
    setLocale(code)
    
    // Navigate to the new path with the locale prefix
    if (code === 'tr') {
        // For Turkish (default locale), remove the prefix
        await router.replace(path)
    } else {
        // For other locales, add the prefix
        await router.replace(`/${code}${path}`)
    }
}

// Watch for route changes to maintain language preference
watch(() => route.path, () => {
    if (process.client) {
        const savedLocale = localStorage.getItem('preferred_locale')
        if (savedLocale && (savedLocale === 'tr' || savedLocale === 'en') && savedLocale !== locale.value) {
            switchLanguage(savedLocale)
        }
    }
}, { immediate: true })

// Initialize locale from localStorage on mount
onMounted(() => {
    if (process.client) {   
        const savedLocale = localStorage.getItem('preferred_locale')
        if (savedLocale && (savedLocale === 'tr' || savedLocale === 'en')) {
            switchLanguage(savedLocale)
        }
    }
})
</script>