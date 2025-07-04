<script setup lang="ts">
import { ref, computed } from 'vue';
import jwtDecode from 'jwt-decode';

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

const { t } = useI18n();

const router = useRouter();

// Get current user from JWT token
const currentUser = ref<{ id: string } | null>(null);
if (process.client) {
  const token = localStorage.getItem('auth_token');
  if (token) {
    try {
      const decoded = jwtDecode(token) as { id: string };
      currentUser.value = decoded;
    } catch (error) {
      console.error('Failed to decode token:', error);
    }
  }
}

type User = {
  id: string;
  username: string;
  email: string;
  created_at: string;
  password: string;
  connects: any[];
  role_id: number;
};

type ApiResponse = {
  data: User[];
  total: number;
  limit: number;
  skip: number;
};

const columns = [
  {
    key: "username",
    label: t('common.username'),
  },
  {
    key: "actions",
  },
];

const searchQuery = ref('');
const currentPage = ref(1);
const itemsPerPage = 10;

const { data: items, error, isFetching, execute: fetchUsers } =
  useApi(() => `/users/search?username=${searchQuery.value}&limit=${itemsPerPage}&skip=${(currentPage.value - 1) * itemsPerPage}`).
    json();

// Filter out current user from the response
const filteredItems = computed(() => {
  if (!items.value || !currentUser.value) return items.value;
  return {
    ...items.value,
    data: items.value.data.filter((user: User) => user.id !== currentUser.value?.id)
  };
});

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

const handleAddConnection = async (id: string) => {
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

// Debounced fetch function
const debouncedFetch = debounce(async () => {
  try {
    await fetchUsers();
  } catch (err: any) {
    console.error('Fetch error:', err);
    showToast(err.message || t('common.error_occurred'), 'error');
  }
}, 300); // 300ms delay

// Watch for search query changes
watch([searchQuery, currentPage], () => {
  debouncedFetch();
});

</script>

<template>
  <div v-if="error" class="p-4 bg-red-100 text-red-700 rounded-lg mb-4">
    {{ error.message || t('common.error_occurred') }}
  </div>
  <div v-else>
    <div class="flex flex-col gap-4">
      <UInput v-model="searchQuery" :placeholder="$t('common.search_users')" icon="i-heroicons-magnifying-glass"
        class="w-full" />

      <UTable :columns="columns" :rows="filteredItems?.data || []" :loading="isFetching" :loading-state="{
        icon: 'i-heroicons-arrow-path-20-solid',
        label: t('common.loading'),
      }">
        <template #username-data="{ row }">
          <UButton variant="link" @click="router.push(`/user/${row.id}`)">
            {{ row.username }}
          </UButton>
        </template>
        <template #actions-data="{ row }">
          <div class="flex justify-end">
            <UButton color="primary" size="sm" class="px-4 py-2 text-sm" @click="() => handleAddConnection(row.id)">
              {{ $t('connect.add_connection') }}
            </UButton>
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

    <div v-if="toast.show" :class="[
      'fixed top-6 left-1/2 transform -translate-x-1/2 px-6 py-3 rounded shadow-lg z-50 transition-all',
      toast.type === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white'
    ]">
      {{ toast.message }}
    </div>
  </div>
</template>
