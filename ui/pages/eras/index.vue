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

const { data: items, error, isFetching, execute: fetchEras } = useApi<{
  data: { eras: Row[] };
}>("/eras").json();

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
    <div class="flex flex-row items-center justify-between">
      <UButton icon="i-heroicons-plus">
        <NuxtLink to="/eras/create/new">{{ t('common.create_new') }}</NuxtLink>
      </UButton>
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchEras"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }">
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
  </div>
</template>
