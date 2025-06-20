<template>
    <div class="space-y-4">
      <div class="flex items-center space-x-2">
        <label for="file-upload" class="custom-file-upload">
          {{ $t('event.select_file_button') }}
        </label>
        <input id="file-upload" type="file" @change="onFileChange" accept="image/*" class="hidden" />
        <span class="text-sm text-gray-500">{{ selectedFile ? selectedFile.name : $t('event.no_file_selected') }}</span>
      </div>
      <button @click="uploadImage" :disabled="!selectedFile">{{ $t('event.upload_image') }}</button>
  
      <div v-if="imageUrl">
        <p>{{ $t('event.uploaded_image') }}</p>
        <img :src="imageUrl" :alt="$t('event.uploaded_image_alt')" class="max-w-xs mt-2" />
      </div>
      
      <!-- Display existing image if provided -->
      <div v-if="modelValue && !imageUrl">
        <p>{{ $t('event.existing_image') }}</p>
        <img :src="modelValue" :alt="$t('event.existing_image_alt')" class="max-w-xs mt-2" />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  const config = useRuntimeConfig()
  
  const props = defineProps<{
    modelValue?: string
  }>()
  
  const selectedFile = ref<File | null>(null)
  const imageUrl = ref<string | null>(null)
  
  const emit = defineEmits(['upload-complete'])
  
  // Initialize with existing value if provided
  if (props.modelValue) {
    imageUrl.value = props.modelValue
  }
  
  function onFileChange(event: Event) {
    const files = (event.target as HTMLInputElement).files
    if (files && files[0]) {
      selectedFile.value = files[0]
    }
  }
  
  async function uploadImage() {
    if (!selectedFile.value) return
  
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('upload_preset', config.public.cloudinaryUploadPreset)
  
    const res = await fetch(`https://api.cloudinary.com/v1_1/${config.public.cloudinaryCloudName}/image/upload`, {
      method: 'POST',
      body: formData,
    })
  
    const data = await res.json()
    imageUrl.value = data.secure_url
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
  