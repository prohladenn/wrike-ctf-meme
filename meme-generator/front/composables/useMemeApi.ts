import type { CreateMemeDto } from '@/types/dto/meme'

export function useMemeApi() {
  const { $api } = useNuxtApp()

  const fetchLastMemes = async () => {
    const response = await $api.get('/memes/last')
    return response.data
  }

  const fetchMeme = async (memeId: number) => {
    const response = await $api.get(`/meme/${memeId}`)
    return response.data
  }

  const fetchUserMemes = async (userId: number) => {
    const response = await $api.get(`/user/${userId}/memes`)
    return response.data
  }

  const createMeme = async (input: CreateMemeDto) => {
    const response = await $api.post('/meme', input)
    return response.data
  }

  const getMemePreview = async (input: CreateMemeDto) => {
    const response = await $api.post('/meme/preview', input, {
      responseType: 'blob',
    })
    return response.data
  }

  const getMemeImageUrl = (memeId: number) => {
    return `${useRuntimeConfig().public.apiBaseUrl}/meme/${memeId}/image`
  }

  return {
    fetchLastMemes,
    fetchMeme,
    fetchUserMemes,
    createMeme,
    getMemePreview,
    getMemeImageUrl,
  }
}
