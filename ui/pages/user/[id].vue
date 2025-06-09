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
}>("/eras?order=desc:time_start").json();

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
  ].sort((a: TimelineItem | TimelineEra, b: TimelineItem | TimelineEra ) => {

    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })
});

const addConnection = async () => {
  try {
    await useApi('/connects').post({
      friend_id: id
    });
    // Başarılı olduğunda kullanıcıya bilgi ver
    alert('Bağlantı başarıyla eklendi');
  } catch (error) {
    console.error('Bağlantı eklenirken hata oluştu:', error);
    alert('Bağlantı eklenirken bir hata oluştu');
  }
};
</script>

<template>
    <div class="min-h-screen bg-white p-8">
      <div class="max-w-6xl mx-auto">
        <!-- Bağlantı Ekle Butonu -->
        <div class="mb-8 flex justify-end">
          <button 
            @click="addConnection"
            class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
          >
            Bağlantı Ekle
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
            <div v-if="item.type === 'EVENT'" class="absolute left-1/2 -translate-x-1/2 w-4 h-4 rounded-full bg-gray-200 z-10"></div>

            <!-- Date Text -->
            <div v-if="item.type === 'EVENT'" class="absolute left-[52%] top-[-1.5rem] text-sm text-gray-600">
              {{ formatDate(item.date)}}
            </div>

            <!-- Cards Container -->
            <div v-if="item.type === 'EVENT'" class="grid grid-cols-2 gap-8">
                <div v-if="timelineData.indexOf(item) % 2 === 1"></div>
                <div  class="flex justify-end">
                  <div class="transform hover:rotate-0 transition-transform duration-200 -rotate-2" :class="timelineData.indexOf(item) % 2 === 0 ? '-rotate-2' : 'rotate-2'">
                    <EventItem  :id="item.id" :date="item.date" :cards="item.cards" :is-left="timelineData.indexOf(item) % 2 === 0" />
                  </div>
                </div>
                <div v-if="timelineData.indexOf(item) % 2 === 0"></div>
            </div>
            <div v-else class="shadow-lg w-full mb-4 px-4 z-20 text-center" :style="{ background: item.color }">
                <div>{{ formatDate(item.date) }}</div>
                <h3 class="font-bold text-6xl text-gray-800">{{ item.name }}</h3>
            </div>
            </div>
        </div>
      </div>
      </div>
    </div>
  </template>