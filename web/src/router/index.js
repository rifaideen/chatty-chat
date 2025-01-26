import { createRouter, createWebHashHistory } from "vue-router";

import HomeView from "@/components/HomeComponent.vue";
import NotFound from "@/components/NotFound.vue";
import chatRoute from "@/modules/chat/chat.route";
import modelsRoute from "@/modules/models/models.route";

const routes = [
  { path: "/", component: HomeView },
  {
    path: "/login",
    component: () => import("@/components/LoginComponent.vue"),
  },
  ...chatRoute,
  ...modelsRoute,
  { path: "/:pathMatch(.*)*", name: "NotFound", component: NotFound },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.path != "/login" && !localStorage.getItem("token")) {
    next("/login");
  } else if (to.path == "/login" && localStorage.getItem("token")) {
    next("/");
  } else {
    next();
  }
});
export default router;
