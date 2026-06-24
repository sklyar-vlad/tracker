<template>
  <div class="main-layout">
    <WelcomeHeader />
  </div>
  <div class="auth-container">
    <div class="auth-card">
      <div class="auth-header">
        <h1>Login</h1>
        <p>Sign in to your account</p>
      </div>

      <form @submit.prevent="handleLogin" class="auth-form">
        <div class="form-group">
          <label for="email">Email Address</label>
          <input id="email" v-model="email" type="email" placeholder="you@example.com" required />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="Enter your password"
            required
          />
        </div>

        <button type="submit" class="btn btn-primary btn-full">Sign In</button>
      </form>

      <div class="auth-divider">
        <span>Don't have an account?</span>
      </div>

      <RouterLink to="/register" class="btn btn-secondary btn-full"> Create Account </RouterLink>
    </div>

    <div class="auth-decoration"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'vue-toastification'
import WelcomeHeader from '@/components/Header/WelcomeHeader.vue'

const router = useRouter()
const route = useRoute()
const toast = useToast()

const email = ref('')
const password = ref('')

onMounted(() => {
  if (route.query.registered === 'true') {
    toast.success('Login successfully!')
  }
})

const handleLogin = async () => {
  try {
    const res = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify({
        email: email.value,
        password: password.value,
      }),
    })

    if (!res.ok) {
      throw new Error(await res.text())
    }

    toast.success('Login successful')

    email.value = ''
    password.value = ''

    await router.push('/')
  } catch (err) {
    toast.error(err instanceof Error ? err.message : 'Login failed')
    console.error(err)
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 20px;
}

/* glass card как header/hero */
.auth-card {
  width: 100%;
  max-width: 420px;

  padding: 48px 40px;

  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);

  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);

  border-radius: 16px;

  box-shadow:
    0 15px 60px rgba(0, 0, 0, 0.35),
    0 0 30px rgba(149, 162, 223, 0.08);

  position: relative;
  z-index: 1;
}

/* glow background как hero */
.auth-container::before {
  content: '';
  position: absolute;

  width: 700px;
  height: 700px;

  background: radial-gradient(circle, rgba(149, 162, 223, 0.18), transparent 60%);

  filter: blur(50px);

  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);

  z-index: 0;
}

/* HEADER TEXT */
.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 30px;
  font-weight: 800;

  color: var(--text-primary);

  text-shadow: 0 0 10px rgba(149, 162, 223, 0.2);
}

.auth-header p {
  font-size: 14px;
  color: var(--text-secondary);
}

/* FORM */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
  margin-bottom: 28px;
}

.form-group label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

/* INPUT — glass style */
.form-group input {
  padding: 12px 14px;

  border-radius: 10px;

  border: 1px solid rgba(255, 255, 255, 0.12);

  background: rgba(255, 255, 255, 0.04);

  color: var(--text-primary);

  backdrop-filter: blur(10px);

  transition: all 0.25s ease;
}

.form-group input:focus {
  outline: none;

  border-color: var(--accent-primary);

  box-shadow:
    0 0 0 3px rgba(100, 200, 255, 0.15),
    0 0 20px rgba(149, 162, 223, 0.15);
}

/* PLACEHOLDER */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
/* BUTTON BASE (как в hero) */
.btn {
  width: 100%;

  padding: 12px 24px;
  border-radius: 10px;

  font-size: 15px;
  font-weight: 700;

  display: flex;
  align-items: center;
  justify-content: center;

  text-decoration: none;

  transition: all 0.25s ease;

  font-family: 'Evolventa', sans-serif;
}

/* PRIMARY — neon gradient */
.btn-primary {
  border: none;

  color: white;

  background: linear-gradient(135deg, var(--accent-primary), var(--accent-primary-dark));

  box-shadow: 0 10px 30px rgba(59, 130, 246, 0.25);
}

.btn-primary:hover {
  transform: translateY(-3px);

  box-shadow:
    0 15px 40px rgba(59, 130, 246, 0.35),
    0 0 30px rgba(149, 162, 223, 0.2);
}

/* SECONDARY — glass */
.btn-secondary {
  background: rgba(255, 255, 255, 0.06);

  color: var(--text-primary);

  border: 1px solid rgba(255, 255, 255, 0.12);

  backdrop-filter: blur(10px);
}

.btn-secondary:hover {
  transform: translateY(-3px);

  border-color: var(--accent-primary);

  box-shadow: 0 0 25px rgba(149, 162, 223, 0.15);
}

.btn-full {
  width: 100%;
}

/* divider */
.auth-divider {
  text-align: center;
  margin-bottom: 20px;
  font-size: 13px;
  color: var(--text-secondary);
}

/* RESPONSIVE */
@media (max-width: 640px) {
  .auth-card {
    padding: 32px 24px;
  }

  .auth-header h1 {
    font-size: 24px;
  }
}
</style>
