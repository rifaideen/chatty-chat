import { ElNotification } from "element-plus";

export const useNotification = (title, message, type, duration = 3000) => {
  return ElNotification({
    message,
    title,
    type,
    duration,
  });
};
