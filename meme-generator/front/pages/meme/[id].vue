<script setup lang="ts">
import type { MemeDto } from '@/types/dto/meme'

const meme = ref<MemeDto | null>(null)

const route = useRoute()
const memeId = Number(route.params.id)

const { fetchMeme, getMemeImageUrl } = useMemeApi()

onMounted(async () => {
  meme.value = await fetchMeme(memeId)
})
</script>

<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" md="8" lg="6">
        <v-card>
          <v-img
            :src="getMemeImageUrl(memeId)"
            aspect-ratio="1"
            class="white--text align-end"
          >
            <v-card-title class="bg-black bg-opacity-50">
              {{ meme?.name }}
            </v-card-title>
          </v-img>
          <v-card-subtitle class="mt-2">
            Created by <strong>@{{ meme?.owner_username }}</strong>
          </v-card-subtitle>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
