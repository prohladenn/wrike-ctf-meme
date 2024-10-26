<script setup lang="ts">
import type { UserDto } from '@/types/dto/user'
import type { MemeDto } from '@/types/dto/meme'
import type { TemplateDto } from '@/types/dto/template'

const user = ref<UserDto | null>(null)
const memes = ref<MemeDto[]>([])
const templates = ref<TemplateDto[]>([])

const route = useRoute()
const router = useRouter()
const userId = Number(route.params.id)

const { fetchUser } = useUserApi()
const { fetchUserMemes, getMemeImageUrl } = useMemeApi()
const { fetchUserTemplates, getTemplateImageUrl } = useTemplateApi()

const goToMeme = (id: number) => {
  router.push(`/meme/${id}`)
}

const goToTemplate = (id: number) => {
  router.push(`/template/${id}`)
}

onMounted(async () => {
  user.value = await fetchUser(userId)

  if (user.value) {
    memes.value = await fetchUserMemes(userId)
    templates.value = await fetchUserTemplates(userId)
  }
})
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">{{ user?.username }}'s Profile</h1>
    <v-expansion-panels>
      <v-expansion-panel>
        <v-expansion-panel-title>Meme List</v-expansion-panel-title>
        <v-expansion-panel-text>
          <v-row>
            <v-col
              v-for="meme in memes"
              :key="meme.id"
              cols="12"
              sm="6"
              md="4"
              lg="3"
            >
              <v-card class="cursor-pointer" @click="goToMeme(meme.id)">
                <v-img
                  :src="getMemeImageUrl(meme.id)"
                  height="200"
                  class="white--text align-end"
                >
                  <v-card-title class="bg-black bg-opacity-50">
                    {{ meme.name }}
                  </v-card-title>
                </v-img>
              </v-card>
            </v-col>
          </v-row>
        </v-expansion-panel-text>
      </v-expansion-panel>
      <v-expansion-panel>
        <v-expansion-panel-title>Template List</v-expansion-panel-title>
        <v-expansion-panel-text>
          <v-row>
            <v-col
              v-for="template in templates"
              :key="template.id"
              cols="12"
              sm="6"
              md="4"
              lg="3"
            >
              <v-card class="cursor-pointer" @click="goToTemplate(template.id)">
                <v-img
                  :src="getTemplateImageUrl(template.id)"
                  height="200"
                  class="white--text align-end"
                >
                  <v-card-title class="bg-black bg-opacity-50">
                    {{ template.name }}
                  </v-card-title>
                </v-img>
              </v-card>
            </v-col>
          </v-row>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-container>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
