<script setup lang="ts">
import type { ClassItem } from '~/types'

definePageMeta({
  layout: 'admin',
})

const { registrations, loading, fetchRegistrations, deleteAndRefresh, batchDelete, getPhotoUrl } = useRegistrations()
const { classes, fetchClasses } = useClasses()
const config = useRuntimeConfig()

const selectedIds = ref<Set<number>>(new Set())
const selectMode = ref(false)
const showFilter = ref(false)

const filterClassId = ref<number | null>(null)
const filterName = ref('')

const photoModalUrl = ref<string | null>(null)
const deletingIds = ref<Set<number>>(new Set())

function toggleSelectMode() {
  selectMode.value = !selectMode.value
  if (!selectMode.value) {
    selectedIds.value = new Set()
  }
}

function toggleSelectAll() {
  if (!registrations.value.length) return
  if (selectedIds.value.size === registrations.value.length) {
    selectedIds.value = new Set()
  } else {
    selectedIds.value = new Set(registrations.value.map(r => r.id))
  }
}

function toggleSelect(id: number) {
  const next = new Set(selectedIds.value)
  if (next.has(id)) {
    next.delete(id)
  } else {
    next.add(id)
  }
  selectedIds.value = next
}

function isAllSelected() {
  return registrations.value.length > 0 && selectedIds.value.size === registrations.value.length
}

async function handleDelete(id: number) {
  deletingIds.value.add(id)
  deletingIds.value = new Set(deletingIds.value)
  await deleteAndRefresh(id, () => fetchClasses())
  deletingIds.value = new Set([...deletingIds.value].filter(x => x !== id))
}

async function handleBatchDelete() {
  if (selectedIds.value.size === 0) return
  const ids = [...selectedIds.value]
  await batchDelete(ids, () => fetchClasses())
  selectedIds.value = new Set()
  selectMode.value = false
}

function handleFilter() {
  showFilter.value = !showFilter.value
}

function applyFilter() {
  fetchRegistrations(filterName.value || undefined, filterClassId.value || undefined)
}

function resetFilter() {
  filterName.value = ''
  filterClassId.value = null
  fetchRegistrations()
}

function viewPhoto(photoPath: string | null) {
  const url = getPhotoUrl(photoPath)
  if (url) {
    photoModalUrl.value = url
  }
}

function closePhotoModal() {
  photoModalUrl.value = null
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  return dateStr.slice(0, 10)
}

function formatIdNumber(idNumber: string): string {
  if (!idNumber) return '-'
  return idNumber.slice(0, 6) + '****' + idNumber.slice(14)
}

function exportCSV() {
  if (!registrations.value.length) return
  const headers = ['ID', '姓名', '性别', '出生日期', '身份证号', '年级', '班级', '家长姓名', '联系电话', '家庭住址', '报名时间']
  const rows = registrations.value.map(r => [
    r.id, r.name, r.gender, formatDate(r.birth_date), r.id_number,
    r.grade, r.class?.name || '', r.parent_name, r.parent_phone, r.address,
    formatDate(r.registration_time),
  ])
  const csv = [headers.join(','), ...rows.map(row => row.map(v => `"${String(v).replace(/"/g, '""')}"`).join(','))].join('\r\n')
  downloadBlob('\uFEFF' + csv, `registrations_${new Date().toISOString().slice(0, 10)}.csv`, 'text/csv;charset=utf-8')
}

function exportExcel() {
  if (!registrations.value.length) return
  const headers = ['ID', '姓名', '性别', '出生日期', '身份证号', '年级', '班级', '家长姓名', '联系电话', '家庭住址', '报名时间']
  const headRow = `<tr>${headers.map(h => `<th>${h}</th>`).join('')}</tr>`
  const dataRows = registrations.value.map(r => `<tr>
    <td>${r.id}</td><td>${r.name}</td><td>${r.gender}</td><td>${formatDate(r.birth_date)}</td>
    <td>${r.id_number}</td><td>${r.grade}</td><td>${r.class?.name || ''}</td>
    <td>${r.parent_name}</td><td>${r.parent_phone}</td><td>${r.address}</td>
    <td>${formatDate(r.registration_time)}</td>
  </tr>`).join('')

  const html = `<html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel">
<head><meta charset="UTF-8"></head><body><table border="1">${headRow}${dataRows}</table></body></html>`
  downloadBlob(html, `registrations_${new Date().toISOString().slice(0, 10)}.xls`, 'application/vnd.ms-excel')
}

