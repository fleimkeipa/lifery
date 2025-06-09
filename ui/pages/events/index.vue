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
  items: Item[];
};

type Item = {
  type: number;
  data: string;
};

const visibilityOptions = [
  { label: t('event.visibilityOpts.public'), value: 1 },
  { label: t('event.visibilityOpts.private'), value: 2 },
  { label: t('event.visibilityOpts.just_me'), value: 3 }
];

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

const selectedRows = ref<Row[]>([]);

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
    <UTable :columns="columns" :rows="items.data" :loading="isFetching" :loading-state="{
      icon: 'i-heroicons-arrow-path-20-solid',
      label: t('common.loading'),
    }" row-selectable v-model:selected="selectedRows" :row-expandable="() => true" show-detail-on-click>
      <template #visibility-data="{ row }">
        {{visibilityOptions.find(option => option.value === row.visibility)?.label || '-'}}
      </template>

      <template #row-details="{ row }">
        <div class="p-4">
          <b>{{ t('event.items') }}:</b>
          <div v-if="row.items && row.items.length > 0" class="mt-4">
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              <div v-for="(item, index) in row.items.filter((item: Item) => item.type === 11)" :key="index"
                class="relative aspect-square">
                <img :src="item.data" :alt="`Photo ${index + 1}`" class="w-full h-full object-cover rounded-lg" />
              </div>
            </div>
          </div>
          <div v-else class="mt-4 text-gray-500">
            {{ t('common.no_items') }}
          </div>
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