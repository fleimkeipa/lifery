<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";

definePageMeta({
  layout: "not-authenticated",
});

const schema = object({
  username: string().required("Required"),
  password: string()
    .min(8, "Must be at least 8 characters")
    .required("Required"),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  username: undefined,
  password: undefined,
});

// async function onSubmit(event: FormSubmitEvent<Schema>) {
//   useApi(`/auth/login`, {

//   }).post(event.data);
// }

const router = useRouter();
const loading = ref(false);
const error = ref(null);



const onSubmit = (event: FormSubmitEvent<Schema>) => {
  useApi(`/auth/login`, {
    afterFetch: () => {
      loading.value = false;
      router.push("/pods");
    },
    onFetchError: ({ error: fetchErr }) => {
      loading.value = false;
      error.value = fetchErr;
    },
  }).post(event.data);
};
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Kubernetes UI</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">Login</h1>
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup label="Username" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup label="Password" name="password">
          <UInput v-model="state.password" type="password" />
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block> Submit </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>
