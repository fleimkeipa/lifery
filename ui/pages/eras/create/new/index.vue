<script setup>
import * as yup from "yup";

definePageMeta({
  middleware: "auth",
});

const { t, locale } = useI18n();

const state = reactive({
  name: null,
  color: null,
  time_start: null,
  time_end: null,
});

const formatDateTime = (dateTime) => {
  if (!dateTime) return null;
  const date = new Date(dateTime);
  return date.toISOString();
};

const formattedTimeStart = computed(() => formatDateTime(state.time_start));
const formattedTimeEnd = computed(() => formatDateTime(state.time_end));

const schema = yup.object({
  name: yup.string().nonNullable(t('common.name')),
  color: yup
    .string()
    .matches(/^#([0-9A-F]{3}){1,2}$/i, t('era.validation.required.color')),
  time_start: yup.date().required("Start time cannot be empty"),
  time_end: yup.date().required("End time cannot be empty"),
});

const router = useRouter();
const loading = ref(false);
const error = ref(null);
const onSubmit = (era) => {
  loading.value = true;

  const formData = {
    ...era.data,
    time_start: formattedTimeStart.value,
    time_end: formattedTimeEnd.value
  };

  useApi("/eras", {
    afterFetch: () => {
      loading.value = false;
      router.push("/eras");
    },
    onFetchError: ({ error: fetchErr }) => {
      loading.value = false;
      error.value = fetchErr;
    },
  }).post(formData);
};
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">{{ t('common.create_new') }}</h1>
    <UForm @submit="onSubmit" style="display: flex; flex-direction: column; gap: 20px" novalidate :state="state"
      :schema="schema" class="mt-8 flex items-start">
      <UFormGroup :label="t('common.name')" name="name">
        <UInput type="text" :placeholder="t('common.write_name')" v-model="state.name" />
      </UFormGroup>

      <UFormGroup :label="t('common.color')" name="color">
        <UInput type="color" placeholder="Color" v-model="state.color" />
      </UFormGroup>

      <UFormGroup :label="t('common.time_start')" name="time_start">
        <UInput type="datetime-local" :placeholder="t(`common.time_start`)" v-model="state.time_start" />
      </UFormGroup>
      <UFormGroup :label="t('common.time_end')" name="time_end">
        <UInput type="datetime-local" :placeholder="t(`common.time_end`)" v-model="state.time_end" />
      </UFormGroup>

      <UButton :loading="loading" type="submit">{{ t('common.create_new') }}</UButton>
      <div v-if="error" class="flex items-center gap-x-2 rounded-lg border px-2 py-1">
        <span @click="error = null" class="cursor-pointer">X</span>
        <span class="text-sm text-red-500">{{ error }}</span>
      </div>
    </UForm>
  </div>
</template>
