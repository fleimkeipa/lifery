<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

const { t, locale } = useI18n();

type Row = {
  id: number;
  name: string;
  description: string;
  visibility: number;
  date: Date;
  items: Item[];
};

type Item = {
  type: number;
  data: string;
};

const visibilityOptions = [
  { label: t('event.visibilityOpts.public'), value: 1 },
  { label: t('event.visibilityOpts.private'), value: 2 },
  { label: t('event.visibilityOpts.just_me'), value: 3 }
];

const columns = [
  {
    key: "name",
    label: t('common.name'),
  },
  {
    key: "description",
    label: t('common.description'),
  },
  {
    key: "visibility",
    label: t('event.visibility'),
  },
  {
    key: "date",
    label: t('common.date'),
  },
  {
    key: "actions",
  },
];

const selectedRows = ref<Row[]>([]);
const currentPage = ref(1);
const itemsPerPage = 30;
const sortOrder = ref('desc');
const sortBy = ref('date');
const searchQuery = ref('');

const toggleSort = (column: string) => {
  if (sortBy.value === column) {
    sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc';
  } else {
    sortBy.value = column;
    sortOrder.value = 'desc';
  }
};

watch([currentPage, sortOrder, sortBy, searchQuery], () => {
  fetchEvents();
});

const getQueryParams = () => {
  const params = new URLSearchParams({
    limit: itemsPerPage.toString(),
    skip: ((currentPage.value - 1) * itemsPerPage).toString(),
    order: `${sortOrder.value}:${sortBy.value}`
  });

  if (searchQuery.value) {
    params.append('name', `like:${searchQuery.value}`);
  }

  return params.toString();
};

const {
  data: items,
  error,
  isFetching,
  execute: fetchEvents,
} = useApi(() => `/events?${getQueryParams()}`).json();

const router = useRouter();


const handleDelete = async (uid: number) => {
  useApi(`/events/${uid}`, {
    afterFetch: () => fetchEvents(),
  }).delete();
};
</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-between mb-4">
      <UButton icon="i-heroicons-plus">
        <NuxtLink to="/events/create/new">{{ t('common.create_new') }}</NuxtLink>
      </UButton>
      <div class="flex items-center gap-4">
        <UInput
          v-model="searchQuery"
          :placeholder="t('common.search')"
          icon="i-heroicons-magnifying-glass"
          class="w-64"
        />
        <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchEvents"></UButton>
      </div>
    </div>
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }" row-selectable v-model:selected="selectedRows" :row-expandable="() => true" show-detail-on-click>
      <template #name-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('name')">
          <span>{{ t('common.name') }}</span>
          <UIcon v-if="sortBy === 'name'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #date-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('date')">
          <span>{{ t('common.date') }}</span>
          <UIcon v-if="sortBy === 'date'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #visibility-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('visibility')">
          <span>{{ t('event.visibility') }}</span>
          <UIcon v-if="sortBy === 'visibility'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #visibility-data="{ row }">
        <span :style="{
          display: 'inline-block',
          width: '10px',
          height: '10px',
          'border-radius': '50%',
          'margin-right': '8px',
          'background':
            row.visibility === 1 ? '#bfdbfe' : // bg-blue-200
              row.visibility === 2 ? '#fef08a' : // bg-yellow-200
                row.visibility === 3 ? '#f9a8d4' : // bg-pink-300
                  '#d1d5db' // gri (default)
        }"></span>
        {{visibilityOptions.find(option => option.value === row.visibility)?.label || '-'}}
      </template>

      <template #name-data="{ row }">
        <UButton variant="link" @click="router.push(`/events/${row.id}`)">
          {{ row.name }}
        </UButton>
      </template>

      <template #row-details="{ row }">
        <div class="p-4">
          <b>{{ t('event.items') }}:</b>
          <div v-if="row.items && row.items.length > 0" class="mt-4">
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              <div v-for="(item, index) in row.items.filter((item: Item) => item.type === 11)" :key="index"
                class="relative aspect-square">
                <img :src="item.data" :alt="`Photo ${index + 1}`" class="w-full h-full object-cover rounded-lg" />
              </div>
            </div>
          </div>
          <div v-else class="mt-4 text-gray-500">
            {{ t('event.no_events') }}
          </div>
        </div>
      </template>

      <template #actions-data="{ row }">
        <div class="flex justify-end">
          <UButton 
            color="red" 
            variant="ghost" 
            icon="i-heroicons-trash-20-solid"
            @click="handleDelete(row.id)"
            :title="t('common.delete')"
          />
        </div>
      </template>
    </UTable>

    <div class="flex justify-end mt-4">
      <UPagination v-model="currentPage" :total="items?.total || 0" :page-count="itemsPerPage" :ui="{
        wrapper: 'flex items-center justify-end',
        base: 'flex items-center gap-1',
        rounded: 'rounded-md',
        default: {
          size: 'sm',
          activeButton: {
            color: 'primary'
          }
        }
      }" />
    </div>
  </div>
</template>