<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { User } from '@/types/user'
import { useUserApi } from '@/composables/useUserApi'

definePageMeta({
  icon: 'mdi-account-group',
  title: 'Users',
  drawerIndex: 2,
})

const users = ref<User[]>([])
const currentPage = ref<number>(1)
const itemsPerPage = ref<number>(10)

const { fetchUsers } = useUserApi()

const router = useRouter()

const goToUser = (id: number) => {
  router.push(`/user/${id}`)
}

const loadUsers = async () => {
  const fetchedUsers = await fetchUsers(itemsPerPage.value, currentPage.value)
  users.value = fetchedUsers
}

onMounted(() => {
  loadUsers()
})

// Watch for page changes
watch(currentPage, () => {
  loadUsers()
})

const nextPage = () => {
  if (users.value.length === itemsPerPage.value) {
    currentPage.value += 1
  }
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value -= 1
  }
}
</script>

<template>
  <v-container>
    <h1 class="text-h4 mb-4">Users</h1>
    <v-list>
      <v-list-item
        v-for="user in users"
        :key="user.id"
        class="cursor-pointer"
        @click="goToUser(user.id)"
      >
        <v-list-item-title>{{ user.username }}</v-list-item-title>
      </v-list-item>
    </v-list>

    <!-- Pagination Controls -->
    <v-row justify="center" class="mt-4">
      <v-btn :disabled="currentPage === 1" class="mr-2" @click="previousPage">
        Previous
      </v-btn>
      <v-btn :disabled="users.length < itemsPerPage" @click="nextPage">
        Next
      </v-btn>
    </v-row>
  </v-container>
</template>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
