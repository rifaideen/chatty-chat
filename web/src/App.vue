<script>
import DashboardLayout from "@/layouts/DashboardLayout.vue";
import { provide } from "vue";
import { useNotification } from "./composables/useNotification";
import { useLoading } from "./composables/useLoading";

export default {
  name: "App",
  components: {
    DashboardLayout,
  },
  setup() {
    provide("notification", useNotification);
    provide("loading", useLoading);
    provide("token", localStorage.getItem("token"));
    provide("user", localStorage.getItem("auth-user"));
  },
};
</script>

<template>
  <DashboardLayout>
    <router-view />
  </DashboardLayout>
</template>
