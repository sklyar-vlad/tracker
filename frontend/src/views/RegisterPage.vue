<template>
  <div class="main-layout">
    <WelcomeHeader />
  </div>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Create Account</h1>
        <p>Start your tracker</p>
      </div>

      <form @submit.prevent="handleRegister" class="auth-form">
        <div class="form-group">
          <input id="name" v-model="username" type="text" required />
          <label for="name">Username</label>
        </div>

        <div class="form-group">
          <input id="email" v-model="email" type="email" required />
          <label for="email">Email Address</label>
        </div>

        <div class="form-group">
          <input
            id="password"
            v-model="password"
            type="password"
            required
            :class="{ 'input-error': passwordError }"
          />
          <label for="password">Password</label>
        </div>

        <div class="form-group">
          <input
            id="confirm-password"
            v-model="confirmPassword"
            type="password"
            required
            :class="{ 'input-error': passwordError }"
          />
          <label for="confirm-password">Confirm Password</label>
        </div>

        <button type="submit" class="btn btn-primary btn-full">Create Account</button>
      </form>

      <div class="auth-divider">
        <span>Already have an account?</span>
      </div>

      <RouterLink to="/login" class="btn btn-secondary btn-full"> Sign In </RouterLink>
    </div>

    <div class="auth-decoration"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import WelcomeHeader from '@/components/Header/WelcomeHeader.vue'

const router = useRouter()
const toast = useToast()

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const passwordError = computed(() => {
  if (!password.value || !confirmPassword.value) return ''

  return password.value !== confirmPassword.value ? 'Passwords do not match' : ''
})

const handleRegister = async () => {
  if (!username.value || !email.value || !password.value || !confirmPassword.value) {
    return
  }

  if (passwordError.value) {
    toast.error(passwordError.value)
    return
  }

  try {
    const res = await fetch('/api/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        email: email.value,
        password: password.value,
      }),
    })

    if (!res.ok) {
      const errText = await res.text()
      throw new Error(errText || 'Registration failed')
    }

    router.push({
      path: '/login',
      query: {
        registered: 'true',
      },
    })
  } catch (err) {
    toast.error(err instanceof Error ? err.message : 'Registration failed')

    console.error('Register error:', err)
  }
}
</script>

<style scoped>
/* =========================
   LOGIN LAYOUT (HERO SYSTEM)
========================= */
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;

  position: relative;
  overflow: hidden;

  padding: 120px 20px 60px;
}

/* glow background */
.auth-container::before {
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
   CARD
========================= */
.auth-card {
  width: 100%;
  max-width: 420px;

  display: flex;
  flex-direction: column;
  gap: 18px;

  padding: clamp(24px, 3vw, 40px);

  background: var(--surface);

  /* 🔥 MAIN FIX */
  border: 1px solid var(--border-subtle);

  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);

  border-radius: 16px;

  box-shadow: var(--shadow-md);

  z-index: 2;
}

/* =========================
   HEADER
========================= */
.auth-header {
  text-align: center;
  margin-bottom: 8px;
}

.auth-header h1 {
  font-size: clamp(24px, 3vw, 36px);
  font-weight: 800;

  color: var(--accent-primary);

  text-shadow:
    0 0 10px rgba(149, 162, 223, 0.15),
    0 0 25px rgba(59, 130, 246, 0.1);
}

.auth-header p {
  font-size: clamp(13px, 1.2vw, 16px);
  color: var(--text-secondary);
}

/* =========================
   FORM
========================= */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  position: relative;
}

.form-group input {
  width: 100%;
  padding: 16px 14px;

  background: var(--surface);

  /* 🔥 FIX */
  border: 1px solid var(--border-default);

  border-radius: 12px;

  color: var(--text-primary);
  font-size: 15px;

  outline: none;

  transition: 0.25s ease;
}

.form-group input:focus {
  border-color: var(--accent-primary);

  box-shadow: 0 0 0 3px var(--border-glow);
}

/* FLOAT LABEL */
.form-group label {
  position: absolute;
  left: 14px;
  top: 50%;

  transform: translateY(-50%);

  color: var(--text-secondary);
  font-size: 14px;

  pointer-events: none;

  transition: 0.2s ease;
  padding: 0 6px;
}

.form-group input:focus ~ label,
.form-group input:valid ~ label {
  top: 0;
  transform: translateY(-50%) scale(0.85);

  color: var(--bg-primary);

  /* вместо hardcoded bg */
  background: var(--accent-primary);

  border-radius: 6px;
}

.input-error {
  border-color: var(--error) !important;
}

.input-error:focus {
  border-color: var(--error) !important;
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15) !important;
}

/* =========================
   DIVIDER
========================= */
.auth-divider {
  text-align: center;
  font-size: 13px;
  color: var(--text-secondary);
  margin: 4px 0;
}

/* =========================
   BUTTONS (SYSTEM UNIFIED)
========================= */
.btn {
  width: 100%;

  padding: 12px 24px;
  border-radius: 10px;

  font-weight: 700;
  text-decoration: none;

  display: flex;
  align-items: center;
  justify-content: center;

  transition: 0.25s ease;

  font-family: 'Evolventa', sans-serif;
}

/* PRIMARY */
.btn-primary {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-dark));

  border: none;
  color: white;

  box-shadow:
    0 10px 30px rgba(59, 130, 246, 0.25),
    0 0 20px rgba(149, 162, 223, 0.15);
}

.btn-primary:hover {
  transform: translateY(-3px);
}

/* SECONDARY */
.btn-secondary {
  background: var(--surface);
  border: 1px solid var(--border-medium);
  color: var(--text-primary);
}

.btn-secondary:hover {
  transform: translateY(-3px);

  border-color: var(--border-strong);

  box-shadow: 0 0 0 1px var(--border-glow);
}

/* =========================
   RESPONSIVE
========================= */
@media (max-width: 768px) {
  .auth-container {
    padding: 100px 16px 40px;
  }

  .auth-card {
    padding: 22px 18px;
  }

  .auth-form {
    width: 100%;
  }
}
</style>
