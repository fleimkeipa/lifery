<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated"
});

const { t } = useI18n();

const isLoading = ref(true);
const errorMessage = ref('');
const successMessage = ref('');

type LinkedInCallbackResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

function getUrlParams() {
  const urlParams = new URLSearchParams(window.location.search);
  return {
    code: urlParams.get('code'),
    error: urlParams.get('error')
  };
}

function checkForErrors(code: string | null, error: string | null) {
  if (error) {
    errorMessage.value = t('login.error.linkedinAuthFailed');
    isLoading.value = false;
    return true;
  }

  if (!code) {
    errorMessage.value = t('login.error.noCodeLinkedIn');
    isLoading.value = false;
    return true;
  }

  return false;
}

async function sendLinkedInRequest(code: string): Promise<LinkedInCallbackResponse> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 30000);

  try {
    const response = await fetch(`${useRuntimeConfig().public.apiBase}/oauth/linkedin/callback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ code }),
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      throw new Error(`HTTP hatası! Durum: ${response.status}`);
    }

    return await response.json() as LinkedInCallbackResponse;
  } finally {
    clearTimeout(timeoutId);
  }
}

function handleSuccessfulLogin(data: LinkedInCallbackResponse) {
  localStorage.setItem('auth_token', data.token);
  localStorage.setItem('username', data.username);
  localStorage.setItem('auth_type', data.type);

  isLoading.value = false;
  successMessage.value = t('login.success.linkedinAuth');

  setTimeout(() => {
    navigateTo('/home');
  }, 2000);
}

function handleError(err: any) {
  console.error('LinkedIn callback hatası:', err);

  if (err instanceof Error) {
    if (err.name === 'AbortError') {
      return;
    }

    if (err.message.includes('Failed to fetch') || err.message.includes('NetworkError')) {
      errorMessage.value = t('login.error.networkError');
    } else {
      errorMessage.value = t('login.error.general');
    }
  } else {
    errorMessage.value = t('login.error.general');
  }

  isLoading.value = false;
}

onMounted(async () => {
  const { code, error } = getUrlParams();

  if (checkForErrors(code, error)) {
    return;
  }

  try {
    const response = await sendLinkedInRequest(code!);

    if (response.type === 'error') {
      errorMessage.value = response.message || t('login.error.general');
      isLoading.value = false;
      return;
    }

    if (response.type === 'linkedin' && response.token) {
      handleSuccessfulLogin(response);
    } else {
      console.error("Beklenmeyen yanıt formatı:", response);
      errorMessage.value = t('login.error.general');
      isLoading.value = false;
    }
  } catch (err) {
    handleError(err);
  }
});
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <div class="text-center">
        <div v-if="isLoading" class="mb-4">
          <UIcon name="i-lucide-loader-2" class="w-8 h-8 animate-spin mx-auto mb-4" />
          <p class="text-gray-600">{{ t('login.processing') }}</p>
        </div>

        <div v-else-if="errorMessage" class="mb-4">
          <UIcon name="i-lucide-x-circle" class="w-8 h-8 text-red-500 mx-auto mb-4" />
          <p class="text-red-600 mb-4">{{ errorMessage }}</p>
          <UButton @click="navigateTo('/login')" color="red" variant="outline">
            {{ t('login.backToLogin') }}
          </UButton>
        </div>

        <div v-else-if="successMessage" class="mb-4">
          <UIcon name="i-lucide-check-circle" class="w-8 h-8 text-green-500 mx-auto mb-4" />
          <p class="text-green-600 mb-4">{{ successMessage }}</p>
          <p class="text-gray-600 text-sm">{{ t('login.redirecting') }}</p>
        </div>
      </div>
    </UCard>
  </div>
</template> 