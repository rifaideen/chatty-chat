<template>
  <el-row :gutter="20">
    <el-col :span="17">
      <img
        src="/logo.svg"
        width="800"
        height="600"
        alt="logo"
        style="margin-top: -200px"
      />
    </el-col>
    <el-col :span="6">
      <el-form :model="form" :rules="rules" label-position="top" ref="formRef">
        <el-form-item label="Username" prop="username" required>
          <el-input v-model="form.username" placeholder="username" />
        </el-form-item>

        <el-form-item label="Password" prop="password" required>
          <el-input
            v-model="form.password"
            placeholder="password"
            type="password"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="login">Login</el-button>
        </el-form-item>
      </el-form>
    </el-col>
  </el-row>
</template>

<script>
import { inject, ref } from "vue";

export default {
  setup() {
    const notification = inject("notification");
    const loading = inject("loading");

    const formRef = ref(null);
    const form = ref({
      username: "",
      password: "",
    });

    const rules = ref({
      username: [
        { required: true, message: "Username is required", trigger: "blur" },
        {
          min: 3,
          max: 15,
          message: "Length should be 3 to 15 characters",
          trigger: "blur",
        },
      ],
      password: [
        { required: true, message: "Password is required", trigger: "blur" },
        {
          min: 4,
          message: "Password must be at least 6 characters",
          trigger: "blur",
        },
      ],
    });

    const login = async () => {
      await formRef.value.validate(async (valid) => {
        if (valid) {
          const loader = loading("Logging in...");

          const response = await fetch("http://localhost:8001/auth/login", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify(form.value),
          });

          if (!response.ok) {
            notification(
              "Login failed",
              "Please check the credentials and try again.",
              "error",
              5000
            );

            loader.close();
            return false;
          }

          const data = await response.json();

          localStorage.setItem("token", data.token);
          localStorage.setItem("auth-user", data.user);
          loader.close();

          form.value.username = "";
          form.value.password = "";

          notification("Login successful", "Welcome back!", "success", 5000);

          setTimeout(() => {
            window.location.reload();
          }, 1000);
        } else {
          notification(
            "Validation Error",
            "Please check the credentials and try again.",
            "error",
            5000
          );

          return false;
        }
      });
    };

    return {
      form,
      formRef,
      login,
      rules,
    };
  },
};
</script>

<style lang="scss" scoped></style>
