<script setup lang="ts">
import { ref, computed } from 'vue';
import { useAuth } from '../../composables/useAuth';

const { user } = useAuth();

definePageMeta({
  middleware: "custom-layout",
});

const { t, locale } = useI18n();

// Check if user is authenticated via Google
const isOAuthUser = computed(() => {
  if (process.client) {
    return localStorage.getItem('auth_type') === 'google' || localStorage.getItem('auth_type') === 'linkedin';
  }
  return false;
});

const form = ref({
  username: user.value?.username || '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const usernameLoading = ref(false);
const passwordLoading = ref(false);
const message = ref<{ show: boolean; text: string; type: 'success' | 'error' }>({
  show: false,
  text: '',
  type: 'success'
});

const showMessage = (text: string, type: 'success' | 'error' = 'success') => {
  message.value = { show: true, text, type };
  setTimeout(() => {
    message.value.show = false;
  }, 3000);
};

const updateUsername = async () => {
  if (!form.value.username.trim()) {
    showMessage(t('settings.username_required'), 'error');
    return;
  }

  usernameLoading.value = true;
  try {
    const { data, error } = await useFetch('/user/username', {
      method: 'PUT',
      baseURL: useRuntimeConfig().public.apiBase,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
        'Content-Type': 'application/json',
      },
      body: {
        username: form.value.username
      }
    });

    if (error.value) {
      console.error('Username update error:', error.value);
      const errorMessage = error.value.data?.message || error.value.message || t('settings.username_update_error');
      showMessage(errorMessage, 'error');
      return;
    }

    if (data.value) {
      showMessage(t('settings.username_updated'), 'success');
      // Update the user object in auth
      if (user.value) {
        user.value.username = form.value.username;
      }
    }
  } catch (error: any) {
    console.error('Username update error:', error);
    showMessage(error?.data?.message || t('settings.username_update_error'), 'error');
  } finally {
    usernameLoading.value = false;
  }
};

const updatePassword = async () => {
  if (!form.value.currentPassword || !form.value.newPassword || !form.value.confirmPassword) {
    showMessage(t('settings.all_fields_required'), 'error');
    return;
  }

  if (form.value.newPassword !== form.value.confirmPassword) {
    showMessage(t('settings.passwords_not_match'), 'error');
    return;
  }

  if (form.value.newPassword.length < 8) {
    showMessage(t('settings.password_min_length'), 'error');
    return;
  }

  passwordLoading.value = true;
  try {
    const { data, error } = await useFetch('/user/password', {
      method: 'PUT',
      baseURL: useRuntimeConfig().public.apiBase,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
        'Content-Type': 'application/json',
      },
      body: {
        current_password: form.value.currentPassword,
        new_password: form.value.newPassword
      }
    });

    if (error.value) {
      console.error('Password update error:', error.value);
      const errorMessage = error.value.data?.message || error.value.message || t('settings.password_update_error');
      showMessage(errorMessage, 'error');
      return;
    }

    if (data.value) {
      showMessage(t('settings.password_updated'), 'success');
      // Clear password fields
      form.value.currentPassword = '';
      form.value.newPassword = '';
      form.value.confirmPassword = '';
    }
  } catch (error: any) {
    console.error('Password update error:', error);
    showMessage(error?.data?.message || t('settings.password_update_error'), 'error');
  } finally {
    passwordLoading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen bg-background p-8">
    <div class="max-w-2xl mx-auto">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <h1 class="text-2xl font-bold mb-6 text-gray-900 dark:text-white">
          {{ t('settings.title') }}
        </h1>

        <!-- Message -->
        <div v-if="message.show" 
             :class="[
               'mb-4 p-4 rounded-md',
               message.type === 'success' ? 'bg-green-100 text-green-700 dark:bg-green-800 dark:text-green-200' : 'bg-red-100 text-red-700 dark:bg-red-800 dark:text-red-200'
             ]">
          {{ message.text }}
        </div>

        <!-- Username Section -->
        <div class="mb-8">
          <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            {{ t('settings.change_username') }}
          </h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {{ t('common.username') }}
              </label>
              <input
                v-model="form.username"
                type="text"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                :placeholder="t('common.username')"
              />
            </div>
            <button
              @click="updateUsername"
              :disabled="usernameLoading || passwordLoading"
              class="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ usernameLoading ? t('common.loading') : t('settings.update_username') }}
            </button>
          </div>
        </div>

        <!-- Password Section -->
        <div v-if="!isOAuthUser">
          <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            {{ t('settings.change_password') }}
          </h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {{ t('settings.current_password') }}
              </label>
              <input
                v-model="form.currentPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                :placeholder="t('settings.current_password')"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {{ t('settings.new_password') }}
              </label>
              <input
                v-model="form.newPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                :placeholder="t('settings.new_password')"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {{ t('settings.confirm_new_password') }}
              </label>
              <input
                v-model="form.confirmPassword"
                type="password"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                :placeholder="t('settings.confirm_new_password')"
              />
            </div>
            <button
              @click="updatePassword"
              :disabled="usernameLoading || passwordLoading"
              class="w-full bg-green-500 text-white py-2 px-4 rounded-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ passwordLoading ? t('common.loading') : t('settings.update_password') }}
            </button>
          </div>
        </div>

        <!-- Google User Info -->
        <div v-if="isOAuthUser" class="mb-8">
          <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
            {{ t('settings.oauth_account_info') }}
          </h2>
          <div class="p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md">
            <p class="text-blue-800 dark:text-blue-200 text-sm">
              {{ t('settings.oauth_account_message') }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template> 