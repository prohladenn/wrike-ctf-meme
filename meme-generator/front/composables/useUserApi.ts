export function useUserApi() {
  const { $api } = useNuxtApp()

  const fetchUser = async (id: number) => {
    const response = await $api.get(`/user/${id}`)
    return response.data
  }

  const fetchUsers = async (limit: number = 10, page: number = 1) => {
    const response = await $api.get(`/users?limit=${limit}&page=${page}`)
    return response.data
  }

  return {
    fetchUser,
    fetchUsers,
  }
}
