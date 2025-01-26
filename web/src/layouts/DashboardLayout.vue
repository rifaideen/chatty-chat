<template>
  <el-container>
    <el-header>
      <el-menu mode="horizontal" :ellipsis="false" router>
        <el-menu-item index="/">
          <img style="width: 100px" src="/logo.svg" />
        </el-menu-item>
        <el-menu-item index="/login" v-if="!isAuthenticated">
          Login
        </el-menu-item>
        <el-menu-item v-if="isAuthenticated" index="/chat/">Chat</el-menu-item>
        <el-menu-item v-if="isAuthenticated" index="/models"
          >Models</el-menu-item
        >
        <el-menu-item v-if="isAuthenticated" @click="logout">
          Logout
        </el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <slot />
    </el-main>
  </el-container>
</template>

<script>
import { computed } from "vue";

export default {
  setup() {
    const isAuthenticated = computed(
      () => localStorage.getItem("token") !== null
    );

    const logout = () => {
      localStorage.removeItem("token");
      window.location.reload();
    };

    return {
      isAuthenticated,
      logout,
    };
  },
};
</script>

<style>
.el-menu--horizontal > .el-menu-item:nth-child(1) {
  margin-right: auto;
}
</style>
