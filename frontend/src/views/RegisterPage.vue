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
          <label for="name">Username</label>
          <input id="name" v-model="username" type="text" placeholder="Piter Parker" required />
        </div>

        <div class="form-group">
          <label for="email">Email Address</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="example@example.com"
            required
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="Create a password"
            required
          />
        </div>

        <div class="form-group">
          <label for="confirm-password">Confirm Password</label>
          <input
            id="confirm-password"
            v-model="confirmPassword"
            type="password"
            placeholder="Confirm your password"
            required
          />
        </div>

        <div v-if="passwordError" class="error-message">
          {{ passwordError }}
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
  if (password.value && confirmPassword.value && password.value !== confirmPassword.value) {
    return 'Passwords do not match'
  }

  return ''
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
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 20px;
}

.auth-container::before {
  content: '';
  position: absolute;

  width: 700px;
  height: 700px;

  background: radial-gradient(circle, rgba(149, 162, 223, 0.18), transparent 60%);

  filter: blur(60px);

  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.auth-card {
  position: relative;
  z-index: 1;

  width: 100%;
  max-width: 440px;

  padding: 48px 40px;

  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);

  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);

  border-radius: 18px;

  box-shadow:
    0 15px 60px rgba(0, 0, 0, 0.35),
    0 0 30px rgba(149, 162, 223, 0.1);
}

/* HEADER */

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 30px;
  font-weight: 800;

  color: var(--text-primary);

  margin-bottom: 8px;

  text-shadow: 0 0 10px rgba(149, 162, 223, 0.2);
}

.auth-header p {
  color: var(--text-secondary);
  font-size: 14px;
}

/* FORM */

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;

  margin-bottom: 28px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: var(--text-primary);
  font-size: 13px;
  font-weight: 600;
}

.form-group input {
  padding: 13px 15px;

  border-radius: 10px;

  border: 1px solid rgba(255, 255, 255, 0.12);

  background: rgba(255, 255, 255, 0.04);

  color: var(--text-primary);

  transition: 0.25s;
}

.form-group input::placeholder {
  color: var(--text-tertiary);
}

.form-group input:focus {
  outline: none;

  border-color: var(--accent-primary);

  box-shadow:
    0 0 0 3px rgba(100, 200, 255, 0.15),
    0 0 25px rgba(149, 162, 223, 0.15);
}

/* ERROR */

.error-message {
  padding: 12px 16px;

  border-radius: 10px;

  background: rgba(255, 70, 70, 0.08);

  border: 1px solid rgba(255, 70, 70, 0.15);

  color: #ff8d8d;

  font-size: 13px;
}

/* BUTTONS */

.btn {
  padding: 13px 24px;

  border-radius: 10px;

  font-weight: 700;

  text-decoration: none;

  display: flex;
  justify-content: center;
  align-items: center;

  transition: 0.25s;

  font-family: 'Evolventa', sans-serif;
}

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

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);

  color: var(--text-primary);

  border: 1px solid rgba(255, 255, 255, 0.15);

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

/* DIVIDER */

.auth-divider {
  text-align: center;

  color: var(--text-secondary);

  font-size: 13px;

  margin-bottom: 20px;
}

/* MOBILE */

@media (max-width: 640px) {
  .auth-card {
    padding: 32px 24px;
  }

  .auth-header h1 {
    font-size: 24px;
  }
}
</style>
