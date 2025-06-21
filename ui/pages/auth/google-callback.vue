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

type GoogleCallbackResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

// URL'den parametreleri al
function getUrlParams() {
  const urlParams = new URLSearchParams(window.location.search);
  return {
    code: urlParams.get('code'),
    error: urlParams.get('error')
  };
}

// Hata durumlarını kontrol et
function checkForErrors(code: string | null, error: string | null) {
  if (error) {
    errorMessage.value = t('login.error.googleAuthFailed');
    isLoading.value = false;
    return true;
  }

  if (!code) {
    errorMessage.value = t('login.error.noCode');
    isLoading.value = false;
    return true;
  }

  return false;
}

// Google API'sine istek gönder
async function sendGoogleRequest(code: string): Promise<GoogleCallbackResponse> {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), 30000); // 30 saniye timeout
  
  try {
    const response = await fetch(`${useRuntimeConfig().public.apiBase}/oauth/google/callback`, {
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

    return await response.json() as GoogleCallbackResponse;
  } finally {
    clearTimeout(timeoutId);
  }
}

// Başarılı giriş işlemlerini yap
function handleSuccessfulLogin(data: GoogleCallbackResponse) {
  // Kullanıcı bilgilerini kaydet
  localStorage.setItem('auth_token', data.token);
  localStorage.setItem('username', data.username);
  
  isLoading.value = false;
  successMessage.value = t('login.success.googleAuth');

  // Ana sayfaya yönlendir
  setTimeout(() => {
    navigateTo('/home');
  }, 2000);
}

// Hata durumlarını işle
function handleError(err: any) {
  console.error('Google callback hatası:', err);
  
  if (err instanceof Error) {
    if (err.name === 'AbortError') {
      return; // Timeout durumunda hata gösterme
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
  // 1. URL'den parametreleri al
  const { code, error } = getUrlParams();
  
  // 2. Hata kontrolü yap
  if (checkForErrors(code, error)) {
    return;
  }

  try {
    // 3. Google API'sine istek gönder
    const response = await sendGoogleRequest(code!);
    
    // 4. Yanıtı kontrol et
    if (response.type === 'error') {
      errorMessage.value = response.message || t('login.error.general');
      isLoading.value = false;
      return;
    }

    // 5. Başarılı giriş ise işlemleri yap
    if (response.type === 'google' && response.token) {
      handleSuccessfulLogin(response);
    } else {
      // Beklenmeyen yanıt formatı
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