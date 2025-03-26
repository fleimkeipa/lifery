<script setup>
import * as yup from "yup";

definePageMeta({
  middleware: "auth",
});

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
  name: yup.string().nonNullable("Name cannot be null"),
  description: yup.string().nullable(),
  visibility: yup
    .number()
    .oneOf([1, 2, 3], "Visibility must be Public (1), Private (2), or JustMe (3)")
    .required("Visibility is required"),
  date: yup.date().nullable(),
  time_start: yup.date().nullable(),
  time_end: yup.date().nullable(),
  items: yup.array().of(
    yup.object({
      type: yup
        .number()
        .oneOf([10, 11, 12, 13], "Type must be Text(10), Photo(11), Video(12), or Voice Record(13)"),
      // .required("Type is required"),
      data: yup.string().required("Item data cannot be empty"),
    })
  ),
});

const push = () => {
  state.items.push({ type: null, data: null });
};
const remove = (idx) => {
  state.items.splice(idx, 1);
};

const visibilityOptions = [
  { label: "Public", value: 1 },
  { label: "Private", value: 2 },
  { label: "Just Me", value: 3 }
];

const typeOptions = [
  { label: "Text", value: 10 },
  { label: "Photo", value: 11 },
  { label: "Video", value: 12 },
  { label: "Voice Record", value: 13 }
];

const router = useRouter();
const loading = ref(false);
const error = ref(null);
const onSubmit = (event) => {
  loading.value = true;
  useApi("/events", {
    afterFetch: () => {
      loading.value = false;
      router.push("/events");
    },
    onFetchError: ({ error: fetchErr }) => {
      loading.value = false;
      error.value = fetchErr;
    },
  }).post(event.data);
};
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold">Create Event</h1>
    <UForm @submit="onSubmit" style="display: flex; 
    flex-direction: column; 
    gap: 20px" novalidate :state="state" :schema="schema" class="mt-8 flex items-start">
      <UFormGroup label="Name" name="name">
        <UInput type="text" placeholder="Name" v-model="state.name" />
      </UFormGroup>
      <UFormGroup label="Description" name="description">
        <UInput type="text" placeholder="Description" v-model="state.description" />
      </UFormGroup>
      <UFormGroup label="Visibility" name="visibility">
        <USelect v-model="state.visibility" placeholder="Select Visibility" :options="visibilityOptions"
          @update:modelValue="(val) => state.visibility = Number(val)" />
      </UFormGroup>
      <UFormGroup label="Date" name="date">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" placeholder="Date" v-model="state.date" />
      </UFormGroup>
      <UFormGroup label="TimeStart" name="time_start">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" placeholder="Time Start" v-model="state.time_start" />
      </UFormGroup>
      <UFormGroup label="TimeEnd" name="time_end">
        <UInput type="date" pattern="\d{4}-\d{2}-\d{2}" placeholder="Time End" v-model="state.time_end" />
      </UFormGroup>

      <div>
        <div class="mb-4 flex flex-row items-center gap-x-4">
          <h1 class="text-xl">Items</h1>
          <UButton @click="push()" size="sm" :ui="{ rounded: 'rounded-full' }" color="blue" icon="i-heroicons-plus">
            Add Item
          </UButton>
        </div>
        <div class="mb-4 flex flex-col gap-y-4" v-if="state.items.length">
          <div :key="item.key" v-for="(item, idx) in state.items"
            class="flex flex-row items-start justify-center gap-x-6">
            <UFormGroup label="Type" :type="`items[${idx}].type`">
              <USelect v-model="item.type" placeholder="Select Type" :options="typeOptions" />
            </UFormGroup>

            <!-- 
            <UFormGroup label="Type" :type="`items[${idx}].type`">
              <UInput type="number" v-model="item.type" />
            </UFormGroup> -->

            <UFormGroup label="Data" :name="`items[${idx}].data`">
              <UInput type="text" placeholder="Data of item" v-model="item.data" />
            </UFormGroup>

            <UButton @click="remove(idx)" size="sm" :ui="{ rounded: 'rounded-full' }" color="red"
              icon="i-heroicons-trash" class="self-center">Remove Item</UButton>
          </div>
        </div>
        <div v-else>
          <span class="text-md font-bold">No item</span>
        </div>
      </div>

      <UButton :loading="loading" type="submit">Submit</UButton>
      <div v-if="error" class="flex items-center gap-x-2 rounded-lg border px-2 py-1">
        <span @click="error = null" class="cursor-pointer">X</span>
        <span class="text-sm text-red-500">{{ error }}</span>
      </div>
    </UForm>
  </div>
</template>
