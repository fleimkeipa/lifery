<script setup lang="ts">
import { object, string, type InferType } from "yup";
import { ref, reactive } from 'vue'
import type { FormSubmitEvent } from "@nuxt/ui";
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});

const { t, locale } = useI18n();
const route = useRoute();

const schema = object({
  newPassword: string()
    .required(t('resetPassword.validation.required.newPassword'))
    .min(8, t('resetPassword.validation.min.newPassword')),
  confirmPassword: string()
    .required(t('resetPassword.validation.required.confirmPassword'))
    .test('passwords-match', t('resetPassword.validation.passwordsMatch'), function (value) {
      return this.parent.newPassword === value
    }),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  newPassword: undefined,
  confirmPassword: undefined,
});

const showPassword = ref(false);
const showConfirmPassword = ref(false);
const isLoading = ref(false);
const successMessage = ref('');
const errorMessage = ref('');

type ResetPasswordResponse = {
  message: string;
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  isLoading.value = true;
  successMessage.value = '';
  errorMessage.value = '';
  
  const token = route.query.token as string;
  
  if (!token) {
    errorMessage.value = 'Invalid reset token';
    isLoading.value = false;
    return;
  }

  const { data, error } = await useFetch<ResetPasswordResponse>("/auth/reset-password", {
    method: "POST",
    baseURL: useRuntimeConfig().public.apiBase,
    body: {
      token,
      new_password: event.data.newPassword,
      confirm_password: event.data.confirmPassword,
    }
  });

  isLoading.value = false;

  if (error.value) {
    console.error(error.value);
    errorMessage.value = 'An error occurred while resetting your password';
    return;
  }

  if (data.value) {
    successMessage.value = 'Password reset successfully! You can now login with your new password.';
    // Redirect to login after 3 seconds
    setTimeout(() => {
      navigateTo('/login');
    }, 3000);
  }
}
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">{{ t('resetPassword.title') }}</h1>
      
      <div v-if="successMessage" class="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
        {{ successMessage }}
      </div>
      
      <div v-if="errorMessage" class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
        {{ errorMessage }}
      </div>
      
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup :label="t('resetPassword.newPassword')" name="newPassword">
          <UInput v-model="state.newPassword" :type="showPassword ? 'text' : 'password'">
            <template #trailing>
              <UButton class="pointer-events-auto" color="gray" variant="link"
                :icon="showPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="showPassword ? 'Hide password' : 'Show password'" :aria-pressed="showPassword" aria-controls="newPassword"
                @click.prevent="showPassword = !showPassword" />
            </template>
          </UInput>
        </UFormGroup>

        <UFormGroup :label="t('resetPassword.confirmPassword')" name="confirmPassword">
          <UInput v-model="state.confirmPassword" :type="showConfirmPassword ? 'text' : 'password'">
            <template #trailing>
              <UButton class="pointer-events-auto" color="gray" variant="link"
                :icon="showConfirmPassword ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="showConfirmPassword ? 'Hide password' : 'Show password'" :aria-pressed="showConfirmPassword" aria-controls="confirmPassword"
                @click.prevent="showConfirmPassword = !showConfirmPassword" />
            </template>
          </UInput>
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block :loading="isLoading">
            {{ t('resetPassword.submit') }}
          </UButton>
        </div>

        <div class="flex w-full justify-center mt-4">
          <UButton :to="{ path: '/login', query: { locale } }" variant="link" class="text-sm">
            {{ t('forgotPassword.backToLogin') }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template> 