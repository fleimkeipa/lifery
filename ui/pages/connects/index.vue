<script setup lang="ts">

definePageMeta({
  middleware: "auth",
});

const { t, locale } = useI18n();

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

const { data: items, error, isFetching, execute: fetchConnects } = useApi<{
  data: { items: Row[] };
}>("/connects").json();

const router = useRouter();

const actions = (row: Row) => [
  [
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
</script>

<template>
  <div v-if="!!error || !items">{{ error }}</div>
  <div v-else>
    <div class="flex flex-row items-center justify-between">
      <UButton icon="i-heroicons-arrow-path" :loading="isFetching" @click="fetchConnects"></UButton>
    </div>
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }">
      <template #status-data="{ row }">
        {{ statusOptions.find(option => option.value === row.status)?.label || '-' }}
      </template>

      <template #actions-data="{ row }">
        <UDropdown :items="actions(row)">
          <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
        </UDropdown>
      </template>
    </UTable>
  </div>
</template>
