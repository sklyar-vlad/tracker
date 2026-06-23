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

  animation: headerAppear .8s ease;
}


.logo-link {
  text-decoration: none;
  display: flex;
  align-items: center;
  font-size: 36px;
  font-weight: 700;

  transition: all .3s ease;
}

.logo-link:hover {
  transform: translateY(-1px) scale(1.02);
}

.self {
  color: #D6DDE0;
  transition: .3s;
}

.slash {
  color: #D6DDE0;
  margin: 0 4px;

  animation: slashGlow 3s infinite ease-in-out;
}

.dev {
  color: #95A2DF;
  font-weight: 800;
  transition: .3s;

  text-shadow:
    0 0 10px rgba(149,162,223,.3),
    0 0 20px rgba(149,162,223,.2);
}

.logo-link:hover .dev {
  color: #b8c3ff;

  text-shadow:
    0 0 15px rgba(149,162,223,.8),
    0 0 40px rgba(149,162,223,.4);
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
  0%,100% {
    opacity: .6;
  }

  50% {
    opacity: 1;
    text-shadow:
      0 0 10px #95A2DF,
      0 0 25px #95A2DF;
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