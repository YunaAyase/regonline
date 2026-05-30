<script setup lang="ts">
const { logout } = useAuth()
const route = useRoute()
const { siteName, fetchSiteSettings } = useSiteSettings()

const now = ref(new Date())

const timeStr = computed(() =>
  now.value.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' }),
)

const dateStr = computed(() =>
  now.value.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric', weekday: 'short' }),
)

const navItems = [
  { to: '/admin', label: '概览', icon: 'i-heroicons-squares-2x2' },
  { to: '/admin/registrations', label: '报名列表', icon: 'i-heroicons-list-bullet' },
  { to: '/admin/classes', label: '班级管理', icon: 'i-heroicons-academic-cap' },
  { to: '/admin/settings', label: '网站设置', icon: 'i-heroicons-cog-6-tooth' },
]

let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  fetchSiteSettings()
  timer = setInterval(() => {
    now.value = new Date()
  }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="admin-app">
    <header class="topbar">
      <div class="topbar-inner">
        <div class="topbar-brand">
          <span class="brand-icon-wrap">
            <UIcon name="i-heroicons-academic-cap" class="brand-icon-svg" />
          </span>
          <div class="brand-meta">
            <span class="brand-name">{{ siteName }}</span>
            <span class="brand-tagline">Admin Console</span>
          </div>
        </div>

        <nav class="topbar-nav">
          <NuxtLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            class="nav-pill"
            :class="{ 'nav-pill--active': route.path === item.to }"
          >
            <UIcon :name="item.icon" class="nav-pill-icon" />
            <span class="nav-pill-text">{{ item.label }}</span>
          </NuxtLink>
        </nav>

        <div class="topbar-extra">
          <div class="datetime">
            <span class="datetime-time">{{ timeStr }}</span>
            <span class="datetime-date">{{ dateStr }}</span>
          </div>
          <button class="nav-logout-btn" @click="logout">
            <UIcon name="i-heroicons-arrow-right-on-rectangle" class="nav-logout-icon" />
            <span>退出</span>
          </button>
        </div>
      </div>
    </header>

    <main class="main-area">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.admin-app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f1f5f9;
}

.topbar {
  height: 56px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #0f172a 0%, #1e3a8a 40%, #2563eb 70%, #38bdf8 100%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 2px 12px rgba(15, 23, 42, 0.2);
}

.topbar-inner {
  height: 100%;
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1.25rem;
  display: flex;
  align-items: center;
}

.topbar-brand {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  flex-shrink: 0;
  margin-right: 1.25rem;
}

.brand-icon-wrap {
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-icon-svg {
  width: 18px;
  height: 18px;
  color: #ffffff;
}

.brand-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.15;
}

.brand-name {
  font-size: 0.9375rem;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
}

.brand-tagline {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.5);
  font-weight: 400;
}

.topbar-nav {
  display: flex;
  align-items: center;
  gap: 0.1875rem;
  flex: 1;
  justify-content: center;
}

.nav-pill {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.4375rem 0.875rem;
  border-radius: 9999px;
  color: rgba(255, 255, 255, 0.68);
  text-decoration: none;
  font-size: 0.8125rem;
  font-weight: 500;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  position: relative;
}

.nav-pill:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
}

.nav-pill--active {
  background: rgba(255, 255, 255, 0.18);
  color: #ffffff;
  box-shadow: 0 0 12px -2px rgba(56, 189, 248, 0.15);
}

.nav-pill-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.nav-pill-text {
  user-select: none;
}

.topbar-extra {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-shrink: 0;
  margin-left: 1.25rem;
}

.datetime {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  line-height: 1.15;
}

.datetime-time {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #ffffff;
  font-variant-numeric: tabular-nums;
}

.datetime-date {
  font-size: 0.625rem;
  color: rgba(255, 255, 255, 0.5);
  white-space: nowrap;
}

.nav-logout-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.4375rem 1rem;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 9999px;
  background: transparent;
  color: #ffffff;
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
}

.nav-logout-btn:hover {
  background: rgba(255, 255, 255, 0.12);
  border-color: rgba(255, 255, 255, 0.5);
}

.nav-logout-icon {
  width: 15px;
  height: 15px;
}

.main-area {
  flex: 1;
  padding: 1.25rem;
  max-width: 1280px;
  width: 100%;
  margin: 0 auto;
}

@media (max-width: 1024px) {
  .nav-pill-text {
    display: none;
  }

  .nav-pill {
    padding: 0.4375rem 0.625rem;
  }

  .topbar-brand {
    margin-right: 0.75rem;
  }

  .topbar-extra {
    margin-left: 0.75rem;
  }
}

@media (max-width: 768px) {
  .topbar {
    height: auto;
  }

  .topbar-inner {
    flex-wrap: wrap;
    gap: 0.375rem;
    padding: 0.5rem 0.75rem;
  }

  .topbar-nav {
    order: 3;
    width: 100%;
    justify-content: flex-start;
    overflow-x: auto;
    padding-bottom: 0.25rem;
  }

  .datetime {
    display: none;
  }

  .main-area {
    padding: 1rem;
  }
}
</style>