function downloadBlob(content: string, filename: string, mime: string) {
  const blob = new Blob([content], { type: mime })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(() => {
  fetchRegistrations()
  fetchClasses()
})
</script>

<template>
  <div class="registrations-page">
    <div class="page-header">
      <h1 class="page-title">报名列表</h1>
      <p class="page-desc">共 {{ registrations.length }} 条报名记录</p>
    </div>

    <div class="action-bar">
      <button class="btn-pill-neutral action-btn" :class="{ 'action-btn--active': showFilter }" @click="handleFilter">
        <UIcon name="i-heroicons-funnel" class="w-4 h-4" />
        筛选
      </button>
      <button class="btn-pill-neutral action-btn" :class="{ 'action-btn--active': selectMode }" @click="toggleSelectMode">
        <UIcon name="i-heroicons-pencil-square" class="w-4 h-4" />
        {{ selectMode ? '取消选择' : '编辑' }}
      </button>
      <button class="btn-pill-neutral action-btn" :disabled="selectedIds.size === 0" @click="handleBatchDelete">
        <UIcon name="i-heroicons-trash" class="w-4 h-4" />
        批量删除
        <span v-if="selectedIds.size > 0" class="badge-count">{{ selectedIds.size }}</span>
      </button>
      <button class="btn-pill-neutral action-btn" @click="exportCSV">
        <UIcon name="i-heroicons-document-text" class="w-4 h-4" />
        导出 CSV
      </button>
      <button class="btn-pill-neutral action-btn" @click="exportExcel">
        <UIcon name="i-heroicons-document-chart-bar" class="w-4 h-4" />
        导出 Excel
      </button>
    </div>

    <div v-if="showFilter" class="filter-panel">
      <div class="filter-row">
        <div class="filter-field">
          <label class="filter-label">班级筛选</label>
          <select v-model.number="filterClassId" class="input-base" style="min-width: 160px;">
            <option :value="null">全部班级</option>
            <option v-for="cls in classes" :key="cls.id" :value="cls.id">{{ cls.name }}</option>
          </select>
        </div>
        <div class="filter-field">
          <label class="filter-label">姓名搜索</label>
          <input v-model="filterName" type="text" class="input-base" placeholder="输入姓名搜索" style="min-width: 160px;" @keyup.enter="applyFilter">
        </div>
        <div class="filter-actions">
          <button class="btn-pill filter-apply-btn" @click="applyFilter">搜索</button>
          <button class="btn-pill-neutral" @click="resetFilter">重置</button>
        </div>
      </div>
    </div>

    <div class="table-wrap">
      <div class="table-scroll">
        <table class="data-table">
          <thead>
            <tr>
              <th v-if="selectMode" class="th-check">
                <input type="checkbox" :checked="isAllSelected()" class="table-check" @change="toggleSelectAll">
              </th>
              <th class="th-id">ID</th>
              <th>姓名</th>
              <th>性别</th>
              <th class="th-date">出生日期</th>
              <th class="th-idnum">身份证号</th>
              <th>年级</th>
              <th>班级</th>
              <th>家长姓名</th>
              <th>联系电话</th>
              <th class="th-address">家庭住址</th>
              <th>户口本图片</th>
              <th class="th-actions">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="loading && !registrations.length">
              <td :colspan="selectMode ? 13 : 12" class="td-empty">
                <UIcon name="i-heroicons-arrow-path" class="w-5 h-5 animate-spin" />
                <span>加载中...</span>
              </td>
            </tr>
            <tr v-else-if="!loading && !registrations.length">
              <td :colspan="selectMode ? 13 : 12" class="td-empty">
                <UIcon name="i-heroicons-inbox" class="w-10 h-10 text-gray-300" />
                <span>暂无报名记录</span>
              </td>
            </tr>
            <tr
              v-for="record in registrations"
              :key="record.id"
              :class="{ 'row-selected': selectedIds.has(record.id) }"
            >
              <td v-if="selectMode" class="td-check" @click="toggleSelect(record.id)">
                <input type="checkbox" :checked="selectedIds.has(record.id)" class="table-check">
              </td>
              <td class="td-mono">{{ record.id }}</td>
              <td class="td-name">{{ record.name }}</td>
              <td>
                <span class="gender-tag" :class="record.gender === '男' ? 'gender-tag--male' : 'gender-tag--female'">
                  {{ record.gender }}
                </span>
              </td>
              <td class="td-mono">{{ formatDate(record.birth_date) }}</td>
              <td class="td-mono td-idnum">{{ formatIdNumber(record.id_number) }}</td>
              <td>{{ record.grade }}</td>
              <td>
                <span class="class-tag">{{ record.class?.name || '-' }}</span>
              </td>
              <td>{{ record.parent_name }}</td>
              <td class="td-mono">{{ record.parent_phone }}</td>
              <td class="td-address" :title="record.address">{{ record.address }}</td>
              <td>
                <button
                  class="photo-btn"
                  :disabled="!record.photo_path"
                  @click="viewPhoto(record.photo_path)"
                >
                  <UIcon name="i-heroicons-eye" class="w-3.5 h-3.5" />
                  查看照片
                </button>
              </td>
              <td>
                <button
                  class="td-delete-btn"
                  :disabled="deletingIds.has(record.id)"
                  @click="handleDelete(record.id)"
                >
                  <UIcon
                    v-if="deletingIds.has(record.id)"
                    name="i-heroicons-arrow-path"
                    class="w-3.5 h-3.5 animate-spin"
                  />
                  <UIcon v-else name="i-heroicons-trash" class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="photoModalUrl" class="photo-modal-mask" @click.self="closePhotoModal">
        <div class="photo-modal-body">
          <button class="photo-modal-close" @click="closePhotoModal">
            <UIcon name="i-heroicons-x-mark" class="w-5 h-5" />
          </button>
          <img :src="photoModalUrl" alt="户口本照片" class="photo-modal-img">
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.registrations-page {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 56px - 2.5rem);
}

