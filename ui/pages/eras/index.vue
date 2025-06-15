<script setup lang="ts">
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

watch([currentPage, sortOrder, sortBy, searchQuery], () => {
  fetchEras();
});

const { data: items, error, isFetching, execute: fetchEras } = useApi(() => `/eras?${getQueryParams()}`).json<{
  data: Row[];
  total: number;
}>();

const router = useRouter();

const actions = (row: Row) => [
  [
    {
      label: t('common.edit'),
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/eras/${row.id}`),
    },
    {
      label: t('common.delete'),
      icon: "i-heroicons-trash-20-solid",
      click: () => handleDelete(row.id),
    },
  ],
];

const handleDelete = async (id: number) => {
  useApi(`/eras/${id}`, {
    afterFetch: () => fetchEras(),
  }).delete();
};
</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
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
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
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
