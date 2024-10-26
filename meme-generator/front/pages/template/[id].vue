<script setup lang="ts">
import type { TemplateDto } from '@/types/dto/template'

const template = ref<TemplateDto | null>(null)
const username = ref<string>('Unknown')
const isLoading = ref(true)

const route = useRoute()
const router = useRouter()
const templateId = Number(route.params.id)

const { fetchTemplate, getTemplateImageUrl } = useTemplateApi()
const { fetchUser } = useUserApi()

const createMeme = () => {
  router.push(`/meme/new/${templateId}`)
}

onMounted(async () => {
  template.value = await fetchTemplate(templateId)
  if (template.value) {
    const user = await fetchUser(template.value.owner_id)
    if (user) {
      username.value = user.username
    }
  }
})

watch(template, () => {
  if (template.value) {
    isLoading.value = false
  }
})
</script>

<template>
  <v-container v-if="!isLoading">
    <v-row justify="center">
      <v-col cols="12" md="8" lg="6">
        <v-card v-if="template">
          <v-img
            :src="getTemplateImageUrl(template.id)"
            aspect-ratio="1"
            class="white--text align-end"
          >
            <v-card-title class="bg-black bg-opacity-50">
              {{ template.name }}
            </v-card-title>
            <v-card-text
              v-if="template.private_info"
              class="bg-grey bg-opacity-50 mt-2 pt-2"
            >
              {{ template.private_info }}
            </v-card-text>
          </v-img>
          <v-card-subtitle class="mt-2">
            Created by <strong>{{ username }}</strong>
          </v-card-subtitle>
          <v-card-actions>
            <v-btn color="primary" @click="createMeme"> Create Meme </v-btn>
          </v-card-actions>
        </v-card>
        <v-alert v-else type="error" dismissible> Template not found. </v-alert>
      </v-col>
    </v-row>
  </v-container>
  <v-container v-else>
    <v-row justify="center">
      <v-col cols="12" md="8" lg="6">
        <v-progress-circular indeterminate color="primary" />
      </v-col>
    </v-row>
  </v-container>
</template>
