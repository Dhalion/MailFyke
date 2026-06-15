<template>
  <div class="flex h-screen">
    <aside class="w-64 border-r bg-white p-4">
      <h2 class="mb-4 text-sm font-semibold uppercase text-gray-400">Mailboxes</h2>
      <ul class="space-y-1">
        <li
          v-for="org in organizations"
          :key="org.id"
          @click="selectOrg(org)"
          class="flex cursor-pointer items-center justify-between rounded px-2 py-1.5 text-sm hover:bg-gray-100"
          :class="{ 'bg-blue-50 font-medium': selectedOrg?.id === org.id }"
        >
          <span>{{ org.name }}</span>
          <span v-if="org.unread_count" class="rounded-full bg-blue-600 px-2 py-0.5 text-xs text-white">
            {{ org.unread_count }}
          </span>
        </li>
      </ul>
      <p v-if="organizations.length === 0" class="text-sm text-gray-400">No organizations</p>
    </aside>
    <main class="flex-1 overflow-y-auto p-6">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useAuthStore } from "../stores/auth"

const auth = useAuthStore()
const organizations = computed(() => auth.organizations)
const selectedOrg = computed(() => auth.selectedOrg)

function selectOrg(org: any) {
  auth.selectedOrg = org
}
</script>
