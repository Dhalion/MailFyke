import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      name: "login",
      component: () => import("../views/Login.vue"),
    },
    {
      path: "/app",
      component: () => import("../layouts/AppLayout.vue"),
      children: [
        {
          path: "mails",
          name: "mails",
          component: () => import("../views/MailList.vue"),
        },
        {
          path: "mails/:id",
          name: "mail-detail",
          component: () => import("../views/MailDetail.vue"),
        },
      ],
    },
  ],
})

export default router
