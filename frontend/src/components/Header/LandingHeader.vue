<template>
  <header class="header">
    <!-- LEFT -->
    <div class="left">
      <RouterLink to="/" class="logo-link">
        <span class="self">{Self}</span>
        <span class="slash">/</span>
        <span class="dev">Dev</span>
      </RouterLink>
    </div>

    <!-- CENTER -->
    <nav class="nav">
      <a href="#">About</a>
      <a href="#">Features</a>
      <a href="#">Tech</a>
      <a href="#">Faq</a>
    </nav>

    <!-- RIGHT -->
    <div class="right">
      <button class="theme-btn" @click="toggleTheme">
        <img :src="theme === 'dark' ? '/icons/sun.svg' : '/icons/moon.svg'" />
      </button>

      <RouterLink to="/login" class="login">
        <span class="login">Log in<img :src="theme === 'dark' ? '/icons/arrowWhite.svg' : '/icons/arrowDark.svg'" /></span>
      </RouterLink>
    </div>
  </header>
</template>

<script setup>
import { useTheme } from '@/composables/useTheme'

const { theme, toggleTheme } = useTheme()
</script>

<style scoped>
/* =========================
   HEADER CONTAINER
========================= */
.header {
  position: fixed;
  top: 10px;
  left: 14px;
  right: 14px;
  z-index: 9999;

  height: 70px;

  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 0 24px;

  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);

  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);

  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.2);
  border-radius: 10px;
}

/* =========================
   LEFT (LOGO)
========================= */
.logo-link {
  font-family: 'Montserrat', sans-serif;
  transition: transform 0.25s ease;
  font-size: clamp(20px, 2vw, 32px);
  font-weight: 700;

  display: flex;
  align-items: center;

  text-decoration: none;
}

.logo-link:hover {
  transform: scale(1.03);
}

.self {
  color: var(--text-primary);
}

.slash {
  color: var(--accent-primary);
  margin: 0 4px;
}

.dev {
  color: var(--accent-primary);
}

/* =========================
   NAV (DESKTOP)
========================= */
.nav {
  display: flex;
  gap: clamp(20px, 4vw, 57px);
  justify-content: center;
  flex: 1;
}

.nav a {
  font-family: 'Hind Madurai', sans-serif;

  font-size: clamp(14px, 1.2vw, 20px);
  font-weight: 700;

  color: var(--text-primary);
  text-decoration: none;

  position: relative;
  transition: 0.3s;
}

.nav a:hover {
  color: var(--accent-primary);
  transform: translateY(-1px);
}

.nav a::after {
  content: '';
  position: absolute;

  bottom: -6px;
  left: 50%;

  width: 0;
  height: 2px;

  background: var(--accent-primary);

  transform: translateX(-50%);
  transition: 0.3s ease;
}

.nav a:hover::after {
  width: 100%;
}

/* =========================
   RIGHT
========================= */
.right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.login {
  font-family: 'Hind Madurai', sans-serif;

  font-size: clamp(16px, 1.2vw, 24px);
  font-weight: 700;

  color: var(--text-primary);
  text-decoration: none;

  display: flex;
  align-items: center;
  gap: 6px;
}

.login-arrow {
  transition: 0.3s;
  opacity: 0.7;
}

.login:hover .login-arrow {
  transform: translateX(4px);
  opacity: 1;
}

/* =========================
   THEME BUTTON
========================= */
.theme-btn {
  width: 56px;
  height: 56px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: transparent;

  border: 1px solid var(--surface-border);
  border-radius: 10px;
  border-color: transparent;

  cursor: pointer;
  transition: all 0.25s ease;
}

.theme-btn:hover {
  transform: translateY(-2px);

  border-color: var(--accent-primary);

  box-shadow:
    0 10px 30px rgba(0, 0, 0, 0.25),
    0 0 20px rgba(100, 200, 255, 0.15);
}

/* =========================
   MOBILE (MAIN FIX)
========================= */
@media (max-width: 768px) {
  .header {
    padding: 0 16px;
    height: 64px;
  }

  .nav {
    display: none;
  }

  .logo-link {
    font-size: 18px;
  }

  .theme-btn {
    width: 44px;
    height: 44px;
  }

  .login {
    font-size: 16px;
  }
}
</style>
