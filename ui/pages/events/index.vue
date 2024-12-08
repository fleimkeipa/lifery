<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

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
    label: "Name",
  },
  {
    key: "description",
    label: "Description",
  },
  {
    key: "visibility",
    label: "Visibility",
  },
  {
    key: "date",
    label: "Date",
  },
  {
    key: "time_start",
    label: "TimeStart",
  },
  {
    key: "time_end",
    label: "TimeEnd",
  },
  {
    key: "items",
    label: "Items",
  },
];

const {
  data: items,
  error,
  isFetching,
  execute: fetchPods,
} = useApi<{ data: { items: Row[] } }>("/events").json();

const router = useRouter();

const actions = (row: Row) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/pods/${row.id}`),
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
    afterFetch: () => fetchPods(),
  }).delete();
};
</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-between">
      <UButton icon="i-heroicons-plus">
        <NuxtLink to="/events/create/new">Create New</NuxtLink>
      </UButton>
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchPods"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data.events" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: 'Loading...',
    }">
      <template #expand="{ row }">
        <div class="p-4">
          <b>Items:</b>
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
