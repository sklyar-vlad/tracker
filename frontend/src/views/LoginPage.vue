<template>
  <div class="main-layout">
    <Header />
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
import Header from '@/components/HeaderAuth.vue'

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
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.auth-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  padding: 48px 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: var(--shadow-lg);
  position: relative;
  z-index: 1;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
}

.auth-header p {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 28px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.form-group input {
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  transition: all 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent-primary);
  box-shadow: 0 0 0 3px var(--accent-primary-light);
}

.form-group input::placeholder {
  color: var(--text-tertiary);
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-family: inherit;
}

.btn-primary {
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-primary-dark));
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(100, 200, 255, 0.3);
}

.btn-primary:active {
  transform: translateY(0);
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--accent-primary);
  border: 2px solid var(--accent-primary);
}

.btn-secondary:hover {
  background: var(--accent-primary);
  color: white;
  transform: translateY(-2px);
}

.btn-full {
  width: 100%;
}

.auth-divider {
  text-align: center;
  margin-bottom: 20px;
  font-size: 13px;
  color: var(--text-secondary);
}

@media (max-width: 640px) {
  .auth-card {
    padding: 32px 24px;
  }

  .auth-header h1 {
    font-size: 24px;
  }
}
</style>
