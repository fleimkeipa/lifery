<template>
    <div class="space-y-4">
      <div class="flex items-center space-x-2">
        <label for="video-upload" class="custom-file-upload">
          {{ $t('event.select_video_button') }}
        </label>
        <input id="video-upload" type="file" @change="onFileChange" accept="video/*" class="hidden" />
        <span class="text-sm text-gray-500">{{ selectedFile ? selectedFile.name : $t('event.no_video_selected') }}</span>
      </div>
      <button @click="uploadVideo" :disabled="!selectedFile">{{ $t('event.upload_video_button') }}</button>
  
      <div v-if="videoUrl">
        <p>{{ $t('event.uploaded_video') }}</p>
        <video :src="videoUrl" controls class="max-w-xs mt-2" />
      </div>
      
      <!-- Display existing video if provided -->
      <div v-if="modelValue && !videoUrl">
        <p>{{ $t('event.existing_video') }}</p>
        <video :src="modelValue" controls class="max-w-xs mt-2" />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  const config = useRuntimeConfig()
  
  const props = defineProps<{
    modelValue?: string
  }>()
  
  const selectedFile = ref<File | null>(null)
  const videoUrl = ref<string | null>(null)
  
  const emit = defineEmits(['upload-complete'])
  
  // Initialize with existing value if provided
  if (props.modelValue) {
    videoUrl.value = props.modelValue
  }
  
  function onFileChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (files && files[0]) {
      selectedFile.value = files[0]
    }
  }
  
  async function uploadVideo() {
    if (!selectedFile.value) return
  
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('upload_preset', config.public.cloudinaryUploadPreset)
  
    const res = await fetch(`https://api.cloudinary.com/v1_1/${config.public.cloudinaryCloudName}/video/upload`, {
      method: 'POST',
      body: formData,
    })
  
    const data = await res.json()
    videoUrl.value = data.secure_url
    emit('upload-complete', data.secure_url)
  }
  </script>
  
  <style scoped>
  .custom-file-upload {
    border: 1px solid #ccc;
    display: inline-block;
    padding: 6px 12px;
    cursor: pointer;
    background-color: #f0f0f0;
    border-radius: 6px;
  }

  button {
    padding: 0.5rem 1rem;
    background-color: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  </style>
  