<template>
  <el-container>
    <el-main>
      <el-carousel
        height="500px"
        autoplay
        trigger="hover"
        indicator-position="outside"
      >
        <el-carousel-item v-for="(image, index) in images" :key="index">
          <img
            :src="image.src"
            :alt="image.alt"
            class="carousel-image"
            height="500"
          />
        </el-carousel-item>
      </el-carousel>
    </el-main>
  </el-container>
</template>

<script>
import { inject, computed } from "vue";

export default {
  name: "HomePage",
  setup() {
    const user = inject("user");

    const images = computed(() => {
      return [
        {
          src: `/logo-grayscale.png`,
          alt: "Logo",
        },
        {
          src: `https://placehold.co/800x500?text=Hello+${user || "User"}`,
          alt: "Image 1",
        },
        {
          src: "https://placehold.co/800x500?text=Welcome+to+Chitty+Chat.",
          alt: "Image 2",
        },
        {
          src: "https://placehold.co/800x500?text=Let's+chat.",
          alt: "Image 3",
        },
        {
          src: "https://placehold.co/800x500?text=Still+havent+decided+what+to+do?",
          alt: "Image 4",
        },
        {
          src: "https://placehold.co/800x500?text=Let+our+AI+assistant+help+you+out.",
          alt: "Image 5",
        },
      ];
    });

    return {
      images,
    };
  },
};
</script>

<style scoped>
.carousel-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
h1 {
  text-align: center;
  margin: 0;
  padding: 1rem 0;
  color: #333;
}
</style>
