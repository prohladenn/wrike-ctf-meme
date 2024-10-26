import type { CreateTemplateDto } from '@/types/dto/template'

export function useTemplateApi() {
  const { $api } = useNuxtApp()

  const fetchLastTemplates = async () => {
    const response = await $api.get('/templates/last')
    return response.data
  }

  const fetchTemplate = async (templateId: number) => {
    const response = await $api.get(`/template/${templateId}`)
    return response.data
  }

  const fetchUserTemplates = async (userId: number) => {
    const response = await $api.get(`/user/${userId}/templates`)
    return response.data
  }

  const createTemplate = async (input: CreateTemplateDto, file: File) => {
    const invalidCharsRegex = /[<>:"/\\|?*\x00-\x1F]/;
    if (invalidCharsRegex.test(file.name)) {
      throw new Error("Invalid file name: contains invalid characters")
    }

    const formData = new FormData()
    formData.append('image', file)
    formData.append('name', input.name)
    formData.append('comment', input.comment)
    const response = await $api.post('/template', formData)
    return response.data
  }

  const getTemplateImageUrl = (templateId: number) => {
    return `${useRuntimeConfig().public.apiBaseUrl}/template/${templateId}/image`
  }

  return {
    fetchLastTemplates,
    fetchTemplate,
    fetchUserTemplates,
    createTemplate,
    getTemplateImageUrl,
  }
}
