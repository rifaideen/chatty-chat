import { ElLoading } from "element-plus";

export const useLoading = (text = "Loading", lock = true) => {
  return ElLoading.service({
    lock,
    text,
    background: "rgba(0, 0, 0, 0.7)",
  });
};
