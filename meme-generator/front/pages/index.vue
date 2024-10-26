<script setup lang="ts">
import type { MemeDto } from '@/types/dto/meme'

definePageMeta({
  icon: 'mdi-emoticon-excited-outline',
  title: 'Memes',
  drawerIndex: 0,
})

const router = useRouter()

const { fetchLastMemes, getMemeImageUrl } = useMemeApi()

const memes = ref<MemeDto[]>([])

const goToMeme = (id: number) => {
  router.push(`/meme/${id}`)
}

onMounted(async () => {
  memes.value = await fetchLastMemes()
})
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Last Memes</h1>
    <v-row>
      <v-col
        v-for="meme in memes"
        :key="meme.id"
        cols="12"
        sm="6"
        md="4"
        lg="3"
      >
        <v-card>
          <v-img
            :src="getMemeImageUrl(meme.id)"
            height="200"
            class="white--text align-end cursor-pointer"
            @click="goToMeme(meme.id)"
          >
            <v-card-title class="bg-black bg-opacity-50">
              {{ meme.name }}
            </v-card-title>
          </v-img>
          <v-card-subtitle class="mt-2">
            By
            <NuxtLink :to="`/user/${meme.owner_id}`" class="font-weight-bold"
              >@{{ meme.owner_username }}</NuxtLink
            >
          </v-card-subtitle>
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