.page-header {
  display: flex;
  align-items: baseline;
  gap: 0.75rem;
  margin-bottom: 0.875rem;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1e293b;
}

.page-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
}

.action-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.875rem;
  flex-wrap: wrap;
}

.action-btn {
  font-size: 0.75rem;
  padding: 0.375rem 0.8125rem;
  position: relative;
}

.action-btn--active {
  background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 100%);
  color: #ffffff;
  border-color: transparent;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.action-btn--active:hover {
  background: linear-gradient(135deg, #1d4ed8 0%, #1e40af 100%);
  color: #ffffff;
}

.badge-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 18px;
  height: 18px;
  border-radius: 9999px;
  background: #ef4444;
  color: #ffffff;
  font-size: 0.625rem;
  font-weight: 700;
  margin-left: 0.25rem;
}

.filter-panel {
  background: #ffffff;
  border: 1px solid #f1f5f9;
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  margin-bottom: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  animation: slideDown 0.2s ease;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

.filter-row {
  display: flex;
  align-items: flex-end;
  gap: 1rem;
  flex-wrap: wrap;
}

.filter-field {
  display: flex;
  flex-direction: column;
  gap: 0.3125rem;
}

.filter-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: #64748b;
}

.filter-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-apply-btn {
  padding: 0.375rem 0.875rem;
  font-size: 0.75rem;
}

.table-wrap {
  flex: 1;
  background: #ffffff;
  border-radius: 0.75rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  overflow: hidden;
}

.table-scroll {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8125rem;
}

.data-table thead {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

.data-table th {
  padding: 0.75rem 0.75rem;
  text-align: left;
  font-weight: 600;
  color: #475569;
  font-size: 0.75rem;
  white-space: nowrap;
  user-select: none;
}

.th-check {
  width: 36px;
  text-align: center;
}

.th-id {
  width: 48px;
}

.th-date {
  width: 100px;
}

.th-idnum {
  width: 120px;
}

.th-actions {
  width: 56px;
  text-align: center;
}

.th-address {
  min-width: 140px;
}

.data-table td {
  padding: 0.625rem 0.75rem;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  vertical-align: middle;
}

.data-table tbody tr {
  transition: background 0.1s ease;
}

.data-table tbody tr:hover {
  background: #f8fafc;
}

.row-selected {
  background: #eff6ff !important;
}

.td-check {
  text-align: center;
  cursor: pointer;
}

.table-check {
  width: 15px;
  height: 15px;
  cursor: pointer;
  accent-color: #2563eb;
}

.td-mono {
  font-variant-numeric: tabular-nums;
  font-family: "SF Mono", "Cascadia Code", "Consolas", monospace;
  font-size: 0.78125rem;
}

.td-name {
  font-weight: 600;
  color: #1e293b;
}

.gender-tag {
  display: inline-flex;
  align-items: center;
  padding: 0.125rem 0.5625rem;
  border-radius: 9999px;
  font-size: 0.6875rem;
  font-weight: 500;
}

.gender-tag--male {
  background: #eff6ff;
  color: #2563eb;
}

.gender-tag--female {
  background: #fdf2f8;
  color: #db2777;
}

.class-tag {
  display: inline-flex;
  padding: 0.125rem 0.5625rem;
  border-radius: 9999px;
  font-size: 0.6875rem;
  font-weight: 500;
  background: #f0fdf4;
  color: #16a34a;
}

.td-address {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.photo-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.625rem;
  border: 1px solid #e2e8f0;
  border-radius: 9999px;
  background: #ffffff;
  color: #2563eb;
  font-size: 0.6875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  white-space: nowrap;
}

.photo-btn:hover:not(:disabled) {
  background: #eff6ff;
  border-color: #93c5fd;
}

.photo-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  color: #94a3b8;
}

.td-delete-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 9999px;
  background: transparent;
  color: #94a3b8;
  cursor: pointer;
  transition: all 0.15s ease;
}

.td-delete-btn:hover:not(:disabled) {
  background: #fef2f2;
  color: #ef4444;
}

.td-delete-btn:disabled {
  cursor: not-allowed;
}

.td-empty {
  text-align: center;
  padding: 3rem 1rem !important;
  color: #94a3b8;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.photo-modal-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.7);
  backdrop-filter: blur(4px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.photo-modal-body {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  background: #ffffff;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 24px 48px -12px rgba(0, 0, 0, 0.3);
}

.photo-modal-close {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: rgba(0, 0, 0, 0.4);
  color: #ffffff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  transition: background 0.15s ease;
}

.photo-modal-close:hover {
  background: rgba(0, 0, 0, 0.6);
}

.photo-modal-img {
  display: block;
  max-width: 80vw;
  max-height: 80vh;
  object-fit: contain;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@media (max-width: 1024px) {
  .th-address {
    min-width: 120px;
  }
}

@media (max-width: 768px) {
  .action-bar {
    justify-content: flex-start;
  }
}
</style>