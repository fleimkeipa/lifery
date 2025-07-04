<script setup lang="ts">
import { useAuth } from '../../composables/useAuth';

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
const { user } = useAuth();

type Row = {
  id: number;
  user_id: string;
  friend_id: string;
  status: number;
  user: {
    username: string;
  };
  friend: {
    username: string;
  };
};

const statusOptions = [
  { label: t('connect.statusOpts.pending'), value: 100 },
  { label: t('connect.statusOpts.accepted'), value: 101 },
  { label: t('connect.statusOpts.rejected'), value: 102 }
];

const columns = [
  // {
  //   key: "id",
  //   label: "ID",
  // },
  // {
  //   key: "user.username",
  //   label: t('common.user'),
  // },
  {
    key: "friend.username",
    label: t('connect.friend_username'), 
  },
  {
    key: "status",
    label: t('connect.status'),
  },
  {
    key: "actions",
  },
];

const currentPage = ref(1);
const itemsPerPage = 30;
const sortOrder = ref('asc');
const sortBy = ref('status');
const searchQuery = ref('');

const getQueryParams = () => {
  const params = new URLSearchParams({
    limit: itemsPerPage.toString(),
    skip: ((currentPage.value - 1) * itemsPerPage).toString(),
    order: `${sortOrder.value}:${sortBy.value}`
  });

  if (searchQuery.value) {
    params.append('username', `${searchQuery.value}`);
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

const { data: items, error, isFetching, execute: fetchConnects } = useApi(() => `/connects?${getQueryParams()}`).json<{
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
    await fetchConnects();
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

const actions = (row: Row) => [
  [
    ...(row.status === 100 && row.user_id != user.value?.id ? [{
      label: t('common.accept'),
      icon: "i-heroicons-check-circle-20-solid",
      click: () => handleAccept(row.id),
    }] : []),
    {
      label: t('common.delete'),
      icon: "i-heroicons-trash-20-solid",
      click: () => handleDelete(row.id),
    },
  ],
];

const handleDelete = async (uid: number) => {
  try {
    await useApi(`/connects/${uid}`).delete();
    showToast(t('common.deleted_successfully'), 'success');
    await fetchConnects();
  } catch (error: any) {
    console.error('Delete error:', error);
    showToast(error.message || t('common.error_occurred'), 'error');
  }
};

const handleAccept = async (uid: number) => {
  try {
    await useApi(`/connects/${uid}`).patch({ status: 101});
    showToast(t('connect.connection_accepted'), 'success');
    await fetchConnects();
  } catch (error: any) {
    console.error('Accept error:', error);
    showToast(error.message || t('common.error_occurred'), 'error');
  }
};
</script>

<template>
  <div v-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg mb-4">
    {{ error.message || t('common.error_occurred') }}
  </div>
  <div v-else>
    <div class="flex flex-row items-center justify-end mb-4">
      <div class="flex items-center gap-4">
        <UInput
          v-model="searchQuery"
          :placeholder="t('common.search')"
          icon="i-heroicons-magnifying-glass"
          class="w-64"
        />
        <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchConnects"></UButton>
      </div>
    </div>
    <UTable :columns="columns" :rows="items?.data || []" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }">
      <template #friend.username-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('friend.username')">
          <span>{{ t('connect.friend_username') }}</span>
          <UIcon v-if="sortBy === 'friend.username'" :name="sortOrder === 'desc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #status-header>
        <div class="flex items-center gap-2 cursor-pointer" @click="toggleSort('status')">
          <span>{{ t('connect.status') }}</span>
          <UIcon v-if="sortBy === 'status'" :name="sortOrder === 'asc' ? 'i-heroicons-arrow-down' : 'i-heroicons-arrow-up'" class="w-4 h-4" />
        </div>
      </template>

      <template #status-data="{ row }">
        <span
          :style="{
            display: 'inline-block',
            width: '10px',
            height: '10px',
            'border-radius': '50%',
            'margin-right': '8px',
            'background':
              row.status === 100 ? '#fef08a' : // yellow
              row.status === 101 ? '#86efac' : // green
              row.status === 102 ? '#fca5a5' : // red
              '#d1d5db'
          }"
        ></span>
        {{ statusOptions.find(option => option.value === row.status)?.label || '-' }}
      </template>

      <template #friend.username-data="{ row }">
        <UButton
          variant="link"
          @click="router.push(`/user/${row.friend.id}`)"
        >
          {{ row.friend.username }}
        </UButton>
      </template>

      <template #actions-data="{ row }">
        <div class="flex gap-2 justify-end">
          <UButton 
            v-if="row.status === 100 && row.user_id != user?.id"
            color="green" 
            variant="ghost" 
            icon="i-heroicons-check-circle-20-solid"
            @click="handleAccept(row.id)"
            :title="t('common.accept')"
          />
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
