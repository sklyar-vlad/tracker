<template>
  <div class="main-layout">
    <LandingHeader />
  </div>

  <section class="hero">
    <div class="hero-content">
      <h1>Track everything that matters</h1>
      <p>Simple, fast and clean tracker for your habits, tasks and goals.</p>

      <div class="hero-actions">
        <RouterLink to="/register" class="btn btn-primary"> Get Started </RouterLink>

        <RouterLink to="/login" class="btn btn-secondary"> Sign In </RouterLink>
      </div>
    </div>

    <div class="hero-decoration"></div>
  </section>
</template>

<script setup lang="ts">
import LandingHeader from '@/components/Header/LandingHeader.vue'
</script>

<style scoped>
.hero {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;

  padding: 120px 20px 60px;
}

/* стеклянная карточка как header */
.hero-content {
  max-width: 720px;
  text-align: center;
  z-index: 2;

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;

  padding: 40px 32px;

  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);

  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);

  border-radius: 16px;

  box-shadow: 0 10px 50px rgba(0, 0, 0, 0.25);
}

/* заголовок = тот же стиль акцента что в header .dev */
.hero h1 {
  font-size: 46px;
  font-weight: 800;
  line-height: 1.1;

  color: var(--text-primary);

  text-shadow:
    0 0 10px rgba(149, 162, 223, 0.25),
    0 0 25px rgba(59, 130, 246, 0.15);
}

/* описание мягкое как secondary text */
.hero p {
  font-size: 16px;
  color: var(--text-secondary);
  max-width: 520px;
}

/* кнопки — делаем как "theme vibe" */
.hero-actions {
  display: flex;
  gap: 14px;
  justify-content: center;
  flex-wrap: wrap;

  margin-top: 10px;
}

/* PRIMARY BUTTON — неоновый градиент как underline в nav */
.btn {
  padding: 12px 24px;
  border-radius: 10px;
  font-weight: 700;
  text-decoration: none;

  display: inline-flex;
  align-items: center;
  justify-content: center;

  transition: all 0.25s ease;

  font-family: 'Evolventa', sans-serif;
}

/* PRIMARY */
.btn-primary {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-primary-dark));

  color: white;

  box-shadow:
    0 10px 30px rgba(59, 130, 246, 0.25),
    0 0 20px rgba(149, 162, 223, 0.15);
}

.btn-primary:hover {
  transform: translateY(-3px) scale(1.03);
  box-shadow:
    0 15px 40px rgba(59, 130, 246, 0.35),
    0 0 30px rgba(149, 162, 223, 0.25);
}

/* SECONDARY — glass как header */
.btn-secondary {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-primary);

  border: 1px solid rgba(255, 255, 255, 0.15);

  backdrop-filter: blur(10px);
}

.btn-secondary:hover {
  transform: translateY(-3px);
  border-color: var(--accent-primary);

  box-shadow: 0 0 25px rgba(149, 162, 223, 0.15);
}

/* декоративный glow фон */
.hero::before {
  content: '';
  position: absolute;

  width: 600px;
  height: 600px;

  background: radial-gradient(circle, rgba(149, 162, 223, 0.15), transparent 60%);

  filter: blur(40px);

  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  z-index: 0;
}

/* responsive */
@media (max-width: 900px) {
  .hero h1 {
    font-size: 32px;
  }

  .hero-content {
    padding: 28px 20px;
  }
}
</style>
