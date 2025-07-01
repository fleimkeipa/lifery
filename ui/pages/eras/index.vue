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
  color: { hex: string; rgb: string };
  time_start: Date;
  time_end: Date;
};

const columns = [
  {
    key: "name",
    label: t('common.name'),
  },
  {
    key: "color",
    label: t('common.color'),
  },
  {
    key: "time_start",
    label: t('common.time_start'),
  },
  {
    key: "time_end",
    label: t('common.time_end'),
  },
  {
    key: "actions",
  },
];

const currentPage = ref(1);
const itemsPerPage = 30;
const sortOrder = ref('desc');
const sortBy = ref('time_start');
const searchQuery = ref('');

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

const toggleSort = (column: string) => {
  if (sortBy.value === column) {
    sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc';
  } else {
    sortBy.value = column;
    sortOrder.value = 'desc';
  }
};

const { data: items, error, isFetching, execute: fetchEras } = useApi(() => `/eras?${getQueryParams()}`).json<{
  data: Row[];
  total: number;
}>();

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
    await fetchEras();
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

const handleDelete = async (id: number) => {
  try {
    await useApi(`/eras/${id}`).delete();
    showToast(t('common.deleted_successfully'), 'success');
    await fetchEras();
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
        <NuxtLink to="/eras/create/new">{{ t('common.create_new') }}</NuxtLink>
      </UButton>
      <div class="flex items-center gap-4">
        <UInput
          v-model="searchQuery"
          :placeholder="t('common.search')"
          icon="i-heroicons-magnifying-glass"
          class="w-64"
        />
        <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchEras"></UButton>
      </div>
    </div>
    <UTable :columns="columns" :rows="items?.data || []" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }">
      <template #name-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('name')">
          <span>{{ t('common.name') }}</span>
          <UIcon v-if="sortBy === 'name'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #time_start-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('time_start')">
          <span>{{ t('common.time_start') }}</span>
          <UIcon v-if="sortBy === 'time_start'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #time_end-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('time_end')">
          <span>{{ t('common.time_end') }}</span>
          <UIcon v-if="sortBy === 'time_end'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #name-data="{ row }">
        <UButton variant="link" @click="router.push(`/eras/${row.id}`)">
          {{ row.name }}
        </UButton>
      </template>
      <template #color-data="{ row }">
        <div class="flex items-center gap-2">
          <div class="w-6 h-6 rounded-full" :style="{ backgroundColor: row.color }"></div>
          <span>{{ row.color }}</span>
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
