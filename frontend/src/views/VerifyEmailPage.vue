<template>
  <div class="main-layout">
    <WelcomeHeader />
  </div>

  <div class="auth-container">
    <div class="auth-card">

      <!-- LOADING -->
      <div v-if="loading" class="auth-header">
        <h1>Verifying Email</h1>
        <p>Please wait while we verify your email address...</p>

        <div class="loader"></div>
      </div>

      <!-- SUCCESS -->
      <div v-else-if="success" class="auth-header">
        <div class="status-icon success">
          ✓
        </div>

        <h1>Email Verified</h1>

        <p>
          Your email has been successfully verified.
          You can now sign in.
        </p>

        <RouterLink to="/login" class="btn btn-primary btn-full">
          Go to Login
        </RouterLink>
      </div>

      <!-- ERROR -->
      <div v-else class="auth-header">
        <div class="status-icon error">
          ✕
        </div>

        <h1>Verification Failed</h1>

        <p>{{ error }}</p>

        <RouterLink to="/login" class="btn btn-secondary btn-full">
          Back to Login
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import WelcomeHeader from '@/components/Header/WelcomeHeader.vue'

const route = useRoute()

const loading = ref(true)
const success = ref(false)
const error = ref('Unable to verify your email.')

onMounted(async () => {
  const token = route.params.token

  if (!token || typeof token !== 'string') {
    loading.value = false
    error.value = 'Verification token is missing.'
    return
  }

  try {
    const res = await fetch(
      `/api/verify/${encodeURIComponent(token)}`,
      {
        method: 'POST',
      }
    )

    if (!res.ok) {
      const text = await res.text()
      throw new Error(text || 'Verification failed')
    }

    success.value = true
  } catch (e) {
    error.value =
      e instanceof Error
        ? e.message
        : 'Verification failed.'
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 120px 20px 60px;
  position: relative;
}

.auth-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
  background: var(--surface);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  backdrop-filter: blur(14px);
  box-shadow: var(--shadow-md);
  text-align: center;
}

.auth-header {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.auth-header h1 {
  font-size: 32px;
  color: var(--accent-primary);
}

.auth-header p {
  color: var(--text-secondary);
  line-height: 1.5;
}

.status-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 34px;
  font-weight: bold;
}

.success {
  background: rgba(34,197,94,.12);
  color: #22c55e;
}

.error {
  background: rgba(239,68,68,.12);
  color: #ef4444;
}

.loader {
  width: 46px;
  height: 46px;
  margin: 12px auto 0;
  border: 4px solid rgba(149,162,223,.2);
  border-top-color: var(--accent-primary);
  border-radius: 50%;
  animation: spin .8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.btn {
  width: 100%;
  padding: 12px 24px;
  border-radius: 10px;
  font-weight: 700;
  text-decoration: none;
  display: flex;
  justify-content: center;
  align-items: center;
  transition: .25s ease;
}

.btn-primary {
  color: white;
  border: none;
  background: linear-gradient(135deg, var(--accent-primary), var(--accent-dark));
}

.btn-secondary {
  background: var(--surface);
  border: 1px solid var(--border-medium);
  color: var(--text-primary);
}

.btn-primary:hover,
.btn-secondary:hover {
  transform: translateY(-3px);
}
</style>