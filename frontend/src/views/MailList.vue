<template>
  <div>
    <p class="text-gray-500" v-if="!selectedOrg">Select an organization to view emails</p>
    <div v-else>
      <h2 class="text-lg font-semibold">{{ selectedOrg.name }}</h2>
      <ul class="mt-4 space-y-2">
        <li v-for="mail in mails" :key="mail.id" class="rounded border p-3">
          <p class="font-medium">{{ mail.subject }}</p>
          <p class="text-sm text-gray-500">{{ mail.mail_from }}</p>
        </li>
      </ul>
      <p v-if="mails.length === 0" class="mt-4 text-gray-400">No emails yet</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useAuthStore } from "../stores/auth"
import { useMailStore } from "../stores/mail"

const auth = useAuthStore()
const mailStore = useMailStore()

const selectedOrg = computed(() => auth.selectedOrg)
const mails = computed(() => mailStore.currentMails)
</script>
