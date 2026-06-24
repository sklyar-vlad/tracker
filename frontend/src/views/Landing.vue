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
      </div>
    </div>

    <div class="hero-decoration"></div>
  </section>
</template>

<script setup lang="ts">
import LandingHeader from '@/components/Header/LandingHeader.vue'
</script>

<style scoped>
/* =========================
   HERO LAYOUT
========================= */
.hero {
  min-height: 85vh;
  display: flex;
  align-items: center;
  justify-content: center;

  position: relative;

  padding: 120px 20px 60px;
  overflow: hidden;
}

/* =========================
   CONTENT CARD
========================= */
.hero-content {
  max-width: 720px;
  width: 100%;

  text-align: center;

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;

  padding: clamp(24px, 3vw, 40px);

  background: var(--surface);
  border: 1px solid var(--surface-border);

  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);

  border-radius: 16px;

  box-shadow: var(--shadow-md);

  z-index: 2;
}

/* =========================
   TITLE (RESPONSIVE FONTS)
========================= */
.hero h1 {
  font-size: clamp(28px, 4vw, 46px);
  font-weight: 800;
  line-height: 1.1;

  color: var(--text-primary);

  text-shadow:
    0 0 10px rgba(149, 162, 223, 0.18),
    0 0 25px rgba(59, 130, 246, 0.12);
}

/* =========================
   DESCRIPTION
========================= */
.hero p {
  font-size: clamp(14px, 1.2vw, 16px);
  color: var(--text-secondary);

  max-width: 520px;
  line-height: 1.5;
}

/* =========================
   ACTION BUTTONS
========================= */
.hero-actions {
  display: flex;
  gap: 14px;
  justify-content: center;
  flex-wrap: wrap;

  margin-top: 10px;
}

/* shared button base */
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
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-dark));

  color: white;

  box-shadow:
    0 10px 30px rgba(59, 130, 246, 0.25),
    0 0 20px rgba(149, 162, 223, 0.15);
}

.btn-primary:hover {
  transform: translateY(-3px);
  box-shadow:
    0 15px 40px rgba(59, 130, 246, 0.35),
    0 0 30px rgba(149, 162, 223, 0.25);
}

/* SECONDARY */
.btn-secondary {
  background: var(--text-secondary);
  color: var(--text-primary);

  border: 1px solid var(--surface-border);

  backdrop-filter: blur(10px);
}

.btn-secondary:hover {
  transform: translateY(-3px);
  border-color: var(--accent-primary);

  box-shadow: 0 0 25px rgba(149, 162, 223, 0.15);
}

/* =========================
   BACKGROUND GLOW (FIXED LAYERS)
========================= */
.hero::before {
  content: '';
  position: absolute;

  width: 600px;
  height: 600px;

  background: radial-gradient(circle, rgba(149, 162, 223, 0.12), transparent 60%);

  filter: blur(50px);

  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  z-index: 0;
  pointer-events: none;
}

/* =========================
   RESPONSIVE
========================= */
@media (max-width: 768px) {
  .hero {
    padding: 100px 16px 40px;
  }

  .hero-content {
    padding: 22px 18px;
  }

  .hero-actions {
    flex-direction: column;
    width: 100%;
  }

  .btn {
    width: 100%;
  }
}
</style>
