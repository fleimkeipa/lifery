<script setup lang="ts">
import { number, object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "@nuxt/ui";
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});

const { t, locale } = useI18n();

const schema = object({
  username: string()
    .required(t('register.validation.required.username'))
    .min(3, t('register.validation.min.username')),
  email: string()
    .required(t('register.validation.required.email'))
    .email(t('register.validation.email')),
  password: string()
    .required(t('register.validation.required.password'))
    .min(6, t('register.validation.min.password')),
  confirm_password: string()
    .required(t('register.validation.required.confirmPassword'))
    .test('passwords-match', t('register.validation.passwordsMatch'), function (value) {
      return this.parent.password === value
    }),
  role_id: number()
    .required('Rol seçimi gereklidir')
});

type Schema = InferType<typeof schema>;

const state = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  role_id: 0
});

type RegisterResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  const { data, error } = await useFetch<RegisterResponse>("/auth/register", {
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
    await navigateTo('/home');
  }
}
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">{{ t('register.title') }}</h1>
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup :label="t('register.username')" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup :label="t('register.email')" name="email">
          <UInput v-model="state.email" type="email" />
        </UFormGroup>

        <UFormGroup :label="t('register.password')" name="password">
          <UInput v-model="state.password" type="password" />
        </UFormGroup>

        <UFormGroup :label="t('register.confirmPassword')" name="confirm_password">
          <UInput v-model="state.confirm_password" type="password" />
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block>{{ t('register.submit') }}</UButton>
        </div>

        <div class="flex w-full justify-center mt-4">
          <UButton :to="{ path: '/login', query: { locale } }" variant="link" class="text-sm">
            {{ t('register.alreadyHaveAccount') }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>
