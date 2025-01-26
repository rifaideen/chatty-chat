<template>
  <el-container>
    <el-container>
      <el-aside width="200px">
        <slot name="sidebar" />
      </el-aside>
      <el-container>
        <el-main class="main" style="height: 580px; margin-bottom: 20px">
          <slot />
        </el-main>
        <el-footer>
          <slot name="footer" />
        </el-footer>
      </el-container>
    </el-container>
  </el-container>
</template>

<style scoped>
footer {
  position: sticky;
  bottom: 0;
}
</style>
