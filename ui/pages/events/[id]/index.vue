<script setup>
import * as yup from "yup";

definePageMeta({
  middleware: "auth",
});

const { t, locale } = useI18n();

const state = reactive({
  name: null,
  description: null,
  visibility: null,
  date: null,
  time_start: null,
  time_end: null,
  items: [],
});

const schema = yup.object({
  name: yup.string().nonNullable(t('common.name')),
  description: yup.string(),
  visibility: yup
    .number()
    .oneOf([1, 2, 3], t('event.validation.one_of.visibility'))
    .required(t('event.validation.required.visibility')),
  date: yup.date().nullable(),
  time_start: yup.date().nullable(),
  time_end: yup.date().nullable(),
  items: yup.array().of(
    yup.object({
      type: yup
        .number()
        .oneOf([10, 11, 12, 13], t('event.validation.one_of.type'))
        .required(t('event.validation.required.type')),
      data: yup.string().required(t('event.validation.required.data')),
    })
  ),
});

const route = useRoute();
const { isFetching } = useApi(`/events/${route.params.id}`, {
  afterFetch: (ctx) => {
    const event = ctx.data.data;
    state.name = event.name;
    state.description = event.description;
    state.visibility = event.visibility;
    state.date = event.date;
    state.time_start = event.time_start;
    state.time_end = event.time_end;
    state.items = event.items;
  },
}).json();

const router = useRouter();
const loading = ref(false);
const error = ref(null);
const onSubmit = (event) => {
  loading.value = true;
  useApi(`/events/${route.params.id}`, {
    afterFetch: () => {
      loading.value = false;
      router.push("/events");
    },
    onFetchError: ({ error: fetchErr }) => {
      loading.value = false;
      error.value = fetchErr;
    },
  }).patch(event.data);
};

const push = () => {
  state.items.push({ type: null, data: null });
};
const remove = (idx) => {
  state.items.splice(idx, 1);
};

const visibilityOptions = [
  { label: t('event.visibilityOpts.public'), value: 1 },
  { label: t('event.visibilityOpts.private'), value: 2 },
  { label: t('event.visibilityOpts.just_me'), value: 3 }
];

const typeOptions = [
  { label: t('event.typeOpts.text'), value: 10 },
  { label: t('event.typeOpts.photo'), value: 11 },
  { label: t('event.typeOpts.video'), value: 12 },
  { label: t('event.typeOpts.voice_record'), value: 13 }
];

</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">{{ t('common.update') }}</h1>
    <div v-if="isFetching">{{ t('common.loading') }}</div>
    <UForm @submit="onSubmit" style="display: flex; flex-direction: column; gap: 20px" novalidate :state="state"
      :schema="schema" class="mt-8 flex items-start">
      <UFormGroup :label="t('common.name')" name="name">
        <UInput type="text" :placeholder="t('common.write_name')" v-model="state.name" />
      </UFormGroup>
      <UFormGroup :label="t('common.description')" name="description">
        <UInput type="text" :placeholder="t(`common.write_desc`)" v-model="state.description" />
      </UFormGroup>
      <UFormGroup :label="t('event.visibility')" name="visibility">
        <USelect v-model="state.visibility" :placeholder="t(`event.select_visibility`)" :options="visibilityOptions"
          @update:modelValue="(val) => state.visibility = Number(val)" />
      </UFormGroup>
      <UFormGroup :label="t('common.date')" name="date">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" :placeholder="t(`common.date`)" v-model="state.date" />
      </UFormGroup>
      <UFormGroup :label="t('common.time_start')" name="time_start">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" placeholder="Time Start" v-model="state.time_start" />
      </UFormGroup>
      <UFormGroup :label="t('common.time_end')" name="time_end">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" placeholder="Time End" v-model="state.time_end" />
      </UFormGroup>

      <div>
        <div class="mb-4 flex flex-row items-center gap-x-4">
          <h1 class="text-xl">{{ t(`event.items`) }}</h1>
          <UButton @click="push()" size="sm" :ui="{ rounded: 'rounded-full' }" color="blue" icon="i-heroicons-plus">
            {{ t(`event.add_item`) }}
          </UButton>
        </div>
        <div class="mb-4 flex flex-col gap-y-4" v-if="state.items.length">
          <div :key="item.key" v-for="(item, idx) in state.items"
            class="flex flex-row items-start justify-center gap-x-6">
            <UFormGroup :label="t('event.item')" :type="`items[${idx}].type`">
              <USelect v-model="item.type" :placeholder="t('event.select_item')" :options="typeOptions"
                @update:modelValue="(val) => item.type = Number(val)" />
            </UFormGroup>

            <UFormGroup :label="t('event.data')" :name="`items[${idx}].data`">
              <UInput type="text" :placeholder="t('event.data_of_item')" v-model="item.data" />
            </UFormGroup>

            <UButton @click="remove(idx)" size="sm" :ui="{ rounded: 'rounded-full' }" color="red"
              icon="i-heroicons-trash" class="self-center">{{ t(`event.remove_item`) }}
            </UButton>
          </div>
        </div>
        <div v-else>
          <span class="text-md font-bold">{{ t(`event.no_item`) }}</span>
        </div>
      </div>

      <UButton :loading="loading" type="submit">{{ t(`event.add_item`) }}</UButton>
      <div v-if="error" class="flex items-center gap-x-2 rounded-lg border px-2 py-1">
        <span @click="error = null" class="cursor-pointer">X</span>
        <span class="text-sm text-red-500">{{ error }}</span>
      </div>
    </UForm>
  </div>
</template>
