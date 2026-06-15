<template>
  <div class="flex min-h-screen items-center justify-center">
    <div class="w-full max-w-sm space-y-4">
      <h1 class="text-2xl font-bold">MailFyke</h1>
      <p class="text-gray-500">Sign in to your account</p>
      <form @submit.prevent="login" class="space-y-3">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          class="w-full rounded border px-3 py-2"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full rounded border px-3 py-2"
        />
        <button
          type="submit"
          class="w-full rounded bg-blue-600 px-3 py-2 text-white"
        >
          Sign in
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useAuthStore } from "../stores/auth"
import { useRouter } from "vue-router"

const email = ref("")
const password = ref("")
const auth = useAuthStore()
const router = useRouter()

async function login() {
  await auth.login(email.value, password.value)
  router.push("/app/mails")
}
</script>
