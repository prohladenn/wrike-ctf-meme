<script setup lang="ts">
import { mergeProps } from 'vue'

const theme = useTheme()
const drawer = useState('drawer')
const route = useRoute()
const breadcrumbs = computed(() => {
  return route!.matched
    .filter(
      (item) =>
        item.meta && item.meta.title && !(item.meta?.breadcrumb === 'hidden'),
    )
    .map((r) => ({
      title: r.meta.title!,
      disabled:
        r.meta?.breadcrumb === 'disabled' || r.path === route.path || false,
      to: r.path,
    }))
})
const isDark = computed({
  get() {
    return theme.global.name.value === 'dark' ? true : false
  },
  set(v) {
    theme.global.name.value = v ? 'dark' : 'light'
  },
})

const authStore = useAuthStore()

const handleLogout = async () => {
  await authStore.logout()
  navigateTo('/login')
}

onMounted(async () => {
  await authStore.fetchUser()
})
</script>

<template>
  <v-app-bar flat>
    <v-app-bar-nav-icon @click="drawer = !drawer" />
    <v-breadcrumbs :items="breadcrumbs" />
    <v-spacer />
    <div id="app-bar" />
    <v-switch
      v-model="isDark"
      color=""
      hide-details
      density="compact"
      inset
      false-icon="mdi-white-balance-sunny"
      true-icon="mdi-weather-night"
      class="opacity-80"
    />
    <v-menu location="bottom">
      <template #activator="{ props: menu }">
        <v-tooltip location="bottom">
          <template #activator="{ props: tooltip }">
            <v-btn icon v-bind="mergeProps(menu, tooltip)" class="ml-1">
              <v-icon
                v-if="!authStore.loggedIn"
                icon="mdi-account-circle"
                size="36"
              />
              <v-avatar v-else color="primary" size="36">
                <v-img :src="authStore.user?.avatar_url" />
              </v-avatar>
            </v-btn>
          </template>
          <span>{{
            authStore.loggedIn ? authStore.user!.username : 'User'
          }}</span>
        </v-tooltip>
      </template>
    </v-menu>
    <v-btn
      v-if="!authStore.loggedIn"
      icon
      title="Login"
      to="/login"
      class="ml-1"
    >
      <v-icon icon="mdi-login" />
    </v-btn>
    <v-btn v-else icon title="Logout" class="ml-1" @click="handleLogout">
      <v-icon icon="mdi-logout" />
    </v-btn>
  </v-app-bar>
</template>
