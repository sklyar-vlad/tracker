<template>
  <header class="header">
    <div class="logo">
      <RouterLink to="/" class="logo-link">
        <span class="self">{Self}</span>
        <span class="slash">/</span>
        <span class="dev">Dev</span>
      </RouterLink>
    </div>
  </header>
</template>

<script setup></script>

<style scoped>
.header {
  position: fixed;
  top: 10px;
  left: 10px;
  right: 10px;
  z-index: 99999;

  height: 60px;

  display: flex;
  align-items: center;

  padding: 0 24px;

  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.12);

  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(14px);

  box-shadow: 0px 10px 40px rgba(0, 0, 0, 0.25);
  border-radius: 12px;

  animation: headerAppear 0.8s ease;
}

.logo-link {
  text-decoration: none;
  display: flex;
  align-items: center;
  font-size: 36px;
  font-weight: 700;

  transition: all 0.3s ease;
}

.logo-link:hover {
  transform: translateY(-1px) scale(1.02);
}

.self {
  color: #d6dde0;
  transition: 0.3s;
}

.slash {
  color: #d6dde0;
  margin: 0 4px;

  animation: slashGlow 3s infinite ease-in-out;
}

.dev {
  color: #95a2df;
  font-weight: 800;
  transition: 0.3s;

  text-shadow:
    0 0 10px rgba(149, 162, 223, 0.3),
    0 0 20px rgba(149, 162, 223, 0.2);
}

.logo-link:hover .dev {
  color: #b8c3ff;

  text-shadow:
    0 0 15px rgba(149, 162, 223, 0.8),
    0 0 40px rgba(149, 162, 223, 0.4);
}

@keyframes headerAppear {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slashGlow {
  0%,
  100% {
    opacity: 0.6;
  }

  50% {
    opacity: 1;
    text-shadow:
      0 0 10px #95a2df,
      0 0 25px #95a2df;
  }
}

@media (max-width: 768px) {
  .nav {
    display: none;
  }

  .logo-link {
    font-size: 24px;
  }
}
</style>
