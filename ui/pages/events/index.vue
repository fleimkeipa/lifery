<script setup lang="ts">
// Debounce utility function
function debounce<T extends (...args: any[]) => any>(func: T, wait: number): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

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

// Watch for errors
watch(error, (newError) => {
  if (newError) {
    console.error('API Error:', newError);
    showToast(newError.message || t('common.error_occurred'), 'error');
  }
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

// Debounced fetch function
const debouncedFetch = debounce(async () => {
  try {
    await fetchEvents();
  } catch (err: any) {
    console.error('Fetch error:', err);
    showToast(err.message || t('common.error_occurred'), 'error');
  }
}, 300); // 300ms delay

// Watch for search query changes
watch([searchQuery, currentPage, sortOrder, sortBy], () => {
  debouncedFetch();
});

const router = useRouter();

const handleDelete = async (uid: number) => {
  try {
    await useApi(`/events/${uid}`).delete();
    showToast(t('common.deleted_successfully'), 'success');
    await fetchEvents();
  } catch (error: any) {
    console.error('Delete error:', error);
    showToast(error.message || t('common.error_occurred'), 'error');
  }
};
</script>

<template>
  <div v-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg mb-4">
    {{ error.message || t('common.error_occurred') }}
  </div>
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
    <UTable :columns="columns" :rows="items?.data || []" :loading="isFetching" :loading-state="{
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
        <UButton variant="link" @click="router.push(`/events/${row.id}`)" :title="row.name">
          {{ row.name.length > 30 ? row.name.substring(0, 30) + '...' : row.name }}
        </UButton>
      </template>

      <template #description-data="{ row }">
        <span :title="row.description">
          {{ row.description && row.description.length > 30 ? row.description.substring(0, 30) + '...' : row.description || '-' }}
        </span>
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

    <div v-if="toast.show" :class="[
      'fixed top-6 left-1/2 transform -translate-x-1/2 px-6 py-3 rounded shadow-lg z-50 transition-all',
      toast.type === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'
    ]">
      {{ toast.message }}
    </div>
  </div>
</template>