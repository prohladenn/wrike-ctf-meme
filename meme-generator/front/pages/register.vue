<script setup lang="ts">
import type { LoginDto } from '~/types/dto/auth'

const { register } = useAuthStore()

const user = ref<LoginDto>({
  username: '',
  password: '',
})

const handleRegister = async () => {
  await register(user.value)
  await navigateTo('/login')
}
</script>

<template>
  <div class="wrapper">
    <v-responsive max-width="600" class="mx-auto">
      <v-card class="pa-2 pt-4">
        <v-card-title>
          <h2>Register</h2>
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
            Already have an account? <NuxtLink to="/login">Login</NuxtLink>
          </p>
          <v-btn type="submit" color="primary" @click="handleRegister"
            >Register</v-btn
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
