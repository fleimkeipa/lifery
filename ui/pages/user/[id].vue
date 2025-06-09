<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

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
}

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString('tr-TR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).replace(/\//g, '.');
};

const { data: eventsData, error, isFetching, execute: fetchEvents } = useApi<{
  data: Row[];
  total: number;
  limit: number;
  skip: number;
}>(`/events?user_id=${id}`).json()

const timelineData = computed<TimelineItem[]>(() => {
  if (!eventsData.value?.data) return [];

  return eventsData.value.data.map((event: Row) => ({
    id: event.id.toString(),
    date: formatDate(event.date.toString()),
    cards: [
      {
        title: event.name,
        description: event.description,
        color: event.visibility === 3 ? 'bg-pink-300' : event.visibility === 2 ? 'bg-yellow-200' : 'bg-blue-200',
        items: event.items
      }
    ]
  }));
});
</script>

<template>
    <div class="min-h-screen bg-white p-8">
      <div class="max-w-6xl mx-auto">
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
          <p>No events found</p>
        </div>
  
        <!-- Timeline Container -->
        <div v-else class="relative">
          <!-- Vertical Line -->
          <div class="absolute left-1/2 h-full w-0.5 bg-gray-200"></div>
  
          <!-- Timeline Items -->
          <div class="space-y-24">
            <div v-for="item in timelineData" :key="item.id" class="relative">
              <!-- Date Marker -->
              <div class="absolute left-1/2 -translate-x-1/2 w-4 h-4 rounded-full bg-gray-200 z-10"></div>
  
              <!-- Date Text -->
              <div class="absolute left-[52%] top-[-1.5rem] text-sm text-gray-600">
                {{ item.date }}
              </div>
  
              <!-- Cards Container -->
              <div class="grid grid-cols-2 gap-8">
                  <div v-if="timelineData.indexOf(item) % 2 === 1"></div>
                  <div class="flex justify-end">
                    <div class="transform hover:rotate-0 transition-transform duration-200 -rotate-2" :class="timelineData.indexOf(item) % 2 === 0 ? '-rotate-2' : 'rotate-2'">
                      <EventItem :id="item.id" :date="item.date" :cards="item.cards" :is-left="timelineData.indexOf(item) % 2 === 0" />
                    </div>
                  </div>
                  <div v-if="timelineData.indexOf(item) % 2 === 0"></div>
                </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>