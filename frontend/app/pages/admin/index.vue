<script setup lang="ts">
definePageMeta({
  layout: 'admin',
})

const config = useRuntimeConfig()

const stats = ref<any>(null)
const classes = ref<any[]>([])
const recentRegistrations = ref<any[]>([])
const loading = ref(true)

async function fetchData() {
  loading.value = true
  try {
    const [s, c, r] = await Promise.allSettled([
      $fetch<any>(`${config.public.apiBase}/stats`),
      $fetch<any>(`${config.public.apiBase}/classes`),
      $fetch<any>(`${config.public.apiBase}/registrations`, { query: { page: 1, limit: 5 } }),
    ])
    if (s.status === 'fulfilled') stats.value = s.value.data
    if (c.status === 'fulfilled') classes.value = c.value.data
    if (r.status === 'fulfilled') recentRegistrations.value = r.value.data?.slice(0, 5) || []
  } finally {
    loading.value = false
  }
}

onMounted(fetchData)

const totalReg = computed(() => stats.value?.total || stats.value?.count || 0)

const capacityPct = computed(() => {
  if (!classes.value.length) return 0
  const total = classes.value.reduce((s, c) => s + c.max_students, 0)
  const used = classes.value.reduce((s, c) => s + (c.current_count || 0), 0)
  return total > 0 ? Math.round((used / total) * 100) : 0
})
</script>

<template>
  <div class="dash">
    <div class="dash-head">
      <h1 class="dash-title">概览</h1>
      <p class="dash-sub">报名系统数据总览</p>
    </div>

    <div v-if="loading" class="dash-loading">
      <UIcon name="i-heroicons-arrow-path" class="w-5 h-5 animate-spin text-blue-500" />
      <span class="text-xs text-gray-500 mt-1">加载中...</span>
    </div>

    <template v-else>
      <div class="kpi-row">
        <div class="kpi">
          <span class="kpi-dot kpi-dot--blue" />
          <div>
            <p class="kpi-label">总报名</p>
            <p class="kpi-val">{{ totalReg }}</p>
          </div>
        </div>
        <div class="kpi">
          <span class="kpi-dot kpi-dot--purple" />
          <div>
            <p class="kpi-label">班级数</p>
            <p class="kpi-val">{{ classes.length }}</p>
          </div>
        </div>
        <div class="kpi">
          <span class="kpi-dot kpi-dot--green" />
          <div>
            <p class="kpi-label">招生进度</p>
            <p class="kpi-val">{{ capacityPct }}%</p>
          </div>
        </div>
        <div class="kpi">
          <span class="kpi-dot kpi-dot--amber" />
          <div>
            <p class="kpi-label">最近报名</p>
            <p class="kpi-val">{{ recentRegistrations.length }} 条</p>
          </div>
        </div>
      </div>

      <div class="panels">
        <div class="panel">
          <div class="panel-head">
            <UIcon name="i-heroicons-academic-cap" class="panel-head-icon" />
            <h3>班级容量</h3>
          </div>
          <div class="panel-list">
            <div v-for="c in classes" :key="c.id" class="panel-item">
              <div class="item-row">
                <span class="item-name">{{ c.name }}</span>
                <span class="item-count">{{ c.current_count || 0 }}/{{ c.max_students }}</span>
              </div>
              <div class="item-bar">
                <div
                  class="item-fill"
                  :class="{ 'item-fill--warn': (c.current_count || 0) / c.max_students > 0.8 }"
                  :style="{ width: c.max_students > 0 ? ((c.current_count || 0) / c.max_students) * 100 + '%' : '0%' }"
                />
              </div>
            </div>
            <div v-if="!classes.length" class="panel-empty">
              <UIcon name="i-heroicons-inbox" class="w-5 h-5 text-gray-400" />
              <span>暂无</span>
            </div>
          </div>
        </div>

        <div class="panel">
          <div class="panel-head">
            <UIcon name="i-heroicons-clock" class="panel-head-icon" />
            <h3>最近报名</h3>
          </div>
          <div class="panel-list">
            <div v-for="r in recentRegistrations" :key="r.id" class="panel-item panel-item--flex">
              <span class="recent-av">{{ r.name?.charAt(0) || '?' }}</span>
              <div>
                <p class="recent-name">{{ r.name }}</p>
                <p class="recent-meta">{{ r.grade || '' }} · {{ new Date(r.registration_time).toLocaleDateString('zh-CN') }}</p>
              </div>
            </div>
            <div v-if="!recentRegistrations.length" class="panel-empty">
              <UIcon name="i-heroicons-inbox" class="w-5 h-5 text-gray-400" />
              <span>暂无</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.dash {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.dash-head {
  margin-bottom: 0;
}

.dash-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #0f172a;
}

.dash-sub {
  font-size: 0.8125rem;
  color: #94a3b8;
  margin-top: 0.125rem;
}

.dash-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 0;
}

.kpi-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.75rem;
}

.kpi {
  background: #fff;
  padding: 0.875rem 1rem;
  border-radius: 0.75rem;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 0.625rem;
}

.kpi-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.kpi-dot--blue   { background: #2563eb; }
.kpi-dot--purple { background: #7c3aed; }
.kpi-dot--green  { background: #059669; }
.kpi-dot--amber  { background: #d97706; }

.kpi-label {
  font-size: 0.75rem;
  color: #64748b;
}

.kpi-val {
  font-size: 1.25rem;
  font-weight: 700;
  color: #0f172a;
  line-height: 1.2;
}

.panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.panel {
  background: #fff;
  border-radius: 0.75rem;
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.panel-head {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.8125rem;
  font-weight: 600;
  color: #1e293b;
}

.panel-head-icon {
  width: 16px;
  height: 16px;
  color: #2563eb;
}

.panel-list {
  padding: 0.5rem 1rem;
}

.panel-item {
  padding: 0.5rem 0;
}

.panel-item + .panel-item {
  border-top: 1px solid #f8fafc;
}

.panel-item--flex {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}

.item-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.375rem;
}

.item-name {
  font-size: 0.8125rem;
  color: #334155;
}

.item-count {
  font-size: 0.75rem;
  color: #94a3b8;
  font-variant-numeric: tabular-nums;
}

.item-bar {
  height: 4px;
  background: #f1f5f9;
  border-radius: 9999px;
}

.item-fill {
  height: 100%;
  border-radius: 9999px;
  background: linear-gradient(90deg, #2563eb, #38bdf8);
  transition: width 0.4s;
}

.item-fill--warn {
  background: linear-gradient(90deg, #d97706, #fbbf24);
}

.recent-av {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
  flex-shrink: 0;
}

.recent-name {
  font-size: 0.8125rem;
  color: #1e293b;
}

.recent-meta {
  font-size: 0.6875rem;
  color: #94a3b8;
  margin-top: 0.0625rem;
}

.panel-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1.5rem 0;
  gap: 0.25rem;
  font-size: 0.8125rem;
  color: #cbd5e1;
}

@media (max-width: 1024px) {
  .kpi-row {
    grid-template-columns: repeat(2, 1fr);
  }
  .panels {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .kpi-row {
    grid-template-columns: repeat(2, 1fr);
    gap: 0.5rem;
  }
  .kpi {
    padding: 0.75rem;
  }
  .kpi-val {
    font-size: 1.125rem;
  }
}
</style>