<script setup lang="ts">
import { useAuth } from '../../composables/useAuth';

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

const toggleSort = (column: string) => {
  if (sortBy.value === column) {
    sortOrder.value = sortOrder.value === 'desc' ? 'asc' : 'desc';
  } else {
    sortBy.value = column;
    sortOrder.value = 'desc';
  }
};

watch([currentPage, sortOrder, sortBy], () => {
  fetchConnects();
});

const { data: items, error, isFetching, execute: fetchConnects } = useApi(() => `/connects?order=${sortOrder.value}:${sortBy.value}&limit=${itemsPerPage}&skip=${(currentPage.value - 1) * itemsPerPage}`).json<{
  data: Row[];
  total: number;
}>();

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
  useApi(`/connects/${uid}`, {
    afterFetch: () => fetchConnects(),
  }).delete();
};

const handleAccept = async (uid: number) => {
  useApi(`/connects/${uid}`, {
    afterFetch: () => fetchConnects(),
  }).patch({ status: 101});
};

</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-end">
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchConnects"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
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
        <UDropdown :items="actions(row)">
          <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
        </UDropdown>
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
