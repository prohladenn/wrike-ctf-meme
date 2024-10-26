<script setup lang="ts">
import type { TemplateDto } from '@/types/dto/template'

definePageMeta({
  icon: 'mdi-panorama-variant-outline',
  title: 'Meme templates',
  drawerIndex: 1,
})

const router = useRouter()

const { fetchLastTemplates, getTemplateImageUrl } = useTemplateApi()

const templates = ref<TemplateDto[]>([])

const goToTemplate = (id: number) => {
  router.push(`/template/${id}`)
}

onMounted(async () => {
  templates.value = await fetchLastTemplates()
})
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="6">
        <h1 class="text-h4 mb-4">Last Templates</h1>
      </v-col>
      <v-col cols="6" class="text-right">
        <v-btn color="primary" to="/template/new">Create new</v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-col
        v-for="template in templates"
        :key="template.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card>
          <v-img
            :src="getTemplateImageUrl(template.id)"
            height="200"
            class="white--text align-end cursor-pointer"
            @click="goToTemplate(template.id)"
          >
            <v-card-title class="bg-black bg-opacity-50">
              {{ template.name }}
            </v-card-title>
          </v-img>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
