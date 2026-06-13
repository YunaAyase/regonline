<script setup lang="ts">
import type { SiteSettings, ServerInfo, UpdateAccountRequest } from '~/types'

definePageMeta({
  layout: 'admin',
})

const config = useRuntimeConfig()

const tabs = [
  { key: 'account', label: '管理员账号', icon: 'i-heroicons-user-circle' },
  { key: 'site', label: '网站信息', icon: 'i-heroicons-globe-alt' },
  { key: 'ocr', label: 'OCR 云识别', icon: 'i-heroicons-camera' },
  { key: 'advanced', label: '高级设置', icon: 'i-heroicons-wrench-screwdriver' },
  { key: 'server', label: '服务器信息', icon: 'i-heroicons-server-stack' },
]

const activeTab = ref('account')

const accountForm = ref<UpdateAccountRequest>({
  username: '',
  old_password: '',
  new_password: '',
})
const confirmPassword = ref('')
const accountSaving = ref(false)
const accountError = ref('')
const accountSuccess = ref('')

const siteForm = ref<SiteSettings>({
  site_name: '',
  site_description: '',
  icp_record: '',
  copyright: '',
  ocr_provider: 'baidu',
  ocr_baidu_api_key: '',
  ocr_baidu_secret_key: '',
  ocr_alibaba_access_key_id: '',
  ocr_alibaba_access_key_secret: '',
})
const siteSaving = ref(false)
const siteError = ref('')
const siteSuccess = ref('')

const ocrSaving = ref(false)
const ocrError = ref('')
const ocrSuccess = ref('')

const serverInfo = ref<ServerInfo | null>(null)
const serverLoading = ref(false)

const dataExportLoading = ref(false)

async function loadSettings() {
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/settings`, {
      credentials: 'include',
    })
    if (res.code === 0 && res.data) {
      siteForm.value = { ...res.data }
    }
  } catch (e) {
    console.error('Failed to load settings', e)
  }
}

async function loadServerInfo() {
  serverLoading.value = true
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/server-info`, {
      credentials: 'include',
    })
    if (res.code === 0) {
      serverInfo.value = res.data
    }
  } catch (e) {
    console.error('Failed to load server info', e)
  } finally {
    serverLoading.value = false
  }
}

async function saveAccount() {
  accountError.value = ''
  accountSuccess.value = ''

  if (!accountForm.value.old_password) {
    accountError.value = '请输入原密码'
    return
  }
  if (accountForm.value.new_password && accountForm.value.new_password !== confirmPassword.value) {
    accountError.value = '两次输入的新密码不一致'
    return
  }

  accountSaving.value = true
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/admin/account`, {
      method: 'PUT',
      credentials: 'include',
      body: accountForm.value,
    })
    if (res.code === 0) {
      accountSuccess.value = '账号信息修改成功'
      accountForm.value.old_password = ''
      accountForm.value.new_password = ''
      confirmPassword.value = ''
    } else {
      accountError.value = res.message || '修改失败'
    }
  } catch (e: any) {
    accountError.value = e.data?.message || '修改失败，请检查网络'
  } finally {
    accountSaving.value = false
  }
}

async function saveSiteSettings() {
  siteError.value = ''
  siteSuccess.value = ''

  siteSaving.value = true
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/settings`, {
      method: 'PUT',
      credentials: 'include',
      body: siteForm.value,
    })
    if (res.code === 0) {
      siteSuccess.value = '网站信息保存成功'
    } else {
      siteError.value = res.message || '保存失败'
    }
  } catch (e: any) {
    siteError.value = e.data?.message || '保存失败，请检查网络'
  } finally {
    siteSaving.value = false
  }
}

