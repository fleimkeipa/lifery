<script setup lang="ts">
const { t } = useI18n();
import { onClickOutside } from '@vueuse/core';

const NOTIFICATION_UNREAD_STATUS = 100;
const NOTIFICATION_READ_STATUS = 101;

interface Notification {
  id: string;
  type: string;
  message: string;
  read: number;
  created_at: string;
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString();
};

const showDropdown = ref(false);
const unreadCount = ref(0);
const errorMessage = ref("");

const getQueryParams = () => {
  const params = new URLSearchParams({
    limit: '30',
    order: 'desc:created_at'
  });
  return params.toString();
};

const { data: notifications, error, isFetching, execute: fetchNotifications } =
  useApi(() =>
    `/notifications?${getQueryParams()}`
  ).json();

const markAsRead = async (id: string) => {
  try {
    await useApi(`/notifications/${id}`, {
      afterFetch: () => fetchNotifications(),
    }).patch({ read: NOTIFICATION_READ_STATUS });
    errorMessage.value = "";
  } catch (error) {
    errorMessage.value = t('notifications.mark_as_read') + ': ' + String((error as any)?.message || error);
    console.error('Failed to mark notification as read:', error);
  }
};

watch(() => notifications.value?.data, (newNotifications) => {
  if (newNotifications) {
    unreadCount.value = newNotifications.filter((notification: Notification) => notification.read === NOTIFICATION_UNREAD_STATUS).length;
  }
}, { immediate: true });

watch(showDropdown, async (val) => {
  if (val) {
    await fetchNotifications();
    errorMessage.value = "";
    if (notifications.value?.data) {
      const unread = notifications.value.data.filter((n: Notification) => n.read === NOTIFICATION_UNREAD_STATUS);
      await Promise.all(unread.map((n: Notification) => markAsRead(n.id)));
    }
  }
});

const notificationRef = ref<HTMLElement | null>(null);

onClickOutside(notificationRef, () => {
  showDropdown.value = false;
});
</script>

<template>
  <div class="relative" ref="notificationRef">
    <UButton @click="showDropdown = !showDropdown" color="gray" variant="ghost" icon="i-heroicons-bell" class="p-2">
      <span v-if="unreadCount > 0"
        class="absolute top-0 right-0 inline-flex items-center justify-center px-2 py-1 text-xs font-bold leading-none text-white transform translate-x-1/2 -translate-y-1/2 bg-red-500 rounded-full">
        {{ unreadCount }}
      </span>
    </UButton>

    <div v-if="showDropdown"
      class="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg overflow-hidden z-50 border border-gray-200">
      <div class="p-4 border-b border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900">{{ t('notifications.title') }}</h3>
      </div>

      <div class="max-h-96 overflow-y-auto">
        <div v-if="isFetching" class="p-4 text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto"></div>
        </div>

        <div v-else-if="errorMessage" class="p-4 text-center text-red-600">
          {{ errorMessage }}
        </div>

        <div v-else-if="!notifications?.data?.length" class="p-4 text-center text-gray-500">
          {{ t('notifications.no_notifications') }}
        </div>

        <div v-else class="divide-y divide-gray-200">
          <div v-for="notification in notifications.data" :key="notification.id" :class="[
            'p-4 transition-colors duration-150',
            notification.read === NOTIFICATION_UNREAD_STATUS ? 'bg-blue-100' : 'bg-white'
          ]">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <p class="text-sm"
                  :class="notification.read === NOTIFICATION_UNREAD_STATUS ? 'text-gray-900 font-bold' : 'text-gray-500'">
                  {{ notification.message }}
                </p>
                <p class="text-xs text-gray-500 mt-1">
                  {{ formatDate(notification.created_at) }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>