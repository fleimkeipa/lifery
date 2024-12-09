<script setup lang="ts">
// definePageMeta({
//   middleware: "auth",
// });

type Row = {
  id: number;
  name: string;
  color: { hex: string; rgb: string };
  time_start: Date;
  time_end: Date;
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
    key: "color",
    label: "Color",
  },
  {
    key: "time_start",
    label: "Time Start",
  },
  {
    key: "time_end",
    label: "Time End",
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
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/eras/${row.id}`),
    },
    {
      label: "Delete",
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
        <NuxtLink to="/eras/create/new">Create New</NuxtLink>
      </UButton>
      <UButton
        icon="i-heroicons-arrow-path"
        :loading="isFetching"
        @click="fetchEras"
      ></UButton>
    </div>
    <UTable
      :columns="columns"
      :rows="items.data.eras"
      :loading="isFetching"
      :loading-state="{
        icon: 'i-heroicons-arrow-path-20-solid',
        label: 'Loading...',
      }"
    >
      <template #actions-data="{ row }">
        <UDropdown :items="actions(row)">
          <UButton
            color="gray"
            variant="ghost"
            icon="i-heroicons-ellipsis-horizontal-20-solid"
          />
        </UDropdown>
      </template>
    </UTable>
  </div>
</template>
