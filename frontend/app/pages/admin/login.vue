<script setup lang="ts">
definePageMeta({
  layout: 'empty',
  middleware: 'guest',
})

const { login, loading, error, username } = useAuth()
const router = useRouter()

const form = reactive({
  username: '',
  password: '',
})

async function handleSubmit() {
  if (!form.username || !form.password) return

  await login(form.username, form.password)

  if (!error.value && username.value) {
    router.push('/admin')
  }
}
</script>

<template>
  <div class="login-wrapper">
    <div class="login-bg-shapes">
      <div class="shape shape-1" />
      <div class="shape shape-2" />
      <div class="shape shape-3" />
    </div>

    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="login-logo">
            <UIcon name="i-heroicons-shield-check" class="w-10 h-10 text-white" />
          </div>
          <h1 class="login-title">管理控制台</h1>
          <p class="login-subtitle">Admin Console</p>
        </div>

        <form class="login-form" @submit.prevent="handleSubmit">
          <div class="form-field">
            <label class="field-label">
              <UIcon name="i-heroicons-user" class="field-label-icon" />
              用户名
            </label>
            <input
              v-model="form.username"
              type="text"
              class="input-base"
              placeholder="请输入用户名"
              autocomplete="username"
            >
          </div>

          <div class="form-field">
            <label class="field-label">
              <UIcon name="i-heroicons-lock-closed" class="field-label-icon" />
              密码
            </label>
            <input
              v-model="form.password"
              type="password"
              class="input-base"
              placeholder="请输入密码"
              autocomplete="current-password"
            >
          </div>

          <div
            v-if="error"
            class="login-error"
          >
            <UIcon name="i-heroicons-exclamation-circle" class="w-5 h-5" />
            <span>{{ error }}</span>
          </div>

          <button
            type="submit"
            class="btn-pill login-submit-btn"
            :disabled="loading"
          >
            <UIcon
              v-if="loading"
              name="i-heroicons-arrow-path"
              class="w-4 h-4 animate-spin"
            />
            <span>{{ loading ? '登录中...' : '登 录' }}</span>
          </button>
        </form>

        <div class="login-footer">
          <NuxtLink to="/" class="back-link">
            <UIcon name="i-heroicons-arrow-left" class="w-4 h-4" />
            返回报名页面
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #172554 0%, #1e3a8a 40%, #2563eb 70%, #38bdf8 100%);
  position: relative;
  overflow: hidden;
}

.login-bg-shapes {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.08;
  background: #ffffff;
}

.shape-1 {
  width: 400px;
  height: 400px;
  top: -100px;
  right: -100px;
}

.shape-2 {
  width: 300px;
  height: 300px;
  bottom: -80px;
  left: -80px;
}

.shape-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 60%;
  transform: translate(-50%, -50%);
}

.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  padding: 1.5rem;
}

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(24px);
  border-radius: 1.25rem;
  box-shadow: 0 24px 48px -12px rgba(0, 0, 0, 0.25);
  padding: 2.5rem 2rem;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-logo {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem;
  box-shadow: 0 8px 24px -4px rgba(37, 99, 235, 0.3);
}

.login-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1e293b;
  line-height: 1.2;
}

.login-subtitle {
  font-size: 0.875rem;
  color: #94a3b8;
  margin-top: 0.25rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 0.4375rem;
}

.field-label-icon {
  width: 15px;
  height: 15px;
  color: #94a3b8;
}

.login-submit-btn {
  width: 100%;
  padding: 0.625rem 1.5rem;
  font-size: 0.9375rem;
  margin-top: 0.5rem;
}

.login-submit-btn .animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.login-error {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 0.75rem;
  font-size: 0.875rem;
}

.login-error svg {
  flex-shrink: 0;
}

.login-footer {
  text-align: center;
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  color: #64748b;
  font-size: 0.875rem;
  text-decoration: none;
  transition: color 0.15s ease;
}

.back-link:hover {
  color: #2563eb;
}

@media (max-width: 480px) {
  .login-card {
    padding: 2rem 1.5rem;
  }

  .login-title {
    font-size: 1.25rem;
  }
}
</style>