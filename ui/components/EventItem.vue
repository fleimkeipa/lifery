<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

enum EventType {
  "STRING" = 10,
  "PHOTO" = 11,
  "VIDEO" = 12,
  "VOICE_RECORD" = 13
}

defineProps<{
    id: string;
    date: string;
    cards: {
        title: string;
        description?: string;
        color: string;
        items?: {
        data: string;
        type: number;
        }[];
    }[];
    isLeft: boolean;
}>()
</script>

<template>
    <div v-for="card in cards" :key="card.title" class="flex gap-x-2">
        <EventCardItem v-if="card.items && isLeft" :items="card.items" :color="card.color" />
        <div :class="[card.color, 'p-4 rounded shadow-lg max-w-[300px] aspect-[4/3] flex flex-col mb-4']">
            <h3 class="font-medium text-gray-800 mb-2">{{ card.title }}</h3>
            <p class="text-sm text-gray-600 overflow-auto">{{ card.description || t('common.no_description') }}</p>
        </div>
        <EventCardItem v-if="card.items && !isLeft" :items="card.items" :color="card.color" />
    </div>
</template>