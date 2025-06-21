<script setup lang="ts">
import { ref } from 'vue'

enum EventType {
  "STRING" = 10,
  "PHOTO" = 11,
  "VIDEO" = 12,
  "VOICE_RECORD" = 13
}

const isPopupOpen = ref(false)
const selectedText = ref('')
const selectedImage = ref('')
const selectedVideo = ref('')

const openPopup = (data: string, type: number) => {
  if (type === EventType.STRING) {
    selectedText.value = data
    selectedImage.value = ''
    selectedVideo.value = ''
  } else if (type === EventType.PHOTO) {
    selectedImage.value = data
    selectedText.value = ''
    selectedVideo.value = ''
  } else if (type === EventType.VIDEO) {
    selectedVideo.value = data
    selectedText.value = ''
    selectedImage.value = ''
  }
  isPopupOpen.value = true
}

defineProps<{
    items: {
        data: string;
        type: number;
    }[];
    color: string;
}>()
</script>

<template>
    <div class="grid grid-cols-3 gap-2 items-center w-full">
        <div v-for="item in items" :key="item.data" class="flex justify-center items-center">
            <div v-if="item.type===EventType.PHOTO">
                <img :src="item.data" 
                     class="max-h-12 w-auto object-contain rounded cursor-pointer hover:opacity-90" 
                     @click="openPopup(item.data, item.type)" />
            </div>
            <div v-if="item.type===EventType.STRING" 
                 :class="color" 
                 class="w-full h-12 rounded shadow-lg flex items-center justify-center text-center text-gray-800 cursor-pointer hover:opacity-90 p-2 break-all"
                 @click="openPopup(item.data, item.type)">
                <p class="truncate">{{ item.data }}</p>
            </div>
            <div v-if="item.type===EventType.VIDEO" 
                 :class="color" 
                 class="w-full h-12 rounded shadow-lg text-center text-gray-800 cursor-pointer hover:opacity-90 relative flex items-center justify-center"
                 @click="openPopup(item.data, item.type)">
                <div class="flex flex-col items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" viewBox="0 0 24 24" fill="currentColor">
                        <path d="M8 5v14l11-7z"/>
                    </svg>
                    <span class="text-xs">Video</span>
                </div>
            </div>
        </div>
    </div>

    <!-- Popup -->
    <Teleport to="body">
        <div v-if="isPopupOpen" class="fixed inset-0 z-50">
            <!-- Arka plan overlay -->
            <div class="absolute inset-0 bg-black z-40 opacity-70" @click="isPopupOpen = false"></div>
            
            <!-- Popup içeriği -->
            <div class="absolute inset-0 flex items-center justify-center z-50">
                <div class="bg-white p-6 rounded-lg shadow-xl max-w-2xl w-full mx-4 flex flex-col">
                    <div class="flex-1">
                        <div v-if="selectedText" class="text-gray-800 whitespace-pre-wrap">{{ selectedText }}</div>
                        <div v-if="selectedImage" class="flex justify-center">
                            <img :src="selectedImage" 
                                 class="max-h-[80vh] max-w-full object-contain" />
                        </div>
                        <div v-if="selectedVideo" class="flex justify-center">
                            <video :src="selectedVideo"
                                controls 
                                 class="max-h-[80vh] max-w-full object-contain" />
                        </div>
                    </div>
                    <div class="mt-4 flex justify-end">
                        <button @click="isPopupOpen = false" 
                                class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300">
                            Kapat
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </Teleport>
</template>