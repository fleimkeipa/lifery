<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

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

interface TimelineItem {
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

const columns = [
  {
    key: "id",
    label: "ID",
  },
  {
    key: "name",
    label: "Name",
  },
  {
    key: "description",
    label: "Description",
  },
  {
    key: "visibility",
    label: "Status",
    formatter: (value: number) => {
      switch (value) {
        case 0:
          return "Private";
        case 1:
          return "Public";
        case 3:
          return "Featured";
        default:
          return "Unknown";
      }
    }
  },
  {
    key: "date",
    label: "Date",
    formatter: (value: Date) => {
      return new Date(value).toLocaleDateString('en-US', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    }
  },
  {
    key: "time_start",
    label: "Start Time",
    formatter: (value: Date) => {
      return new Date(value).toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  },
  {
    key: "time_end",
    label: "End Time",
    formatter: (value: Date) => {
      return new Date(value).toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  },
  {
    key: "items",
    label: "Media",
    formatter: (items: any[]) => {
      if (!items?.length) return 'No media';
      return `${items.length} items`;
    }
  },
  {
    key: "actions",
  },
];

const formatDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return date.toLocaleDateString('tr-TR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).replace(/\//g, '.');
};

const { data: eventsData, error, isFetching, execute: fetchEvents } = useApi<{
  data: {
    events: Row[];
    total: number;
    limit: number;
    skip: number;
  };
}>("/events?order=asc:date").json();

const timelineData = computed<TimelineItem[]>(() => {
  if (!eventsData.value?.data?.events) return [];

  return eventsData.value.data.events.map((event: Row) => ({
    id: event.id.toString(),
    date: formatDate(event.date.toString()),
    cards: [
      {
        title: event.name,
        description: event.description,
        color: event.visibility === 3 ? 'bg-pink-300' : 'bg-blue-200',
        items: event.items
      }
    ]
  }));
});

const router = useRouter();

const actions = (row: Row) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/events/${row.id}`),
    },
    {
      label: "Delete",
      icon: "i-heroicons-trash-20-solid",
      click: () => handleDelete(row.id),
    },
  ],
];

const handleDelete = async (uid: number) => {
  useApi(`/events/${uid}`, {
    afterFetch: () => fetchEvents(),
  }).delete();
};
</script>

<template>
  <div class="min-h-screen bg-white p-8">
    <div class="max-w-6xl mx-auto">
      <!-- Timeline Container -->
      <div class="relative">
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
              <!-- Left Side -->
              <div class="flex justify-end">
                <div v-if="timelineData.indexOf(item) % 2 === 0"
                  class="transform -rotate-2 hover:rotate-0 transition-transform duration-200">
                  <div v-for="card in item.cards" :key="card.title"
                    :class="[card.color, 'p-4 rounded shadow-lg max-w-[300px] aspect-[4/3] flex flex-col mb-4']">
                    <h3 class="font-medium text-gray-800 mb-2">{{ card.title }}</h3>
                    <p class="text-sm text-gray-600 overflow-auto">{{ card.description || 'No description available' }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- Right Side -->
              <div class="flex">
                <div v-if="timelineData.indexOf(item) % 2 === 1"
                  class="transform rotate-2 hover:rotate-0 transition-transform duration-200">
                  <div v-for="card in item.cards" :key="card.title"
                    :class="[card.color, 'p-4 rounded shadow-lg max-w-[300px] aspect-[4/3] flex flex-col mb-4']">
                    <h3 class="font-medium text-gray-800 mb-2">{{ card.title }}</h3>
                    <p class="text-sm text-gray-600 overflow-auto">{{ card.description || 'No description available' }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.timeline-item:not(:last-child)::after {
  content: '';
  position: absolute;
  left: 50%;
  top: 2rem;
  transform: translateX(-50%);
  width: 2px;
  height: calc(100% + 2rem);
  background-color: #e5e7eb;
}
</style>
