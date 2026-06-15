import { defineStore } from "pinia"
import { ref } from "vue"

interface Organization {
  id: string
  name: string
  slug: string
  role: string
  unread_count: number
}

interface User {
  id: string
  email: string
  is_admin: boolean
}

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(null)
  const user = ref<User | null>(null)
  const organizations = ref<Organization[]>([])
  const selectedOrg = ref<Organization | null>(null)

  async function login(email: string, password: string) {
    // TODO: implement API call
    token.value = "placeholder"
  }

  function logout() {
    token.value = null
    user.value = null
    organizations.value = []
    selectedOrg.value = null
  }

  return { token, user, organizations, selectedOrg, login, logout }
})