async function saveOCR() {
  ocrError.value = ''
  ocrSuccess.value = ''

  ocrSaving.value = true
  try {
    const body: Record<string, string> = {
      ocr_provider: siteForm.value.ocr_provider,
    }
    if (siteForm.value.ocr_provider === 'alibaba') {
      body.ocr_alibaba_access_key_id = siteForm.value.ocr_alibaba_access_key_id
      body.ocr_alibaba_access_key_secret = siteForm.value.ocr_alibaba_access_key_secret
    } else {
      body.ocr_baidu_api_key = siteForm.value.ocr_baidu_api_key
      body.ocr_baidu_secret_key = siteForm.value.ocr_baidu_secret_key
    }

    const res = await $fetch<any>(`${config.public.apiBase}/settings`, {
      method: 'PUT',
      credentials: 'include',
      body,
    })
    if (res.code === 0) {
      ocrSuccess.value = 'OCR 配置保存成功'
    } else {
      ocrError.value = res.message || '保存失败'
    }
  } catch (e: any) {
    ocrError.value = e.data?.message || '保存失败，请检查网络'
  } finally {
    ocrSaving.value = false
  }
}

const backupLoading = ref(false)
const resetLoading = ref(false)
const qrcodeLoading = ref(false)
const advancedMsg = ref('')
const advancedMsgType = ref<'success' | 'error'>('success')

const qrcodeUrl = ref('')
const qrcodePreview = ref('')
const qrcodeDownloading = ref(false)
const serverIPLoading = ref(false)
const serverIPType = ref('')

async function fetchServerIP() {
  serverIPLoading.value = true
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/server-ip`, {
      credentials: 'include',
    })
    if (res.code === 0 && res.data?.url) {
      qrcodeUrl.value = res.data.url
      const typeMap: Record<string, string> = {
        public_ipv4: '公网 IPv4',
        public_ipv6: '公网 IPv6',
        private_ipv4: '内网 IPv4',
      }
      serverIPType.value = typeMap[res.data.ip_type] || ''
      return true
    }
    return false
  } catch {
    console.warn('Failed to fetch server IP')
    return false
  } finally {
    serverIPLoading.value = false
  }
}

async function handleBackup() {
  backupLoading.value = true
  advancedMsg.value = ''
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/backup`, {
      method: 'POST',
      credentials: 'include',
    })
    if (res.code === 0) {
      advancedMsg.value = `备份成功：${res.data.file}`
      advancedMsgType.value = 'success'
    } else {
      advancedMsg.value = res.message || '备份失败'
      advancedMsgType.value = 'error'
    }
  } catch (e: any) {
    advancedMsg.value = e.data?.message || '备份失败，请检查网络'
    advancedMsgType.value = 'error'
  } finally {
    backupLoading.value = false
  }
}

async function handleReset() {
  if (!confirm('⚠️ 确定要重置数据库吗？所有数据将被清空！')) return
  if (!confirm('再次确认：此操作不可逆，将删除所有报名、班级和设置数据。')) return

  resetLoading.value = true
  advancedMsg.value = ''
  try {
    const res = await $fetch<any>(`${config.public.apiBase}/reset-db`, {
      method: 'POST',
      credentials: 'include',
    })
    if (res.code === 0) {
      advancedMsg.value = `数据库已重置，旧数据已备份为 ${res.data.backup}`
      advancedMsgType.value = 'success'
    } else {
      advancedMsg.value = res.message || '重置失败'
      advancedMsgType.value = 'error'
    }
  } catch (e: any) {
    advancedMsg.value = e.data?.message || '重置失败，请检查网络'
    advancedMsgType.value = 'error'
  } finally {
    resetLoading.value = false
  }
}

async function generateQRCode() {
  if (!qrcodeUrl.value.trim()) {
    advancedMsg.value = ''
    const ok = await fetchServerIP()
    if (!ok || !qrcodeUrl.value.trim()) {
      advancedMsg.value = '无法获取服务器地址，请手动输入 URL'
      advancedMsgType.value = 'error'
      return
    }
  }

  qrcodeLoading.value = true
  advancedMsg.value = ''
  try {
    const response = await fetch(`${config.public.apiBase}/qrcode/generate`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ url: qrcodeUrl.value }),
    })

    if (!response.ok) {
      throw new Error('生成失败')
    }

    const blob = await response.blob()
    qrcodePreview.value = URL.createObjectURL(blob)
    advancedMsg.value = '二维码生成成功'
    advancedMsgType.value = 'success'
  } catch (e: any) {
    advancedMsg.value = e.message || '生成失败，请检查网络'
    advancedMsgType.value = 'error'
  } finally {
    qrcodeLoading.value = false
  }
}

