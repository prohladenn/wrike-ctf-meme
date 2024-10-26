<script setup lang="ts">
import type { CreateTemplateDto } from '~/types/dto/template'
const newTemplate = ref({} as CreateTemplateDto)
const image = ref<File | null>(null)
const formRef = ref()

const { createTemplate } = useTemplateApi()
const router = useRouter()

const createTemplateHandler = async () => {
  if (!formRef.value.validate()) {
    return
  }

  if (!image.value) return

  const createdTemplate = await createTemplate(newTemplate.value, image.value)
  router.push(`/template/${createdTemplate.id}`)
}
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Create New Template</h1>
    <v-form ref="formRef" @submit.prevent="createTemplateHandler">
      <v-text-field
        v-model="newTemplate.name"
        label="Template Name"
        :rules="[(v) => !!v || 'Name is required']"
        prepend-icon="mdi-rename"
        required
      />
      <v-text-field
        v-model="newTemplate.comment"
        label="Comment"
        prepend-icon="mdi-text"
        required
      />
      <v-file-input
        v-model="image"
        label="Upload Image"
        accept="image/*"
        :rules="[(v) => !!v || 'Image is required']"
        required
      />
      <v-btn type="submit" color="primary">Create</v-btn>
    </v-form>
  </v-container>
</template>
