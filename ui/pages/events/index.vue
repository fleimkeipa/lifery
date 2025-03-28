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
  time_start: Date;
  time_end: Date;
  items: [];
};

const columns = [
  {
    key: "id",
    label: "ID",
  },
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
    key: "time_start",
    label: t('common.time_start'),
  },
  {
    key: "time_end",
    label: t('common.time_end'),
  },
  {
    key: "items",
    label: t('event.items'),
  },
  {
    key: "actions",
  },
];

const {
  data: items,
  error,
  isFetching,
  execute: fetchEvents,
} = useApi<{ data: { items: Row[] } }>("/events").json();

const router = useRouter();

const actions = (row: Row) => [
  [
    {
      label: t('common.edit'),
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/events/${row.id}`),
    },
    {
      label: t('common.delete'),
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
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-between">
      <UButton icon="i-heroicons-plus">
        <NuxtLink to="/events/create/new">{{ t('common.create_new') }}</NuxtLink>
      </UButton>
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchEvents"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data.events" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: 'Loading...',
    }">
      <template #expand="{ row }">
        <div class="p-4">
          <b>{{ t('event.items') }}:</b>
          <pre>{{ row.items }}</pre>
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