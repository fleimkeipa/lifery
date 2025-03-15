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

type LoginResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  const { data, error } = await useFetch<LoginResponse>("/auth/login", {
    method: "POST",
    baseURL: useRuntimeConfig().public.apiBase,
    body: event.data
  });

  if (error.value) {
    console.error(error.value);
    return;
  }

  if (data.value) {
    // Store the token in localStorage
    localStorage.setItem('auth_token', data.value.token);
    localStorage.setItem('username', data.value.username);
    
    // Navigate to the events page
    await navigateTo('/events');
  }
}
</script>

<template>
  <div
    class="flex h-screen w-full flex-col items-center justify-center gap-y-16"
  >
    <h1 class="text-6xl font-bold">Kubernetes UI</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">Login</h1>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
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
