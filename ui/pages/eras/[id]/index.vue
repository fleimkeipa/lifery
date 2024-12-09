<script setup>
import * as yup from "yup";

definePageMeta({
  middleware: "auth",
});

const state = reactive({
  name: null,
  color: null,
  time_start: null,
  time_end: null,
});

const schema = yup.object({
  name: yup.string().nonNullable("Name cannot be null"),
  color: yup
    .string()
    .matches(/^#([0-9A-F]{3}){1,2}$/i, "Color must be a valid hex color"),
  time_start: yup.date().required("Start time cannot be empty"),
  time_end: yup.date().required("End time cannot be empty"),
});

const route = useRoute();
const { isFetching } = useApi(`/eras/${route.params.id}`, {
  afterFetch: (ctx) => {
    const era = ctx.data.data;
    state.name = era.name;
    state.color = era.color;
    state.time_start = era.time_start;
    state.time_end = era.time_end;
  },
}).json();

const router = useRouter();
const loading = ref(false);
const error = ref(null);
const onSubmit = (era) => {
  loading.value = true;
  useApi(`/eras/${route.params.id}`, {
    afterFetch: () => {
      loading.value = false;
      router.push("/eras");
    },
    onFetchError: ({ error: fetchErr }) => {
      loading.value = false;
      error.value = fetchErr;
    },
  }).patch(era.data);
};
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">Update Era</h1>
    <div v-if="isFetching">Loading...</div>
    <UForm
      @submit="onSubmit"
      style="display: flex; flex-direction: column; gap: 20px"
      novalidate
      :state="state"
      :schema="schema"
      class="mt-8 flex items-start"
    >
      <UFormGroup label="Name" name="name">
        <UInput type="text" placeholder="Name" v-model="state.name" />
      </UFormGroup>

      <UFormGroup label="Color" name="color">
        <UInput type="color" placeholder="Color" v-model="state.color" />
      </UFormGroup>

      <UFormGroup label="TimeStart" name="time_start">
        <UInput
          type="date"
          pattern="\d{4}-\d{2}-\d{2}"
          placeholder="Time Start"
          v-model="state.time_start"
        />
      </UFormGroup>
      <UFormGroup label="TimeEnd" name="time_end">
        <UInput
          type="date"
          pattern="\d{4}-\d{2}-\d{2}"
          placeholder="Time End"
          v-model="state.time_end"
        />
      </UFormGroup>

      <UButton :loading="loading" type="submit">Submit</UButton>
      <div v-if="error" class="flex items-center gap-x-2 rounded-lg border px-2 py-1">
        <span @click="error = null" class="cursor-pointer">X</span>
        <span class="text-sm text-red-500">{{ error }}</span>
      </div>
    </UForm>
  </div>
</template>
