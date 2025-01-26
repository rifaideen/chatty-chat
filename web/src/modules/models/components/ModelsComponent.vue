<template>
  <ChatLayout>
    <el-row v-if="!chats.length">
      <el-col>
        <h1 class="title">Hello, {{ user }}</h1>
        <h1 class="sub-title">
          Here for first time? Pull the AI model you want to chat with.
        </h1>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="18">
        <el-row v-for="(chat, i) in chats" :key="i">
          <el-col v-if="chat.role == 'user'" :span="18" :push="6" class="user">
            {{ chat.content }}
          </el-col>
          <el-col v-else :span="22" class="assistant">
            <div
              v-if="chat.content != ''"
              v-html="md.render(chat.content)"
            ></div>
          </el-col>
        </el-row>
        <el-skeleton
          v-show="loading"
          style="margin-top: 10px"
          :rows="2"
          animated
        />
      </el-col>
    </el-row>

    <template #footer>
      <el-row>
        <el-col v-for="model of models" :key="model.model" :span="6">
          <el-button type="primary" size="large" @click="pull(model.model)"
            >Pull {{ model.name }}</el-button
          >
        </el-col>
      </el-row>
    </template>
  </ChatLayout>
</template>

<script>
import markdownit from "markdown-it";
import hljs from "highlight.js";
import ChatLayout from "@/layouts/ChatLayout.vue";
import { inject, ref, onMounted, watch } from "vue";
import { onBeforeRouteLeave } from "vue-router";

export default {
  components: {
    ChatLayout,
  },
  setup() {
    const notification = inject("notification");
    const user = inject("user");
    const token = inject("token");

    const loading = ref(false);
    const receiving = ref(false);
    const response = ref("");
    const chats = ref([]);
    const ws = ref(null);

    const models = ref([
      {
        model: "phi",
        name: "Phi",
      },
      {
        model: "llama3.2",
        name: "llama3.2",
      },
      {
        model: "gemma",
        name: "Gemma",
      },
    ]);

    const stopLoading = () => {
      loading.value = false;
    };

    const md = markdownit({
      html: true,
      linkify: true,
      typographer: true,
      highlight: function (str, lang) {
        if (lang && hljs.getLanguage(lang)) {
          try {
            return hljs.highlight(str, { language: lang }).value;
          } catch (__) {}
        }

        return ""; // use external default escaping
      },
    });

    const connect = () => {
      ws.value = new WebSocket(`ws://localhost:8003/ws?token=${token}`);

      ws.value.onopen = () => {
        receiving.value = false;
      };

      ws.value.onmessage = async (event) => {
        stopLoading();

        try {
          const { type, data, done } = JSON.parse(event.data);

          if (type == "notification") {
            notification("New Message", data, "success");
          } else if (type == "pull") {
            response.value +=
              data.status +
              " [downloading " +
              data.completed +
              " of total " +
              data.total +
              "]\n\n";
            receiving.value = !done;
          }
        } catch (error) {
          console.warn("Error parsing message", error);
        }
      };

      ws.value.onclose = () => {
        console.log("WebSocket closed");
      };

      ws.value.onerror = (error) => {
        stopLoading();
        console.error("WebSocket error:", error);
        notification(
          "Connection Error",
          "Oops! Something went wrong. Please try again later.",
          "error"
        );
      };
    };

    const pull = (model) => {
      if (
        model &&
        ws.value &&
        ws.value.readyState == WebSocket.OPEN &&
        !receiving.value
      ) {
        loading.value = true;

        chats.value.push(
          {
            role: "user",
            content: `pull ${model}`,
          },
          {
            role: "assistant",
            content: "",
          }
        );

        const payload = {
          type: "pull",
          model,
        };

        ws.value.send(JSON.stringify(payload));
        response.value = "";
      }
    };

    watch(
      () => response.value,
      (v) => {
        if (v != "") {
          const last = chats.value.length - 1;
          const chat = chats.value[last];

          // Make sure the last chat's role is 'assistant'
          if (chat && chat.role == "assistant") {
            chat.content = v;
          }
        }
      }
    );

    onMounted(() => {
      connect();
    });

    onBeforeRouteLeave(() => {
      if (ws.value) {
        ws.value.close();
      }
    });

    return {
      user,
      chats,
      loading,
      receiving,
      pull,
      md,
      models,
    };
  },
};
</script>

<style scoped>
.input-area {
  padding: 10px;
  border-top: 1px solid #e0e0e0;
  background-color: #ffffff;
  box-shadow: 0 -2px 5px rgba(0, 0, 0, 0.1);
  border-radius: 5px;
}

.user {
  padding: 10px;
  background: var(--el-color-primary-light-9);
  margin-bottom: 10px;
}

.assistant {
  padding: 10px;
  background: var(--el-color-info-light-9);
  margin-bottom: 10px;
  line-height: var(--el-font-line-height-primary);
  font-size: var(--el-font-size-medium);
}

.form {
  position: relative;
  bottom: 10%;
  width: 100%;
}

h1.title {
  font-size: 56px;
  background: -webkit-linear-gradient(#5082ee, #d36679);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
h1.sub-title {
  font-size: 56px;
  color: var(--el-color-info);
}

.input-area {
  padding: 10px;
  border-top: 1px solid #e0e0e0;
  background-color: #ffffff;
  box-shadow: 0 -2px 5px rgba(0, 0, 0, 0.1);
  border-radius: 5px;
  width: 75%;
}
.disclaimer {
  font-size: 14px;
  color: var(--el-color-info);
}
</style>
