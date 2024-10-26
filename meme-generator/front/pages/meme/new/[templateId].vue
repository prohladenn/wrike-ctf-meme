<script setup lang="ts">
import type { TemplateDto } from '@/types/dto/template'
import type { CreateMemeDto } from '~/types/dto/meme'

const template = ref<TemplateDto | null>(null)
const formRef = ref()

const route = useRoute()
const router = useRouter()
const templateId = Number(route.params.templateId)

const previewImageUrl = ref<string | null>(null)

const newMeme = ref({
  name: '',
  caption: '',
  color: 'white',
} as CreateMemeDto)

const { fetchTemplate, getTemplateImageUrl } = useTemplateApi()
const { createMeme, getMemePreview } = useMemeApi()

onMounted(async () => {
  template.value = await fetchTemplate(templateId)
})

const createMemeHandler = async () => {
  if (!formRef.value.validate()) {
    return
  }

  if (!template.value) return

  if (template.value.id !== 0) {
    newMeme.value.template_id = template.value.id.toString()
  }

  const createdMeme = await createMeme(newMeme.value)
  router.push(`/meme/${createdMeme.id}`) // Redirect to home page or memes list
}

const generatePreview = async () => {
  if (!formRef.value.validate()) {
    return
  }

  if (!template.value) return

  newMeme.value.template_id = template.value.id.toString()

  try {
    const blob = await getMemePreview(newMeme.value)

    if (!(blob instanceof Blob)) {
      console.error('Received data is not a Blob:', blob)
      return
    }

    if (previewImageUrl.value) {
      URL.revokeObjectURL(previewImageUrl.value)
    }
    previewImageUrl.value = URL.createObjectURL(blob)
  } catch (error) {
    console.error('Error generating preview:', error)
    Notify.error('Failed to generate meme preview.')
  }
}
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Create New Meme</h1>
    <v-form ref="formRef" @submit.prevent="createMemeHandler">
      <v-row>
        <v-col cols="12" md="12">
          <v-text-field
            v-model="newMeme.name"
            label="Meme Name"
            :rules="[(v) => !!v || 'Name is required']"
            prepend-icon="mdi-rename"
            required
          />
          <v-textarea
            v-model="newMeme.caption"
            label="Caption"
            :rules="[(v) => !!v || 'Caption is required']"
            prepend-icon="mdi-text"
            required
          />
          <v-select
            v-model="newMeme.color"
            :items="['black', 'white']"
            label="Text Color"
            prepend-icon="mdi-palette"
            :rules="[(v) => !!v || 'Color is required']"
            required
          />
          <v-btn color="primary" @click="generatePreview" class="mt-4 mr-2">
            Generate Preview
          </v-btn>
          <v-btn color="success" @click="createMemeHandler" class="mt-4">
            Create Meme
          </v-btn>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" md="6">
          <h3 class="text-h6">Template Image</h3>
          <v-img
            v-if="template"
            :src="getTemplateImageUrl(template.id)"
            class="meme-image my-2"
          />
        </v-col>
        <v-col cols="12" md="6">
          <h3 class="text-h6">Preview Image</h3>
          <v-img
            v-if="previewImageUrl"
            :src="previewImageUrl"
            class="meme-image my-2"
          />
        </v-col>
      </v-row>
    </v-form>
  </v-container>
</template>

<style scoped>
.meme-image {
  max-width: 100%;
  width: 300px; /* Adjust this value as needed */
  height: auto;
  display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>
