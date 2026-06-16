<template>
  <div class="page">
    <div class="card">
      <h1>Hi, its my first site</h1>
      <h2>I'm try to input this frontend in docker</h2>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const count = ref<number>(0)
</script>

<style>
.page {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: var(--background);
}
.card {
  color: #03451e;
}
</style>