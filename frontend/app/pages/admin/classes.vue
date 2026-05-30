<script setup lang="ts">
import type { ClassItem } from '~/types'

definePageMeta({
  layout: 'admin',
})

const config = useRuntimeConfig()
const { fetchClasses: fetchAllClasses } = useClasses()

const classes = ref<ClassItem[]>([])
const loading = ref(false)
const showForm = ref(false)
const editingId = ref<number | null>(null)
const formError = ref('')
const formSuccess = ref('')

const form = reactive({
  name: '',
  description: '',
  max_students: 30,
  min_age: 5,
  max_age: 18,
})

const togglingIds = ref<Set<number>>(new Set())
const deletingIds = ref<Set<number>>(new Set())

async function loadClasses() {
  loading.value = true
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/classes`, { credentials: 'include' })
    if (res.code === 0) {
      classes.value = res.data || []
    }
  } catch (e) {
    console.error('Failed to load classes', e)
  } finally {
    loading.value = false
  }
}

function openAdd() {
  showForm.value = true
  editingId.value = null
  form.name = ''
  form.description = ''
  form.max_students = 30
  form.min_age = 5
  form.max_age = 18
  formError.value = ''
  formSuccess.value = ''
}

function openEdit(cls: ClassItem) {
  showForm.value = true
  editingId.value = cls.id
  form.name = cls.name
  form.description = cls.description
  form.max_students = cls.max_students
  form.min_age = cls.min_age
  form.max_age = cls.max_age
  formError.value = ''
  formSuccess.value = ''
}

function closeForm() {
  showForm.value = false
  editingId.value = null
}

async function handleSubmit() {
  formError.value = ''
  formSuccess.value = ''

  if (!form.name.trim()) {
    formError.value = '班级名称不能为空'
    return
  }
  if (form.min_age >= form.max_age) {
    formError.value = '年龄下限必须小于年龄上限'
    return
  }

  try {
    let res
    if (editingId.value !== null) {
      res = await $fetch<any>(`${config.public.apiBase}/classes/${editingId.value}`, {
        method: 'PUT',
        credentials: 'include',
        body: {
          name: form.name.trim(),
          description: form.description.trim(),
          max_students: form.max_students,
          min_age: form.min_age,
          max_age: form.max_age,
        },
      })
    } else {
      res = await $fetch<any>(`${config.public.apiBase}/classes`, {
        method: 'POST',
        credentials: 'include',
        body: {
          name: form.name.trim(),
          description: form.description.trim(),
          max_students: form.max_students,
          min_age: form.min_age,
          max_age: form.max_age,
        },
      })
    }
    if (res.code === 0) {
      formSuccess.value = editingId.value !== null ? '班级信息已更新' : '班级创建成功'
      setTimeout(closeForm, 800)
      loadClasses()
      fetchAllClasses()
    } else {
      formError.value = res.message || '操作失败'
    }
  } catch (e: any) {
    formError.value = e.data?.message || '操作失败，请检查网络'
  }
}

async function handleToggle(cls: ClassItem) {
  togglingIds.value.add(cls.id)
  togglingIds.value = new Set(togglingIds.value)
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/classes/${cls.id}/toggle`, {
      method: 'PUT',
      credentials: 'include',
      body: { enabled: !cls.enabled },
    })
    if (res.code === 0) {
      loadClasses()
      fetchAllClasses()
    }
  } catch (e) {
    console.error('Toggle failed', e)
  } finally {
    togglingIds.value = new Set([...togglingIds.value].filter(x => x !== cls.id))
  }
}