async function downloadQRCode() {
  if (!qrcodePreview.value) {
    advancedMsg.value = '请先生成二维码'
    advancedMsgType.value = 'error'
    return
  }

  qrcodeDownloading.value = true
  try {
    const link = document.createElement('a')
    link.href = qrcodePreview.value
    link.download = 'qrcode.png'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    advancedMsg.value = '二维码已下载'
    advancedMsgType.value = 'success'
  } catch (e: any) {
    advancedMsg.value = '下载失败'
    advancedMsgType.value = 'error'
  } finally {
    qrcodeDownloading.value = false
  }
}

function formatUptime(seconds: number): string {
  const d = Math.floor(seconds / 86400)
  const h = Math.floor((seconds % 86400) / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const parts: string[] = []
  if (d > 0) parts.push(`${d}天`)
  if (h > 0) parts.push(`${h}小时`)
  parts.push(`${m}分钟`)
  return parts.join(' ')
}

watch(activeTab, (tab) => {
  if (tab === 'server') {
    loadServerInfo()
  }
})

onMounted(() => {
  loadSettings()
})
</script>

<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">网站设置</h1>
      <p class="page-description">系统配置与参数管理</p>
    </div>

    <div class="settings-body">
      <div class="tabs-sidebar">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="tab-btn"
          :class="{ 'tab-btn--active': activeTab === tab.key }"
          @click="activeTab = tab.key"
        >
          <UIcon :name="tab.icon" class="tab-btn-icon" />
          <span class="tab-btn-text">{{ tab.label }}</span>
        </button>
      </div>

      <div class="tab-content">
        <div v-if="activeTab === 'account'" class="tab-panel">
          <div class="panel-header">
            <h2 class="panel-title">管理员账号</h2>
            <p class="panel-desc">修改管理员用户名与密码，密码使用 bcrypt 加密存储</p>
          </div>

          <div class="form-card">
            <div class="form-group">
              <label class="form-label">用户名</label>
              <input
                v-model="accountForm.username"
                type="text"
                class="input-base"
                placeholder="新用户名"
                style="max-width: 400px;"
              >
            </div>

            <div class="form-group">
              <label class="form-label"><span class="form-required">*</span> 原密码</label>
              <input
                v-model="accountForm.old_password"
                type="password"
                class="input-base"
                placeholder="输入原密码以验证身份"
                style="max-width: 400px;"
              >
            </div>

            <div class="form-group">
              <label class="form-label">新密码</label>
              <input
                v-model="accountForm.new_password"
                type="password"
                class="input-base"
                placeholder="留空则保持原密码不变"
                style="max-width: 400px;"
              >
            </div>

            <div class="form-group">
              <label class="form-label">确认新密码</label>
              <input
                v-model="confirmPassword"
                type="password"
                class="input-base"
                placeholder="再次输入新密码"
                style="max-width: 400px;"
              >
            </div>

            <div v-if="accountError" class="form-message form-message--error">
              <UIcon name="i-heroicons-exclamation-circle" class="msg-icon" />
              {{ accountError }}
            </div>
            <div v-if="accountSuccess" class="form-message form-message--success">
              <UIcon name="i-heroicons-check-circle" class="msg-icon" />
              {{ accountSuccess }}
            </div>

            <button class="btn-pill-neutral form-submit-btn" :disabled="accountSaving" @click="saveAccount">
              <UIcon v-if="accountSaving" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
              保存修改
            </button>
          </div>
        </div>

        <div v-if="activeTab === 'site'" class="tab-panel">
          <div class="panel-header">
            <h2 class="panel-title">网站信息</h2>
            <p class="panel-desc">修改网站的基础信息，所有数据存储在 SQLite 中</p>
          </div>

          <div class="form-card">
            <div class="form-group">
              <label class="form-label">网站名称</label>
              <input
                v-model="siteForm.site_name"
                type="text"
                class="input-base"
                placeholder="例如：2025年秋季招生报名系统"
                style="max-width: 480px;"
              >
            </div>

            <div class="form-group">
              <label class="form-label">网站描述</label>
              <textarea
                v-model="siteForm.site_description"
                class="input-base input-textarea"
                placeholder="简要介绍网站用途"
                rows="3"
                style="max-width: 480px;"
              />
            </div>

            <div class="form-group">
              <label class="form-label">ICP 备案号</label>
              <input
                v-model="siteForm.icp_record"
                type="text"
                class="input-base"
                placeholder="例如：京ICP备XXXXXXXX号"
                style="max-width: 480px;"
              >
            </div>

            <div class="form-group">
              <label class="form-label">版权声明</label>
              <input
                v-model="siteForm.copyright"
                type="text"
                class="input-base"
                placeholder="例如：&copy; 2025 RegOnline. All rights reserved."
                style="max-width: 480px;"
              >
            </div>

            <div v-if="siteError" class="form-message form-message--error">
              <UIcon name="i-heroicons-exclamation-circle" class="msg-icon" />
              {{ siteError }}
            </div>
            <div v-if="siteSuccess" class="form-message form-message--success">
              <UIcon name="i-heroicons-check-circle" class="msg-icon" />
              {{ siteSuccess }}
            </div>

            <button class="btn-pill-neutral form-submit-btn" :disabled="siteSaving" @click="saveSiteSettings">
              <UIcon v-if="siteSaving" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
              保存设置
            </button>
          </div>
        </div>

        <div v-if="activeTab === 'ocr'" class="tab-panel">
          <div class="panel-header">
            <h2 class="panel-title">OCR 云识别配置</h2>
            <p class="panel-desc">配置百度云/阿里云 OCR 识别服务，用于自动识别户口本照片中的身份证号</p>
          </div>

          <div class="form-card">
            <div class="form-group">
              <label class="form-label">OCR 服务商</label>
              <div class="ocr-provider-row">
                <label class="ocr-provider-option" :class="{ 'ocr-provider-option--active': siteForm.ocr_provider === 'baidu' }">
                  <input v-model="siteForm.ocr_provider" type="radio" value="baidu" class="sr-only" />
                  <span class="ocr-provider-label">百度云 OCR</span>
                  <span class="ocr-provider-desc">免费 1000次/月</span>
                </label>
                <label class="ocr-provider-option" :class="{ 'ocr-provider-option--active': siteForm.ocr_provider === 'alibaba' }">
                  <input v-model="siteForm.ocr_provider" type="radio" value="alibaba" class="sr-only" />
                  <span class="ocr-provider-label">阿里云 OCR</span>
                  <span class="ocr-provider-desc">免费 200次/月</span>
                </label>
              </div>
            </div>

            <!-- 百度云配置 -->
            <template v-if="siteForm.ocr_provider === 'baidu'">
              <div class="form-group">
                <label class="form-label">API Key</label>
                <input
                  v-model="siteForm.ocr_baidu_api_key"
                  type="text"
                  class="input-base"
                  placeholder="请输入百度云 API Key"
                  style="max-width: 480px;"
                >
                <p class="form-hint">
                  在
                  <a href="https://console.bce.baidu.com/ai/#/ai/ocr/overview/index" target="_blank" class="form-link">百度智能云控制台</a>
                  中创建应用获取
                </p>
              </div>
              <div class="form-group">
                <label class="form-label">Secret Key</label>
                <input
                  v-model="siteForm.ocr_baidu_secret_key"
                  type="password"
                  class="input-base"
                  placeholder="请输入百度云 Secret Key"
                  style="max-width: 480px;"
                >
              </div>
            </template>

            <!-- 阿里云配置 -->
            <template v-if="siteForm.ocr_provider === 'alibaba'">
              <div class="form-group">
                <label class="form-label">AccessKey ID</label>
                <input
                  v-model="siteForm.ocr_alibaba_access_key_id"
                  type="text"
                  class="input-base"
                  placeholder="请输入阿里云 AccessKey ID"
                  style="max-width: 480px;"
                >
                <p class="form-hint">
                  在
                  <a href="https://ram.console.aliyun.com/manage/ak" target="_blank" class="form-link">阿里云 RAM 访问控制</a>
                  中创建 AccessKey 获取
                </p>
              </div>
              <div class="form-group">
                <label class="form-label">AccessKey Secret</label>
                <input
                  v-model="siteForm.ocr_alibaba_access_key_secret"
                  type="password"
                  class="input-base"
                  placeholder="请输入阿里云 AccessKey Secret"
                  style="max-width: 480px;"
                >
              </div>
            </template>

            <div v-if="ocrError" class="form-message form-message--error">
              <UIcon name="i-heroicons-exclamation-circle" class="msg-icon" />
              {{ ocrError }}
            </div>
            <div v-if="ocrSuccess" class="form-message form-message--success">
              <UIcon name="i-heroicons-check-circle" class="msg-icon" />
              {{ ocrSuccess }}
            </div>

            <button class="btn-pill-neutral form-submit-btn" :disabled="ocrSaving" @click="saveOCR">
              <UIcon v-if="ocrSaving" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
              保存配置
            </button>
          </div>
        </div>

        <div v-if="activeTab === 'advanced'" class="tab-panel">
          <div class="panel-header">
            <h2 class="panel-title">高级设置</h2>
            <p class="panel-desc">数据管理与系统维护工具</p>
          </div>

          <div class="advanced-grid">
            <div class="adv-card">
              <div class="adv-card-left">
                <div class="adv-icon-wrap adv-icon-wrap--amber">
                  <UIcon name="i-heroicons-archive-box" class="adv-icon" />
                </div>
                <div>
                  <h3 class="adv-card-title">数据库备份</h3>
                  <p class="adv-card-desc">将当前 SQLite 数据库复制到 backups 目录</p>
                </div>
              </div>
              <button class="btn-pill-neutral" :disabled="backupLoading" @click="handleBackup">
                <UIcon v-if="backupLoading" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
                {{ backupLoading ? '备份中...' : '立即备份' }}
              </button>
            </div>

            <div class="adv-card">
              <div class="adv-card-left">
                <div class="adv-icon-wrap adv-icon-wrap--red">
                  <UIcon name="i-heroicons-arrow-path-rounded-square" class="adv-icon" />
                </div>
                <div>
                  <h3 class="adv-card-title">重置数据库</h3>
                  <p class="adv-card-desc">删除所有数据并重建数据库，自动备份后再重置</p>
                </div>
              </div>
              <button class="btn-pill-neutral adv-btn--danger" :disabled="resetLoading" @click="handleReset">
                <UIcon v-if="resetLoading" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
                {{ resetLoading ? '重置中...' : '重置数据库' }}
              </button>
            </div>

            <div class="adv-card adv-card--full">
              <div class="qrcode-section">
                <div class="qrcode-header">
                  <div class="adv-icon-wrap adv-icon-wrap--green">
                    <UIcon name="i-heroicons-qr-code" class="adv-icon" />
                  </div>
                  <div>
                    <h3 class="adv-card-title">报名二维码</h3>
                    <p class="adv-card-desc">生成报名页面的二维码，扫描后可直接进入报名页面</p>
                  </div>
                </div>

                <div class="qrcode-body">
                  <div class="qrcode-input-group">
                    <label class="form-label">报名页面 URL</label>
                    <div class="qrcode-url-row">
                      <input
                        v-model="qrcodeUrl"
                        type="text"
                        class="input-base"
                        placeholder="留空点击生成将自动检测服务器地址"
                        style="flex: 1;"
                        :disabled="serverIPLoading"
                      >
                      <span v-if="serverIPType" class="qrcode-ip-badge">{{ serverIPType }}</span>
                    </div>
                    <button class="btn-pill-neutral" :disabled="qrcodeLoading" @click="generateQRCode">
                      <UIcon v-if="qrcodeLoading" name="i-heroicons-arrow-path" class="w-3.5 h-3.5 animate-spin" />
                      {{ qrcodeLoading ? '生成中...' : '生成二维码' }}
                    </button>
                  </div>

                  <div v-if="qrcodePreview" class="qrcode-preview">
                    <img :src="qrcodePreview" alt="QR Code" class="qrcode-image">
                    <button class="btn-pill-neutral" :disabled="qrcodeDownloading" @click="downloadQRCode">
                      <UIcon v-if="qrcodeDownloading" name="i-heroicons-arrow-down-tray" class="w-3.5 h-3.5 animate-spin" />
                      {{ qrcodeDownloading ? '下载中...' : '下载二维码' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="advancedMsg" class="adv-msg" :class="advancedMsgType === 'error' ? 'adv-msg--error' : 'adv-msg--success'">
            <UIcon :name="advancedMsgType === 'error' ? 'i-heroicons-exclamation-circle' : 'i-heroicons-check-circle'" class="msg-icon" />
            {{ advancedMsg }}
          </div>
        </div>

        <div v-if="activeTab === 'server'" class="tab-panel">
          <div class="panel-header">
            <h2 class="panel-title">服务器信息</h2>
            <p class="panel-desc">实时系统运行状态</p>
          </div>

          <div v-if="serverLoading" class="loading-state">
            <UIcon name="i-heroicons-arrow-path" class="loading-icon animate-spin" />
            <span>加载中...</span>
          </div>

          <div v-else-if="serverInfo" class="server-grid">
            <div class="server-card">
              <div class="server-card-label">Go 版本</div>
              <div class="server-card-value">{{ serverInfo.go_version }}</div>
            </div>
            <div class="server-card">
              <div class="server-card-label">操作系统</div>
              <div class="server-card-value">{{ serverInfo.os }} / {{ serverInfo.arch }}</div>
            </div>
            <div class="server-card">
              <div class="server-card-label">CPU 核心数</div>
              <div class="server-card-value">{{ serverInfo.cpu_cores }}</div>
            </div>
            <div class="server-card">
              <div class="server-card-label">Goroutines</div>
              <div class="server-card-value">{{ serverInfo.goroutines }}</div>
            </div>
            <div class="server-card">
              <div class="server-card-label">内存占用</div>
              <div class="server-card-value">{{ serverInfo.memory_alloc }}</div>
            </div>
            <div class="server-card">
              <div class="server-card-label">数据库大小</div>
              <div class="server-card-value">{{ serverInfo.db_size }}</div>
            </div>
            <div class="server-card server-card--wide">
              <div class="server-card-label">运行时长</div>
              <div class="server-card-value">{{ formatUptime(serverInfo.uptime_seconds) }}</div>
            </div>
          </div>

          <div v-else class="loading-state">
            <span class="loading-text">暂无数据</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  min-height: calc(100vh - 56px - 2.5rem);
}

.page-header {
  margin-bottom: 1.25rem;
}

.page-title {
  font-size: 1.375rem;
  font-weight: 700;
  color: #1e293b;
}

.page-description {
  font-size: 0.8125rem;
  color: #94a3b8;
  margin-top: 0.125rem;
}

.settings-body {
  display: flex;
  gap: 1.25rem;
  flex: 1;
}

.tabs-sidebar {
  width: 190px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 0.1875rem;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5625rem 0.875rem;
  border-radius: 9999px;
  border: none;
  background: transparent;
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  text-align: left;
  width: 100%;
}

.tab-btn:hover {
  background: #e2e8f0;
  color: #334155;
}

.tab-btn--active {
  background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 100%);
  color: #ffffff;
  box-shadow: 0 2px 10px rgba(37, 99, 235, 0.3);
}

.tab-btn-icon {
  width: 17px;
  height: 17px;
  flex-shrink: 0;
}

.tab-btn-text {
  white-space: nowrap;
}

.tab-content {
  flex: 1;
  min-width: 0;
}

.tab-panel {
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

.panel-header {
  margin-bottom: 0.75rem;
}

.panel-title {
  font-size: 1.0625rem;
  font-weight: 600;
  color: #1e293b;
}

.panel-desc {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 0.125rem;
}

.form-card {
  background: #ffffff;
  border-radius: 0.75rem;
  padding: 1.25rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.02);
  border: 1px solid #f1f5f9;
  display: flex;
  flex-direction: column;
  gap: 0.875rem;
}

.form-group {
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

.form-message {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8125rem;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4375rem;
}

.form-message--error {
  color: #b91c1c;
  background: #fef2f2;
}

.form-message--success {
  color: #15803d;
  background: #f0fdf4;
}

.msg-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.form-submit-btn {
  align-self: flex-start;
  margin-top: 0.25rem;
}

.advanced-grid {
  display: grid;
  gap: 0.75rem;
}

.adv-card {
  background: #ffffff;
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.adv-card-left {
  display: flex;
  align-items: center;
  gap: 0.875rem;
}

.adv-icon-wrap {
  width: 38px;
  height: 38px;
  border-radius: 0.625rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.adv-icon-wrap--blue {
  background: #eff6ff;
  color: #2563eb;
}

.adv-icon-wrap--amber {
  background: #fffbeb;
  color: #d97706;
}

.adv-icon-wrap--green {
  background: #f0fdf4;
  color: #16a34a;
}

.adv-icon-wrap--red {
  background: #fef2f2;
  color: #dc2626;
}

.adv-icon {
  width: 19px;
  height: 19px;
}

.adv-card-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
}

.adv-card-desc {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 0.0625rem;
}

.adv-hint {
  font-size: 0.75rem;
  color: #cbd5e1;
  font-weight: 500;
}

.adv-btn--danger {
  color: #dc2626;
  border-color: #fecaca;
}

.adv-btn--danger:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fca5a5;
}

.adv-msg {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  font-size: 0.8125rem;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4375rem;
  margin-top: 0.75rem;
}

.adv-msg--error {
  color: #b91c1c;
  background: #fef2f2;
}

.adv-msg--success {
  color: #15803d;
  background: #f0fdf4;
}

.adv-card--full {
  grid-column: 1 / -1;
}

.qrcode-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.qrcode-header {
  display: flex;
  align-items: center;
  gap: 0.875rem;
}

.qrcode-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.qrcode-input-group {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.qrcode-url-row {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.qrcode-url-row .input-base {
  flex: 1;
}

.qrcode-ip-badge {
  font-size: 0.6875rem;
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  white-space: nowrap;
  flex-shrink: 0;
  background: #dcfce7;
  color: #059669;
  border: 1px solid #bbf7d0;
}

.qrcode-preview {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 0.5rem;
  border: 1px dashed #e2e8f0;
}

.qrcode-image {
  width: 256px;
  height: 256px;
  border-radius: 0.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.server-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.75rem;
}

.server-card {
  background: #ffffff;
  border-radius: 0.625rem;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.server-card--wide {
  grid-column: 1 / -1;
}

.server-card-label {
  font-size: 0.6875rem;
  color: #94a3b8;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.server-card-value {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1e293b;
  word-break: break-all;
}

.loading-state {
  background: #ffffff;
  border-radius: 0.75rem;
  padding: 2.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  border: 1px solid #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: #94a3b8;
  font-size: 0.875rem;
}

.loading-icon {
  width: 18px;
  height: 18px;
}

.loading-text {
  color: #cbd5e1;
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}

.ocr-provider-row {
  display: flex;
  gap: 0.75rem;
  max-width: 480px;
}

.ocr-provider-option {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  padding: 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.75rem;
  cursor: pointer;
  transition: all var(--transition-fast);
  background: #ffffff;
}

.ocr-provider-option:hover {
  border-color: #bfdbfe;
  background: #f8fafc;
}

.ocr-provider-option--active {
  border-color: #2563eb;
  background: #eff6ff;
}

.ocr-provider-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #1e293b;
}

.ocr-provider-desc {
  font-size: 0.6875rem;
  color: #94a3b8;
}

.ocr-provider-option--active .ocr-provider-desc {
  color: #60a5fa;
}

.form-hint {
  font-size: 0.75rem;
  color: #94a3b8;
  margin-top: 0.375rem;
}

.form-link {
  color: #2563eb;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.form-link:hover {
  color: #1d4ed8;
}

@media (max-width: 1024px) {
  .server-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .settings-body {
    flex-direction: column;
  }

  .tabs-sidebar {
    width: 100%;
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: 0.375rem;
  }

  .tab-btn {
    white-space: nowrap;
    flex-shrink: 0;
  }

  .server-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .server-grid {
    grid-template-columns: 1fr;
  }
}
</style>