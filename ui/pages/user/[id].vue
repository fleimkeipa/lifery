<script setup lang="ts">
import { ref } from 'vue';

definePageMeta({
  middleware: "auth",
});

const { t } = useI18n();

const route = useRoute();
const id = route.params.id

interface Row {
  id: number;
  name: string;
  description: string;
  visibility: number;
  date: Date;
  time_start: Date;
  time_end: Date;
  items: {
    data: string;
    type: number;
  }[];
}

export interface TimelineItem {
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
  type: 'EVENT'
}

interface TimelineEra {
  id: string;
  color: string;
  name: string;
  time_end: string;
  time_start: string;
  created_at: string;
  updated_at: string;
  date: string;
  type: 'ERA'
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString('tr-TR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).replace(/\//g, '.');
};

const getTextColor = (hexColor: string) => {
  // Remove # if it exists
  const hex = hexColor.replace('#', '');
  
  // Convert hex to RGB
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);
  
  // Calculate relative luminance
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
  
  // Return dark text for light backgrounds, light text for dark backgrounds
  return luminance > 0.5 ? 'text-gray-900' : 'text-white';
};

const { data: eventsData, error, isFetching, execute: fetchEvents } = useApi<{
  data: Row[];
  total: number;
  limit: number;
  skip: number;
}>(`/events?user_id=${id}`).json()

const { data: erasData, error: errorEras, isFetching: isFetchingEras, execute: fetchEras } = useApi<{
  data: Row[];
  total: number;
  limit: number;
  skip: number;
}>(`/eras?order=desc:time_start&user_id=${id}`).json();

const timelineData = computed<(TimelineItem | TimelineEra)[]>(() => {
  if (!eventsData.value?.data) return [];

  return [
    ...eventsData.value.data.map((event: Row) => ({
      id: event.id.toString(),
      date: event.date,
      cards: [
        {
          title: event.name,
          description: event.description,
          color: event.visibility === 3 ? 'bg-pink-300' : event.visibility === 2 ? 'bg-yellow-200' : 'bg-blue-200',
          items: event.items
        }
      ],
      type: 'EVENT'
    })),
    ...erasData.value.data.map((event: TimelineEra) => ({
      id: event.id,
      color: event.color,
      name: event.name,
      date: event.time_start,
      type: 'ERA'
    })),
    ...erasData.value.data.map((event: TimelineEra) => ({
      id: event.id,
      color: event.color,
      name: event.name,
      date: event.time_end,
      type: 'ERA'
    }))
  ].sort((a: TimelineItem | TimelineEra, b: TimelineItem | TimelineEra) => {

    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
});

const toast = ref<{ show: boolean; message: string; type: 'success' | 'error' }>({
  show: false,
  message: '',
  type: 'success'
});

function showToast(message: string, type: 'success' | 'error' = 'success') {
  toast.value = { show: true, message, type };
  setTimeout(() => {
    toast.value.show = false;
  }, 3000);
}

const addConnection = async () => {
  try {
    await useApi('/connects').post({
      friend_id: id
    });
    showToast(t('connect.connection_added'), 'success');
  } catch (error) {
    console.error(t('connect.connection_error'), error);
    showToast(t('connect.connection_error'), 'error');
  }
};
</script>

<template>
  <div class="min-h-screen bg-background p-8">
    <div class="max-w-6xl mx-auto">
      <!-- Bağlantı Ekle Butonu -->
      <div class="mb-8 flex justify-end">
        <button @click="addConnection"
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
          {{ t(`connect.add_connection`) }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="isFetching" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center text-red-600 p-4">
        <p>Error loading events: {{ error }}</p>
        <button @click="() => fetchEvents()" class="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
          Retry
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="!timelineData.length" class="text-center text-gray-600 p-4">
        <p>{{ t('event.no_events') }}</p>
      </div>

      <!-- Timeline Container -->
      <div v-else class="relative">
        <!-- Vertical Line -->
        <div class="absolute left-1/2 h-full w-0.5 bg-gray-200"></div>

        <!-- Timeline Items -->
        <div class="space-y-24">
          <div v-for="item in timelineData" :key="item.id" class="relative">
            <!-- Date Marker -->
            <div v-if="item.type === 'EVENT'"
              class="absolute left-1/2 -translate-x-1/2 w-4 h-4 rounded-full bg-gray-200 z-10"></div>

            <!-- Date Text -->
            <div v-if="item.type === 'EVENT'" class="absolute left-[52%] top-[-1.5rem] text-sm text-foreground">
              {{ formatDate(item.date) }}
            </div>

            <!-- Cards Container -->
            <div v-if="item.type === 'EVENT'" class="grid grid-cols-2 gap-8">
              <div v-if="timelineData.indexOf(item) % 2 === 1"></div>
              <div class="flex justify-end">
                <div class="transform hover:rotate-0 transition-transform duration-200 -rotate-2"
                  :class="timelineData.indexOf(item) % 2 === 0 ? '-rotate-2' : 'rotate-2'">
                  <EventItem :id="item.id" :date="item.date" :cards="item.cards"
                    :is-left="timelineData.indexOf(item) % 2 === 0" />
                </div>
              </div>
              <div v-if="timelineData.indexOf(item) % 2 === 0"></div>
            </div>
            <div v-else class="shadow-lg w-full mb-4 px-4 z-20 text-center flex justify-between items-center"
              :style="{ background: item.color }">
              <div class="font-medium text-xs" :class="getTextColor(item.color)">{{ formatDate(item.date) }}</div>
              <h3 class="font-bold text-sm" :class="getTextColor(item.color)">{{ item.name }}</h3>
            </div>
          </div>
        </div>
      </div>

      <div v-if="toast.show" :class="[
        'fixed top-6 left-1/2 transform -translate-x-1/2 px-6 py-3 rounded shadow-lg z-50 transition-all',
        toast.type === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'
      ]">
        {{ toast.message }}
      </div>
    </div>
  </div>
</template>