async function handleDelete(cls: ClassItem) {
  deletingIds.value.add(cls.id)
  deletingIds.value = new Set(deletingIds.value)
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/classes/${cls.id}`, {
      method: 'DELETE',
      credentials: 'include',
    })
    if (res.code === 0) {
      loadClasses()
      fetchAllClasses()
    }
  } catch (e) {
    console.error('Delete failed', e)
  } finally {
    deletingIds.value = new Set([...deletingIds.value].filter(x => x !== cls.id))
  }
}

onMounted(() => {
  loadClasses()
})
</script>

<template>
  <div class="classes-page">
    <div class="page-header">
      <div>
        <h1 class="page-title">班级管理</h1>
        <p class="page-desc">共 {{ classes.length }} 个班级</p>
      </div>
      <button class="btn-pill" @click="openAdd">
        <UIcon name="i-heroicons-plus" class="w-4 h-4" />
        添加班级
      </button>
    </div>

    <Transition name="form-panel">
      <div v-if="showForm" class="form-panel">
        <div class="form-panel-inner">
          <h2 class="form-panel-title">{{ editingId !== null ? '编辑班级' : '添加班级' }}</h2>
          <button class="form-panel-close" @click="closeForm">
            <UIcon name="i-heroicons-x-mark" class="w-4 h-4" />
          </button>
        </div>

        <div class="form-grid">
          <div class="form-field">
            <label class="form-label"><span class="form-required">*</span> 班级名称</label>
            <input v-model="form.name" type="text" class="input-base" placeholder="例如：一年级1班" maxlength="50">
          </div>
          <div class="form-field">
            <label class="form-label">班级描述</label>
            <input v-model="form.description" type="text" class="input-base" placeholder="简要说明（可选）" maxlength="255">
          </div>
          <div class="form-field">
            <label class="form-label"><span class="form-required">*</span> 最大报名人数</label>
            <input v-model.number="form.max_students" type="number" class="input-base" min="1" placeholder="30">
          </div>
          <div class="form-field">
            <label class="form-label"><span class="form-required">*</span> 年龄下限</label>
            <input v-model.number="form.min_age" type="number" class="input-base" min="0" placeholder="5">
          </div>
          <div class="form-field">
            <label class="form-label"><span class="form-required">*</span> 年龄上限</label>
            <input v-model.number="form.max_age" type="number" class="input-base" min="0" placeholder="18">
          </div>
        </div>

        <div v-if="formError" class="form-msg form-msg--error">
          <UIcon name="i-heroicons-exclamation-circle" class="msg-icon" />
          {{ formError }}
        </div>
        <div v-if="formSuccess" class="form-msg form-msg--success">
          <UIcon name="i-heroicons-check-circle" class="msg-icon" />
          {{ formSuccess }}
        </div>

        <div class="form-actions">
          <button class="btn-pill-neutral" @click="closeForm">取消</button>
          <button class="btn-pill" @click="handleSubmit">
            <UIcon v-if="!formError && !formSuccess" name="i-heroicons-check" class="w-4 h-4" />
            {{ editingId !== null ? '保存修改' : '创建班级' }}
          </button>
        </div>
      </div>
    </Transition>

    <div v-if="!loading && classes.length === 0" class="empty-state">
      <UIcon name="i-heroicons-academic-cap" class="empty-icon" />
      <p class="empty-title">暂无班级</p>
      <p class="empty-desc">点击右上角「添加班级」开始管理</p>
      <button class="btn-pill empty-btn" @click="openAdd">
        <UIcon name="i-heroicons-plus" class="w-4 h-4" />
        添加第一个班级
      </button>
    </div>

    <div v-else class="class-grid">
      <div
        v-for="cls in classes"
        :key="cls.id"
        class="class-card"
        :class="{ 'class-card--disabled': !cls.enabled }"
      >
        <div class="class-card-header">
          <div class="class-card-name">
            {{ cls.name }}
            <span class="status-tag" :class="cls.enabled ? 'status-tag--on' : 'status-tag--off'">
              {{ cls.enabled ? '启用' : '停用' }}
            </span>
          </div>
          <button class="class-card-edit" @click="openEdit(cls)">
            <UIcon name="i-heroicons-pencil-square" class="w-4 h-4" />
          </button>
        </div>

        <p v-if="cls.description" class="class-card-desc">{{ cls.description }}</p>

        <div class="class-card-stats">
          <div class="stat-item">
            <span class="stat-label">报名人数</span>
            <span class="stat-value">
              <span :class="cls.current_count >= cls.max_students ? 'stat-value--full' : ''">
                {{ cls.current_count }}
              </span>
              / {{ cls.max_students }}
            </span>
          </div>
          <div class="stat-item">
            <span class="stat-label">年龄范围</span>
            <span class="stat-value">{{ cls.min_age }} ~ {{ cls.max_age }} 岁</span>
          </div>
        </div>

        <div class="class-card-bar">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: `${Math.min(100, (cls.current_count / cls.max_students) * 100)}%` }" />
          </div>
        </div>

        <div class="class-card-actions">
          <button
            class="card-action-btn"
            :class="cls.enabled ? 'card-action-btn--off' : 'card-action-btn--on'"
            :disabled="togglingIds.has(cls.id)"
            @click="handleToggle(cls)"
          >
            <UIcon v-if="togglingIds.has(cls.id)" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
            <UIcon v-else :name="cls.enabled ? 'i-heroicons-pause-circle' : 'i-heroicons-play-circle'" class="w-3.5 h-3.5" />
            {{ cls.enabled ? '停用' : '启用' }}
          </button>
          <button
            class="card-action-btn card-action-btn--delete"
            :disabled="deletingIds.has(cls.id)"
            @click="handleDelete(cls)"
          >
            <UIcon v-if="deletingIds.has(cls.id)" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
            <UIcon v-else name="i-heroicons-trash" class="w-3.5 h-3.5" />
            删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.classes-page {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 56px - 2.5rem);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1e293b;
}

.page-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
  margin-top: 0.125rem;
}

.form-panel {
  background: #ffffff;
  border: 1px solid #f1f5f9;
  border-radius: 0.75rem;
  padding: 1.25rem;
  margin-bottom: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.form-panel-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.form-panel-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
}

.form-panel-close {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.form-panel-close:hover {
  background: #e2e8f0;
  color: #334155;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.75rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.form-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #334155;
}

.form-required {
  color: #ef4444;
}

.form-msg {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8125rem;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4375rem;
  margin-top: 0.5rem;
}

.form-msg--error {
  color: #b91c1c;
  background: #fef2f2;
}

.form-msg--success {
  color: #15803d;
  background: #f0fdf4;
}

.msg-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.form-panel-enter-active,
.form-panel-leave-active {
  transition: all 0.25s ease;
}

.form-panel-enter-from,
.form-panel-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.empty-state {
  background: #ffffff;
  border: 1px solid #f1f5f9;
  border-radius: 0.75rem;
  padding: 3rem 2rem;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
}

.empty-icon {
  width: 48px;
  height: 48px;
  color: #cbd5e1;
  margin-bottom: 1rem;
}

.empty-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #334155;
  margin-bottom: 0.375rem;
}

.empty-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
  margin-bottom: 1.25rem;
}

.empty-btn {
  margin: 0 auto;
}

.class-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 0.75rem;
}

.class-card {
  background: #ffffff;
  border: 1px solid #f1f5f9;
  border-radius: 0.75rem;
  padding: 1.25rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  transition: all 0.2s ease;
  border-left: 3px solid #2563eb;
}

.class-card:hover {
  box-shadow: 0 4px 16px -4px rgba(37, 99, 235, 0.08);
}

.class-card--disabled {
  opacity: 0.65;
  border-left-color: #cbd5e1;
}

.class-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.class-card-name {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-tag {
  display: inline-flex;
  padding: 0.0625rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.625rem;
  font-weight: 600;
}

.status-tag--on {
  background: #f0fdf4;
  color: #16a34a;
}

.status-tag--off {
  background: #f1f5f9;
  color: #94a3b8;
}

.class-card-edit {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.class-card-edit:hover {
  background: #eff6ff;
  color: #2563eb;
}

.class-card-desc {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-bottom: 0.75rem;
  line-height: 1.4;
}

.class-card-stats {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 0.75rem;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.stat-label {
  font-size: 0.625rem;
  color: #94a3b8;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.stat-value {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
}

.stat-value--full {
  color: #dc2626;
}

.class-card-bar {
  margin-bottom: 0.75rem;
}

.progress-bar {
  height: 4px;
  background: #f1f5f9;
  border-radius: 9999px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563eb, #38bdf8);
  border-radius: 9999px;
  transition: width 0.4s ease;
}

.class-card-actions {
  display: flex;
  gap: 0.5rem;
  border-top: 1px solid #f1f5f9;
  padding-top: 0.625rem;
}

.card-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.625rem;
  border: 1px solid #e2e8f0;
  border-radius: 9999px;
  background: #ffffff;
  color: #64748b;
  font-size: 0.6875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.card-action-btn:hover:not(:disabled) {
  border-color: #93c5fd;
  color: #2563eb;
}

.card-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.card-action-btn--on {
  border-color: #bbf7d0;
  color: #16a34a;
}

.card-action-btn--on:hover {
  background: #f0fdf4;
  border-color: #86efac;
}

.card-action-btn--delete {
  border-color: #fecaca;
  color: #94a3b8;
}

.card-action-btn--delete:hover {
  background: #fef2f2;
  border-color: #fca5a5;
  color: #dc2626;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .class-grid {
    grid-template-columns: 1fr;
  }
}
</style>