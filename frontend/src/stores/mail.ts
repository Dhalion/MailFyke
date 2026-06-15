import { defineStore } from "pinia"
import { ref } from "vue"

interface Email {
  id: string
  organization_id: string
  mail_from: string
  rcpt_to: string[]
  subject: string
  body_html?: string
  body_text?: string
  raw_eml?: string
  read: boolean
  has_attachments: boolean
  received_at: string
  size_bytes: number
}

export const useMailStore = defineStore("mail", () => {
  const currentMails = ref<Email[]>([])
  const currentMail = ref<Email | null>(null)
  const unreadCounts = ref<Record<string, number>>({})

  return { currentMails, currentMail, unreadCounts }
})
