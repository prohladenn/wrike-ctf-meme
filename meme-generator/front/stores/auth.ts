import { defineStore } from 'pinia'

import type { LoginDto, RegisterDto } from '@/types/dto/auth'
import type { User } from '@/types/user'

export const useAuthStore = defineStore('auth', {
  state: () => {
    return {
      user: null as User | null,
    }
  },
  actions: {
    async login(credentials: LoginDto) {
      const { $api } = useNuxtApp()
      await $api.post('/login', credentials)
      await this.fetchUser()
    },
    async register(credentials: RegisterDto) {
      const { $api } = useNuxtApp()
      await $api.post('/register', credentials)
    },
    async logout() {
      const { $api } = useNuxtApp()
      await $api.post('/logout')
      this.user = null
    },
    async fetchUser() {
      const { $api } = useNuxtApp()
      const response = await $api.get('/user/me')
      this.user = response.data
      this.user.avatar_url = `https://ui-avatars.com/api/?name=${this.user.username}&background=0D8ABC&color=fff`
    },
  },
  getters: {
    loggedIn(state) {
      return !!state.user
    },
  },
})
