<template>
  <div class="main-layout">
    <Header />
  </div>
  
  <section class="hero">
    <div class="hero-content">
      <h1>Track everything that matters</h1>
      <p>
        Simple, fast and clean tracker for your habits, tasks and goals.
      </p>

      <div class="hero-actions">
        <RouterLink to="/register" class="btn btn-primary">
          Get Started
        </RouterLink>

        <RouterLink to="/login" class="btn btn-secondary">
          Sign In
        </RouterLink>
      </div>
    </div>

    <div class="hero-decoration"></div>
  </section>
</template>

<script setup lang="ts">
import Header from '@/components/HeaderAuth.vue'
</script>

<style scoped>

.hero {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.hero h1 {
  font-size: 44px;
  font-weight: 800;
  color: var(--text-primary);
  margin-bottom: 16px;
  line-height: 1.1;
}

.hero p {
  font-size: 16px;
  color: var(--text-secondary);
  margin-bottom: 28px;
}

.hero-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;

  margin-top: 8px;
}

.hero-content {
  max-width: 720px;
  text-align: center;
  z-index: 2;

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}


/* buttons — ты уже используешь эти классы */
.btn {
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: 600;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(
    135deg,
    var(--accent-primary),
    var(--accent-primary-dark)
  );
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--accent-primary);
  border: 2px solid var(--accent-primary);
}

.btn-secondary:hover {
  background: var(--accent-primary);
  color: white;
}

/* responsive */
@media (max-width: 900px) {
  .features {
    grid-template-columns: 1fr;
  }

  .hero h1 {
    font-size: 32px;
  }
}
</style>