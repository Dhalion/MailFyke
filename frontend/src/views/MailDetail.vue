<template>
  <div v-if="mail">
    <h1 class="text-xl font-bold">{{ mail.subject }}</h1>
    <p class="text-sm text-gray-500">From: {{ mail.mail_from }}</p>
    <p class="text-sm text-gray-500">To: {{ mail.rcpt_to?.join(", ") }}</p>
    <div class="mt-4" v-html="mail.body_html || mail.body_text"></div>
  </div>
  <p v-else class="text-gray-400">Select an email to view</p>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useMailStore } from "../stores/mail"

const mailStore = useMailStore()
const mail = computed(() => mailStore.currentMail)
</script>
