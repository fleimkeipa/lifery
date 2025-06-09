<template>
    <div class="space-y-4">
      <input type="file" @change="onFileChange" accept="video/*" />
      <button @click="uploadImage" :disabled="!selectedFile">Yükle</button>
  
      <div v-if="imageUrl">
        <p>Yüklenen Video:</p>
        <video :src="imageUrl" controls class="max-w-xs mt-2" />
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  const config = useRuntimeConfig()
  
  const selectedFile = ref<File | null>(null)
  const imageUrl = ref<string | null>(null)
  
  const emit = defineEmits(['upload-complete'])
  
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
  
    const res = await fetch(`https://api.cloudinary.com/v1_1/${config.public.cloudinaryCloudName}/video/upload`, {
      method: 'POST',
      body: formData,
    })
  
    const data = await res.json()
    imageUrl.value = data.secure_url
    emit('upload-complete', data.secure_url)
  }
  </script>
  
  <style scoped>
  button {
    padding: 0.5rem 1rem;
    background-color: #3b82f6;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
  }
  </style>
  