<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { useAuth } from '../../composables/useAuth';

const { user } = useAuth();

definePageMeta({
  middleware: "custom-layout",
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
  date: string;
  type: 'ERA'
}

interface ApiResponse {
  data: Row[];
  total: number;
  limit: number;
  skip: number;
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
  const hex = hexColor.replace('#', '');
  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255;
  return luminance > 0.5 ? 'text-gray-900' : 'text-white';
};

// Pagination state for events
const limit = ref(10);
const skip = ref(0);
const total = ref(0);
const hasMore = ref(true);
const isLoadingMore = ref(false);
const error = ref<string | null>(null);

// Store all events data
const allEventsData = ref<Row[]>([]);
const allErasData = ref<TimelineEra[]>([]);

// Fetch events with pagination
const fetchEventsWithPagination = async (reset = false) => {
  if (reset) {
    skip.value = 0;
    allEventsData.value = [];
    hasMore.value = true;
    error.value = null;
  }

  if (!hasMore.value || isLoadingMore.value) return;

  isLoadingMore.value = true;
  
  try {
    const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
    const config = useRuntimeConfig();
    
    const response = await $fetch<ApiResponse>(`${config.public.apiBase}/events?user_id=${id}&order=desc:date&limit=${limit.value}&skip=${skip.value}`, {
      headers: {
        'Authorization': token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json',
      }
    });
    
    if (reset) {
      allEventsData.value = response.data;
    } else {
      allEventsData.value = [...allEventsData.value, ...response.data];
    }
    
    total.value = response.total;
    skip.value += limit.value;
    hasMore.value = skip.value < total.value;
  } catch (err) {
    console.error('Error fetching events:', err);
    error.value = err instanceof Error ? err.message : 'An error occurred while fetching events';
  } finally {
    isLoadingMore.value = false;
  }
};

// Fetch eras
const { data: erasData, execute: fetchEras } = useApi(() => `eras?order=desc:time_start&user_id=${id}`).json();

// Watch eras data and update allErasData
watch(erasData, (newErasData) => {
  if (newErasData?.data) {
    allErasData.value = newErasData.data;
  }
});

// Intersection observer for infinite scroll
const loadMoreTrigger = ref<HTMLElement>();
let observer: IntersectionObserver | null = null;

const setupIntersectionObserver = () => {
  if (loadMoreTrigger.value && !observer) {
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && hasMore.value && !isLoadingMore.value) {
            fetchEventsWithPagination();
          }
        });
      },
      { threshold: 0.1 }
    );

    observer.observe(loadMoreTrigger.value);
  }
};

// Initial fetch and setup
onMounted(() => {
  fetchEventsWithPagination(true);
  fetchEras();
  
  nextTick(() => {
    setupIntersectionObserver();
  });
});

onUnmounted(() => {
  if (observer) {
    observer.disconnect();
    observer = null;
  }
});

const timelineData = computed<(TimelineItem | TimelineEra)[]>(() => {
  if (!allEventsData.value.length) return [];

  return [
    ...allEventsData.value.map((event: Row) => ({
      id: event.id.toString(),
      date: typeof event.date === 'string' ? event.date : new Date(event.date).toISOString(),
      cards: [
        {
          title: event.name,
          description: event.description,
          color: event.visibility === 3 ? 'bg-pink-300' : event.visibility === 2 ? 'bg-yellow-200' : 'bg-blue-200',
          items: event.items
        }
      ],
      type: 'EVENT' as const
    })),
    ...allErasData.value.map((event: TimelineEra) => ({
      id: event.id,
      color: event.color,
      name: event.name,
      time_end: event.time_end,
      time_start: event.time_start,
      date: event.time_start,
      type: 'ERA' as const
    })),
    ...allErasData.value.map((event: TimelineEra) => ({
      id: event.id,
      color: event.color,
      name: event.name,
      time_end: event.time_end,
      time_start: event.time_start,
      date: event.time_end,
      type: 'ERA' as const
    }))
  ].sort((a: TimelineItem | TimelineEra, b: TimelineItem | TimelineEra) => {
    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
});

// Watch for timelineData changes to setup observer when DOM is ready
watch(timelineData, () => {
  nextTick(() => {
    setupIntersectionObserver();
  });
}, { immediate: true });

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
      <div v-if="user" class="mb-8 flex justify-end">
        <button @click="addConnection"
          class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors">
          {{ t(`connect.add_connection`) }}
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="isLoadingMore && !allEventsData.length" class="flex justify-center items-center h-64">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center text-red-600 p-4">
        <p>Error loading events: {{ error }}</p>
        <button @click="() => fetchEventsWithPagination(true)" class="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
          Retry
        </button>
      </div>

      <!-- Empty State -->
      <div v-else-if="!timelineData.length && !isLoadingMore" class="text-center text-gray-600 p-4">
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

        <!-- Load More Trigger -->
        <div ref="loadMoreTrigger" class="h-20 flex items-center justify-center">
          <div v-if="isLoadingMore" class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
          <div v-else-if="!hasMore && allEventsData.length > 0" class="text-gray-500 text-sm">
            {{ t('event.no_more_events') }}
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