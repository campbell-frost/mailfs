import { createRouter, createWebHistory } from "vue-router";
import { useAuth } from "@/lib/auth";

import Home from "@/pages/home.vue";
import Login from "@/pages/login.vue";

const routes = [
  { path: "/", name: "home", component: Home },
  { path: "/login", name: "login", component: Login },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});

// gate frontend routes
router.beforeEach(async to => {
  const { user, loaded } = useAuth();
  await loaded;

  if (!user.value && to.name !== "login") return { name: "login" };
  if (user.value && to.name === "login") return { name: "home" };
});
