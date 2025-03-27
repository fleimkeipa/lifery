<script setup lang="ts">

definePageMeta({
  middleware: "auth",
});

type Row = {
  id: number;
  user_id: string;
  friend_id: string;
  status: number;
};

const columns = [
  {
    key: "id",
    label: "ID",
  },
  {
    key: "user_id",
    label: "User ID",
  },
  {
    key: "friend_id",
    label: "Friend ID",
  },
  {
    key: "status",
    label: "Status",
  },
  {
    key: "actions",
  },
];

const { data: items, error, isFetching, execute: fetchConnects } = useApi<{
  data: { items: Row[] };
}>("/connects").json();

const router = useRouter();

const actions = (row: Row) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => router.push(`/connects/${row.id}`),
    },
    {
      label: "Delete",
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
</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-between">
      <UButton icon="i-heroicons-plus">
        <NuxtLink to="/connects/create/new">Create New</NuxtLink>
      </UButton>
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchConnects"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data.connects" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: 'Loading...',
    }">
      <template #actions-data="{ row }">
        <UDropdown :items="actions(row)">
          <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
        </UDropdown>
      </template>
    </UTable>
  </div>
</template>
