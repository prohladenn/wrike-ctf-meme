import { defineNuxtPlugin } from '#app'
import axios from 'axios'

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl

  const api = axios.create({
    baseURL: apiBaseUrl,
    withCredentials: true, // For session cookie auth
  })

  // Add a request interceptor for CSRF protection
  api.interceptors.request.use(
    (config) => {
      const csrfToken = document.cookie.split('; ').find(row => row.startsWith('csrfToken='));
      if (csrfToken) {
        config.headers['X-CSRF-Token'] = csrfToken.split('=')[1];
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Add a response interceptor
  api.interceptors.response.use(
    (response) => {
      // If the response is successful, just return it
      return response
    },
    (error) => {
      // Handle error responses
      const { response } = error

      if (response) {
        if (response.status === 401) {
          // Check if the request URL is not '/user/me' before redirecting
          if (response.config.url !== '/user/me') {
            nuxtApp.$router.push('/login')
          }
        } else {
          // Other errors
          const errorMessage = response.data.error || 'An error occurred'

          // Use the global Notify.error(msg) function
          if (
            typeof Notify !== 'undefined' &&
            typeof Notify.error === 'function'
          ) {
            Notify.error(errorMessage)
          } else {
            console.error('Notify.error is not available')
          }
        }
      } else {
        // Network error or other issues
        if (
          typeof Notify !== 'undefined' &&
          typeof Notify.error === 'function'
        ) {
          Notify.error('Network error: Please check your internet connection.')
        } else {
          console.error('Notify.error is not available')
        }
      }

      // Return a rejected promise to prevent further processing
      return Promise.reject(error)
    },
  )

  return {
    provide: {
      api,
    },
  }
})
