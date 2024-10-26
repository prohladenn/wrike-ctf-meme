<script setup lang="ts">
import type { LoginDto } from '~/types/dto/auth'

const { login } = useAuthStore()

const user = ref<LoginDto>({
  username: '',
  password: '',
})

const handleLogin = async () => {
  await login(user.value)
  await navigateTo('/')
}
</script>

<template>
  <div class="wrapper">
    <v-responsive max-width="600" class="mx-auto">
      <v-card class="pa-2 pt-4">
        <v-card-title>
          <h2>Login</h2>
        </v-card-title>
        <v-card-text>
          <v-form>
            <v-text-field
              v-model="user.username"
              label="Username"
              color="required"
            />
            <v-text-field
              v-model="user.password"
              label="Password"
              type="password"
              required
            />
          </v-form>
          <p class="mb-2">
            Not registered yet? <NuxtLink to="/register">Register</NuxtLink>
          </p>
          <v-btn type="submit" color="primary" @click="handleLogin"
            >Login</v-btn
          >
        </v-card-text>
      </v-card>
    </v-responsive>
  </div>
</template>

<style scoped>
.wrapper {
  position: relative;
  top: calc(50vh - 330px);
  text-align: center;
}
</style>
