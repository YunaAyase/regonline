<script setup lang="ts">
import type { RegistrationRequest, ClassItem, OCRNumberCandidate } from '~/types'

const { classes, loading: classesLoading, fetchClasses } = useClasses()
const {
  submitting,
  successMessage,
  errorMessage,
  submitRegistration,
} = useRegistration()
const { siteSettings, siteName, fetchSiteSettings } = useSiteSettings()
const { acquire: acquireWakeLock, release: releaseWakeLock } = useWakeLock()

const config = useRuntimeConfig()

const householdPhoto = ref<File | null>(null)
const householdPreview = ref<string | null>(null)
const idNumberLoading = ref(false)
const idNumberError = ref('')
const showManualInput = ref(false)

const showBottomSheet = ref(false)
const ocrCandidates = ref<OCRNumberCandidate[]>([])

const form = reactive<RegistrationRequest>({
  name: '',
  gender: '',
  birth_date: '',
  grade: '',
  class_id: 0 as any,
  parent_name: '',
  parent_phone: '',
  address: '',
  id_number: '',
})

const step = ref(1)
const started = ref(false)

const steps = [
  { num: 1, label: '个人信息', icon: 'i-heroicons-user' },
  { num: 2, label: '家长信息', icon: 'i-heroicons-user-group' },
  { num: 3, label: '户口本与身份证', icon: 'i-heroicons-document-check' },
  { num: 4, label: '确认提交', icon: 'i-heroicons-check-circle' },
]

const gradeOptions = [
  { label: '请选择年级', value: '', disabled: true },
  { label: '学前', value: '学前' },
  { label: '一年级', value: '一年级' },
  { label: '二年级', value: '二年级' },
  { label: '三年级', value: '三年级' },
  { label: '四年级', value: '四年级' },
  { label: '五年级', value: '五年级' },
  { label: '六年级', value: '六年级' },
  { label: '七年级', value: '七年级' },
  { label: '八年级', value: '八年级' },
  { label: '九年级', value: '九年级' },
  { label: '高一', value: '高一' },
  { label: '高二', value: '高二' },
  { label: '高三', value: '高三' },
  { label: '大学', value: '大学' },
]

const classOptions = computed(() => {
  const disabledOpt = { label: '请选择班级', value: 0, disabled: true }
  const opts = classes.value
    .filter(cls => cls.enabled)
    .map(cls => ({
      label: `${cls.name} (${cls.current_count}/${cls.max_students})`,
      value: cls.id,
      description: cls.description,
      min_age: cls.min_age,
      max_age: cls.max_age,
    }))
  return [disabledOpt, ...opts]
})

onMounted(async () => {
  await fetchClasses()
  fetchSiteSettings()
  acquireWakeLock(config.public.apiBase as string)
})

onUnmounted(() => {
  releaseWakeLock()
})

function startRegistration() {
  started.value = true
}

async function handleHouseholdPhotoChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    householdPhoto.value = file
    if (householdPreview.value) URL.revokeObjectURL(householdPreview.value)
    householdPreview.value = URL.createObjectURL(file)
    idNumberError.value = ''
    autoRecognizeID(file)
  }
}

async function autoRecognizeID(file: File) {
  idNumberLoading.value = true
  idNumberError.value = ''
  showManualInput.value = false
  try {
    const formData = new FormData()
    formData.append('photo', file)
    const res = await $fetch<{ code: number; data: { id_number: string; candidates: OCRNumberCandidate[] }; message: string }>(
      `${config.public.apiBase}/ocr/recognize`,
      {
        method: 'POST',
        body: formData,
      }
    )
    if (res.code === 0 && res.data?.candidates && res.data.candidates.length > 0) {
      ocrCandidates.value = res.data.candidates
      
      const isMobile = /Mobi|Android|iPhone|iPad|iPod/i.test(navigator.userAgent)
      if (isMobile) {
        showBottomSheet.value = true
        
        const idCandidate = res.data.candidates.find(c => c.is_id_number)
        if (idCandidate) {
          form.id_number = idCandidate.value
        }
      } else {
        if (res.data.id_number) {
          form.id_number = res.data.id_number
        } else {
          idNumberError.value = '未识别到身份证号，请手动输入'
          showManualInput.value = true
        }
      }
    } else {
      idNumberError.value = res.message || 'OCR 识别失败，请手动输入身份证号'
      showManualInput.value = true
    }
  } catch {
    idNumberError.value = 'OCR 识别失败，请手动输入身份证号'
    showManualInput.value = true
  } finally {
    idNumberLoading.value = false
  }
}

function selectCandidate(candidate: OCRNumberCandidate) {
  form.id_number = candidate.value
  showBottomSheet.value = false
}

function closeBottomSheet() {
  showBottomSheet.value = false
}

function removeHouseholdPhoto() {
  if (householdPreview.value) {
    URL.revokeObjectURL(householdPreview.value)
  }
  householdPhoto.value = null
  householdPreview.value = null
  form.id_number = ''
  idNumberError.value = ''
  showManualInput.value = false
}

function toggleManualInput() {
  showManualInput.value = !showManualInput.value
  if (!showManualInput.value) {
    form.id_number = ''
    idNumberError.value = ''
  }
}

function onPhoneInput(event: Event) {
  const target = event.target as HTMLInputElement
  let value = target.value.replace(/\D/g, '')
  if (value.length > 11) value = value.substring(0, 11)
  if (value.length > 3 && value.length <= 7) {
    value = value.substring(0, 3) + ' ' + value.substring(3)
  } else if (value.length > 7) {
    value = value.substring(0, 3) + ' ' + value.substring(3, 7) + ' ' + value.substring(7)
  }
  form.parent_phone = value
}

function validateStep1() {
  if (!form.name || !form.gender || !form.birth_date || !form.grade || !form.class_id) {
    return false
  }
  return true
}

function validateStep2() {
  if (!form.parent_name || !form.parent_phone || !form.address) {
    return false
  }
  return true
}

function validateStep3() {
  if (!householdPhoto.value || !form.id_number) return false
  if (form.id_number.length !== 18) return false
  const idNumber = form.id_number as string
  const birthDate = form.birth_date as string
  const birthStr = birthDate.replace(/-/g, '')
  if (idNumber.substring(6, 14) !== birthStr) return false
  const genderDigit = parseInt(idNumber.charAt(16))
  const genderFromID = genderDigit % 2 === 1 ? '男' : '女'
  if (form.gender !== genderFromID) return false
  return true
}

function nextStep() {
  if (step.value === 1 && !validateStep1()) return
  if (step.value === 2 && !validateStep2()) return
  if (step.value === 3 && !validateStep3()) return
  if (step.value < 4) {
    step.value++
  }
}

function prevStep() {
  if (step.value > 1) {
    step.value--
  }
}

async function handleSubmit() {
  await submitRegistration({ ...form }, householdPhoto.value || undefined)

  if (!errorMessage.value) {
    form.name = ''
    form.gender = ''
    form.birth_date = ''
    form.grade = ''
    form.class_id = 0 as any
    form.parent_name = ''
    form.parent_phone = ''
    form.address = ''
    form.id_number = ''
    removeHouseholdPhoto()
    step.value = 1
  }
}
</script>

<template>
  <div class="page-wrapper">
    <header class="header-section">
      <div class="header-inner">
        <div class="header-brand">
          <div class="header-logo">
            <UIcon name="i-heroicons-academic-cap" class="w-8 h-8 text-white" />
          </div>
          <div>
            <h1 class="header-title">{{ siteName }}</h1>
            <p class="header-subtitle">Online Registration System</p>
          </div>
        </div>
        <div class="header-actions">
          <NuxtLink to="/admin" class="btn-pill-outline">
            <UIcon name="i-heroicons-arrow-right-on-rectangle" class="w-3.5 h-3.5" />
            登录
          </NuxtLink>
        </div>
      </div>
    </header>

    <main class="main-container">
      <div class="content-wrapper">
        <div v-if="!started" class="welcome-page">
          <div class="welcome-card">
            <div class="welcome-icon">
              <UIcon name="i-heroicons-academic-cap" class="w-16 h-16" />
            </div>
            <div class="welcome-title-group">
              <span class="welcome-subtitle">欢迎进入</span>
              <h1 class="welcome-title">{{ siteSettings?.site_name || '在线报名系统' }}</h1>
            </div>
            <p v-if="siteSettings?.site_description" class="welcome-desc">
              {{ siteSettings.site_description }}
            </p>
            <button class="btn-pill btn-lg" @click="startRegistration">
              <span>开始报名</span>
              <UIcon name="i-heroicons-arrow-right" class="w-5 h-5" />
            </button>
          </div>
        </div>

        <template v-else>
          <div class="stepper">
          <div
            v-for="(s, index) in steps"
            :key="s.num"
            class="stepper-item"
          >
            <div
              class="stepper-circle"
              :class="[step >= s.num ? 'stepper-circle--active' : '', step === s.num ? 'stepper-circle--current' : '']"
            >
              <UIcon v-if="step > s.num" name="i-heroicons-check" class="w-4 h-4" />
              <UIcon v-else-if="step === s.num" :name="s.icon" class="w-4 h-4" />
              <span v-else>{{ s.num }}</span>
            </div>
            <span class="stepper-label" :class="step >= s.num ? 'stepper-label--active' : ''">
              {{ s.label }}
            </span>
            <div
              v-if="index < steps.length - 1"
              class="stepper-line"
              :class="step > s.num ? 'stepper-line--active' : ''"
            />
          </div>
        </div>

        <div class="form-card">
          <div v-if="step === 1">
            <div class="section-header">
              <UIcon name="i-heroicons-user" class="section-icon" />
              <h2 class="section-title">个人信息</h2>
              <span class="section-step">步骤 1/4</span>
            </div>
            <div class="form-grid">
              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 姓名</label>
                <input v-model="form.name" type="text" class="input-glass" placeholder="请输入姓名">
              </div>

              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 性别</label>
                <SelectInput
                  v-model="form.gender"
                  :options="[{ label: '请选择性别', value: '', disabled: true }, { label: '男', value: '男' }, { label: '女', value: '女' }]"
                  placeholder="请选择性别"
                />
              </div>

              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 出生日期</label>
                <DateInput v-model="form.birth_date" placeholder="请选择出生日期" />
              </div>

              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 年级</label>
                <SelectInput
                  v-model="form.grade"
                  :options="gradeOptions"
                  placeholder="请选择年级"
                />
              </div>

              <div class="form-field form-field--full">
                <label class="field-label"><span class="required-star">*</span> 班级</label>
                <SelectInput
                  v-model="form.class_id"
                  :options="classOptions"
                  placeholder="请选择班级"
                />
              </div>
            </div>
          </div>

          <div v-if="step === 2">
            <div class="section-header">
              <UIcon name="i-heroicons-user-group" class="section-icon" />
              <h2 class="section-title">家长信息</h2>
              <span class="section-step">步骤 2/4</span>
            </div>
            <div class="form-grid">
              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 家长姓名</label>
                <input v-model="form.parent_name" type="text" class="input-glass" placeholder="请输入家长姓名">
              </div>

              <div class="form-field">
                <label class="field-label"><span class="required-star">*</span> 家长联系电话</label>
                <div class="phone-input-wrapper">
                  <span class="phone-prefix">+86</span>
                  <input
                    v-model="form.parent_phone"
                    type="tel"
                    class="input-glass phone-input"
                    placeholder="请输入手机号"
                    maxlength="13"
                    @input="onPhoneInput"
                  >
                </div>
              </div>

              <div class="form-field form-field--full">
                <label class="field-label"><span class="required-star">*</span> 家庭住址</label>
                <textarea v-model="form.address" class="input-glass input-textarea" placeholder="请输入详细家庭住址" rows="3" />
              </div>
            </div>
          </div>

          <div v-if="step === 3">
            <div class="section-header">
              <UIcon name="i-heroicons-document-check" class="section-icon" />
              <h2 class="section-title">户口本与身份证</h2>
              <span class="section-step">步骤 3/4</span>
            </div>

            <div class="id-section">
              <h3 class="id-section-title">户口本照片</h3>
              <div v-if="householdPreview" class="id-photo-area">
                <img :src="householdPreview" alt="户口本照片" class="id-photo-img">
                <button type="button" class="id-photo-remove" @click="removeHouseholdPhoto">
                  <UIcon name="i-heroicons-trash" class="w-4 h-4" />
                  删除照片
                </button>
              </div>
              <div v-else class="id-photo-upload">
                <UIcon name="i-heroicons-camera" class="id-photo-icon" />
                <p class="id-photo-text">拍摄或上传户口本照片</p>
                <p class="id-photo-hint">系统将自动识别身份证号</p>
                <div class="id-photo-actions">
                  <label class="id-photo-btn id-photo-btn--camera">
                    <UIcon name="i-heroicons-camera" class="w-4 h-4" />
                    拍照
                    <input
                      type="file"
                      accept="image/*"
                      capture="environment"
                      class="sr-only"
                      @change="handleHouseholdPhotoChange"
                    >
                  </label>
                  <label class="id-photo-btn id-photo-btn--gallery">
                    <UIcon name="i-heroicons-photo" class="w-4 h-4" />
                    选择照片
                    <input
                      type="file"
                      accept="image/*"
                      class="sr-only"
                      @change="handleHouseholdPhotoChange"
                    >
                  </label>
                </div>
              </div>
            </div>

            <div class="id-section id-section--inline">
              <h3 class="id-section-title">身份证号</h3>
              <div v-if="idNumberLoading" class="id-number-loading">
                <UIcon name="i-heroicons-arrow-path" class="w-5 h-5 animate-spin text-blue-500" />
                <span>OCR 识别中...</span>
              </div>
              <template v-else>
                <div class="id-number-row">
                  <input
                    v-model="form.id_number"
                    type="text"
                    class="input-glass"
                    :class="{ 'input-glass--error': idNumberError }"
                    placeholder="请输入18位身份证号"
                    maxlength="18"
                    :readonly="!showManualInput && !!form.id_number && !idNumberError"
                  >
                  <button
                    v-if="householdPhoto && !idNumberLoading"
                    type="button"
                    class="btn-pill-outline btn-sm"
                    @click="toggleManualInput"
                  >
                    {{ showManualInput ? '隐藏' : '手动输入' }}
                  </button>
                </div>
                <div v-if="idNumberError" class="id-number-error">
                  <UIcon name="i-heroicons-exclamation-circle" class="w-4 h-4" />
                  <span>{{ idNumberError }}</span>
                </div>
              </template>
            </div>
          </div>

          <div v-if="step === 4">
            <div class="section-header">
              <UIcon name="i-heroicons-check-circle" class="section-icon" />
              <h2 class="section-title">确认信息</h2>
              <span class="section-step">步骤 4/4</span>
            </div>

            <div class="review-section">
              <h3 class="review-title">个人信息</h3>
              <div class="review-grid">
                <div class="review-item">
                  <span class="review-label">姓名</span>
                  <span class="review-value">{{ form.name }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">性别</span>
                  <span class="review-value">{{ form.gender }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">出生日期</span>
                  <span class="review-value">{{ form.birth_date }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">身份证号</span>
                  <span class="review-value">{{ form.id_number }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">年级</span>
                  <span class="review-value">{{ form.grade }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">班级</span>
                  <span class="review-value">
                    {{ classes.find(c => c.id === form.class_id)?.name || '-' }}
                  </span>
                </div>
              </div>
            </div>

            <div class="review-section">
              <h3 class="review-title">家长信息</h3>
              <div class="review-grid">
                <div class="review-item">
                  <span class="review-label">家长姓名</span>
                  <span class="review-value">{{ form.parent_name }}</span>
                </div>
                <div class="review-item">
                  <span class="review-label">联系电话</span>
                  <span class="review-value">{{ form.parent_phone }}</span>
                </div>
                <div class="review-item review-item--full">
                  <span class="review-label">家庭住址</span>
                  <span class="review-value">{{ form.address }}</span>
                </div>
              </div>
            </div>

            <div class="review-section">
              <h3 class="review-title">户口本照片</h3>
              <div v-if="householdPreview" class="photo-preview-review">
                <img :src="householdPreview" alt="户口本照片">
              </div>
              <div v-else class="photo-placeholder">
                <UIcon name="i-heroicons-photo" class="w-8 h-8 text-gray-400" />
                <span class="text-gray-400">未上传照片</span>
              </div>
            </div>

            <div v-if="successMessage" class="alert alert--success">
              <UIcon name="i-heroicons-check-circle" class="alert-icon" />
              <span>{{ successMessage }}</span>
            </div>
            <div v-if="errorMessage" class="alert alert--error">
              <UIcon name="i-heroicons-exclamation-circle" class="alert-icon" />
              <span>{{ errorMessage }}</span>
            </div>
          </div>
        </div>

        <div class="action-bar">
          <button v-if="step > 1" class="btn-pill-neutral" @click="prevStep">
            <UIcon name="i-heroicons-arrow-left" class="w-4 h-4" />
            上一步
          </button>
          <div class="action-spacer" />
          <button v-if="step < 4" class="btn-pill action-btn" :disabled="step === 1 ? !validateStep1() : step === 2 ? !validateStep2() : step === 3 ? !validateStep3() : false" @click="nextStep">
            下一步
            <UIcon name="i-heroicons-arrow-right" class="w-4 h-4" />
          </button>
          <button v-if="step === 4" class="btn-pill action-btn" :disabled="submitting" @click="handleSubmit">
            <UIcon v-if="submitting" name="i-heroicons-arrow-path" class="w-4 h-4 animate-spin" />
            <template v-else>
              提交报名
              <UIcon name="i-heroicons-paper-airplane" class="w-4 h-4" />
            </template>
          </button>
        </div>
        </template>
      </div>
    </main>

    <footer class="footer-section">
      <p v-if="siteSettings?.icp_record || siteSettings?.copyright" class="footer-text">
        <span v-if="siteSettings?.icp_record">{{ siteSettings.icp_record }}</span>
        <span v-if="siteSettings?.icp_record && siteSettings?.copyright" class="footer-sep">·</span>
        <span v-if="siteSettings?.copyright">{{ siteSettings.copyright }}</span>
      </p>
    </footer>

    <Teleport to="body">
      <Transition name="bottom-sheet">
        <div v-if="showBottomSheet" class="bottom-sheet-overlay" @click="closeBottomSheet">
          <div class="bottom-sheet" @click.stop>
            <div class="bottom-sheet-header">
              <div class="bottom-sheet-indicator" />
              <h3 class="bottom-sheet-title">选择识别到的数字</h3>
              <p class="bottom-sheet-desc">请点击选择正确的身份证号</p>
            </div>

            <div class="bottom-sheet-content">
              <div v-if="ocrCandidates.length === 0" class="bottom-sheet-empty">
                <UIcon name="i-heroicons-exclamation-triangle" class="w-8 h-8 text-gray-400" />
                <p>未识别到数字，请手动输入</p>
              </div>

              <div v-else class="candidate-list">
                <button
                  v-for="(candidate, index) in ocrCandidates"
                  :key="index"
                  class="candidate-item"
                  :class="{ 'candidate-item--id': candidate.is_id_number }"
                  @click="selectCandidate(candidate)"
                >
                  <div class="candidate-main">
                    <span class="candidate-value">{{ candidate.value }}</span>
                    <span class="candidate-label">{{ candidate.label }}</span>
                  </div>
                  <UIcon v-if="candidate.is_id_number" name="i-heroicons-check-circle" class="candidate-badge" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.page-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header-section {
  background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 35%, #3b82f6 65%, #38bdf8 100%);
  padding: 1.5rem 0;
  box-shadow: 0 4px 20px -4px rgba(30, 58, 138, 0.15);
}

.header-inner {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-brand {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-actions {
  display: flex;
  align-items: center;
}

.header-logo {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  line-height: 1.2;
}

.header-subtitle {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.7);
  font-weight: 400;
}

.main-container {
  flex: 1;
  padding: 2rem 1rem;
}

.content-wrapper {
  max-width: 680px;
  margin: 0 auto;
}

.welcome-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 200px);
  padding: 2rem 1.5rem;
}

.welcome-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
  padding: 3rem 2.5rem;
  background: var(--card-bg-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--card-border);
  border-radius: 1.5rem;
  box-shadow: 0 8px 32px -4px rgba(30, 58, 138, 0.12);
  text-align: center;
  max-width: 480px;
  width: 100%;
}

.welcome-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1e3a8a, #3b82f6, #38bdf8);
  color: #ffffff;
  box-shadow: 0 4px 16px -2px rgba(37, 99, 235, 0.3);
}

.welcome-title-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
}

.welcome-subtitle {
  font-size: 0.875rem;
  color: #94a3b8;
  font-weight: 400;
  letter-spacing: 0.05em;
}

.welcome-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: #1e293b;
  line-height: 1.3;
}

.welcome-desc {
  font-size: 1rem;
  color: #64748b;
  line-height: 1.6;
  max-width: 360px;
}

.btn-lg {
  padding: 0.75rem 2.5rem;
  font-size: 1.125rem;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.stepper {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 2rem;
  padding: 0 1rem;
}

.stepper-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  position: relative;
}

.stepper-circle {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #e5e7eb;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: 600;
  flex-shrink: 0;
  transition: all 0.25s ease;
}

.stepper-circle--active {
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: #ffffff;
}

.stepper-circle--current {
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.15);
}

.stepper-label {
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
  white-space: nowrap;
  transition: color 0.25s ease;
}

.stepper-label--active {
  color: #1e3a8a;
}

.stepper-line {
  width: 32px;
  height: 2px;
  background: #e5e7eb;
  margin: 0 0.5rem;
  transition: background 0.25s ease;
}

.stepper-line--active {
  background: linear-gradient(90deg, #2563eb, #38bdf8);
}

.form-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(16px);
  border-radius: 1rem;
  border: 1px solid rgba(255, 255, 255, 0.7);
  box-shadow: 0 4px 24px -4px rgba(30, 58, 138, 0.08), 0 1px 4px -1px rgba(30, 58, 138, 0.04);
  padding: 2rem;
  margin-bottom: 1.5rem;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.section-icon {
  width: 24px;
  height: 24px;
  color: #2563eb;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1e293b;
  flex: 1;
}

.section-step {
  font-size: 0.75rem;
  color: #6b7280;
  background: #f1f5f9;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-weight: 500;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.25rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field--full {
  grid-column: 1 / -1;
}

.field-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #374151;
}

.required-star {
  color: #ef4444;
  margin-right: 1px;
}

.id-section {
  margin-bottom: 1.5rem;
}

.id-section:last-child {
  margin-bottom: 0;
}

.id-section--inline {
  display: flex;
  flex-direction: column;
}

.id-section-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #2563eb;
  margin-bottom: 0.75rem;
  padding-left: 0.5rem;
  border-left: 3px solid #2563eb;
}

.id-photo-area {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 0.75rem;
  border: 1px solid #e5e7eb;
}

.id-photo-img {
  max-width: 120px;
  max-height: 80px;
  border-radius: 0.5rem;
  border: 1px solid #d1d5db;
}

.id-photo-remove {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.75rem;
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  font-family: inherit;
}

.id-photo-remove:hover {
  background: #fecaca;
}

.id-photo-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 2rem 1.5rem;
  background: #f8fafc;
  border-radius: 0.75rem;
  border: 2px dashed #d1d5db;
  text-align: center;
  transition: all 0.2s ease;
}

.id-photo-upload:hover {
  border-color: #2563eb;
  background: #eff6ff;
}

.id-photo-icon {
  width: 40px;
  height: 40px;
  color: #9ca3af;
}

.id-photo-text {
  font-size: 0.9375rem;
  color: #374151;
  font-weight: 500;
}

.id-photo-hint {
  font-size: 0.8125rem;
  color: #9ca3af;
}

.id-photo-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-top: 1rem;
}

.id-photo-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.625rem 1.5rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.3);
}

.id-photo-btn--camera {
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: #ffffff;
}

.id-photo-btn--camera:hover {
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.4);
  transform: translateY(-1px);
}

.id-photo-btn--gallery {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #ffffff;
}

.id-photo-btn--gallery:hover {
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
  transform: translateY(-1px);
}

.phone-input-wrapper {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.phone-prefix {
  font-size: 0.875rem;
  color: #6b7280;
  font-weight: 500;
  white-space: nowrap;
  padding-left: 0.5rem;
}

.phone-input {
  flex: 1;
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
  border: 0;
}

.id-number-loading {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: #eff6ff;
  border-radius: 0.625rem;
  font-size: 0.875rem;
  color: #2563eb;
  font-weight: 500;
}

.id-number-row {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.id-number-row .input-glass {
  flex: 1;
}

.input-glass--error {
  border-color: #ef4444 !important;
}

.id-number-error {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  margin-top: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: #fef2f2;
  color: #dc2626;
  border-radius: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: 0.8125rem;
  white-space: nowrap;
  flex-shrink: 0;
}

.review-section {
  margin-bottom: 1.5rem;
}

.review-section:last-child {
  margin-bottom: 0;
}

.review-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #2563eb;
  margin-bottom: 0.75rem;
  padding-left: 0.5rem;
  border-left: 3px solid #2563eb;
}

.review-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.review-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border-radius: 0.625rem;
}

.review-item--full {
  grid-column: 1 / -1;
}

.review-label {
  font-size: 0.75rem;
  color: #6b7280;
  font-weight: 500;
}

.review-value {
  font-size: 0.875rem;
  color: #1e293b;
  font-weight: 500;
}

.photo-preview-review {
  display: flex;
  justify-content: center;
}

.photo-preview-review img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 0.75rem;
  border: 1px solid #e5e7eb;
}

.photo-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 2rem;
  background: #f8fafc;
  border-radius: 0.75rem;
  border: 1px dashed #d1d5db;
  justify-content: center;
}

.alert {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  margin-top: 1.25rem;
}

.alert--success {
  background: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.alert--error {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fecaca;
}

.alert-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.action-bar {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0 0.5rem;
}

.action-spacer {
  flex: 1;
}

.action-btn {
  min-width: 140px;
  font-size: 0.9375rem;
  padding: 0.5625rem 1.5rem;
}

.action-btn:disabled {
  opacity: 0.5;
}

.footer-section {
  text-align: center;
  padding: 1.5rem;
  color: #9ca3af;
  font-size: 0.75rem;
}

.footer-text {
  display: inline-flex;
  align-items: center;
  gap: 0;
}

.footer-sep {
  margin: 0 0.625rem;
  font-weight: 300;
}

@media (max-width: 640px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .review-grid {
    grid-template-columns: 1fr;
  }

  .form-card {
    padding: 1.25rem;
  }

  .stepper-label {
    display: none;
  }

  .stepper-line {
    width: 24px;
  }

  .header-title {
    font-size: 1.25rem;
  }

  .id-number-row {
    flex-direction: column;
    align-items: stretch;
  }

  .id-photo-area {
    flex-direction: column;
    text-align: center;
  }

  .id-photo-actions {
    flex-direction: column;
    gap: 0.75rem;
  }

  .id-photo-btn {
    justify-content: center;
    width: 100%;
  }

  .phone-input-wrapper {
    flex-wrap: wrap;
  }
}

.bottom-sheet-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 100000;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.bottom-sheet {
  width: 100%;
  max-width: 640px;
  max-height: 80vh;
  background: #ffffff;
  border-radius: 1rem 1rem 0 0;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.bottom-sheet-header {
  padding: 1rem 1.25rem 0.75rem;
  text-align: center;
  border-bottom: 1px solid #f1f5f9;
  flex-shrink: 0;
}

.bottom-sheet-indicator {
  width: 36px;
  height: 4px;
  background: #d1d5db;
  border-radius: 9999px;
  margin: 0 auto 0.75rem;
}

.bottom-sheet-title {
  font-size: 1.0625rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.bottom-sheet-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
}

.bottom-sheet-content {
  flex: 1;
  overflow-y: auto;
  padding: 1rem 1.25rem;
  padding-bottom: calc(1rem + env(safe-area-inset-bottom, 0px));
}

.bottom-sheet-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 2rem 0;
  color: #94a3b8;
  font-size: 0.875rem;
}

.candidate-list {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
}

.candidate-item {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.875rem 1rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 0.75rem;
  cursor: pointer;
  transition: all 0.15s ease;
  width: 100%;
  text-align: left;
  font-family: inherit;
}

.candidate-item:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.candidate-item--id {
  background: #f0fdf4;
  border-color: #bbf7d0;
}

.candidate-item--id:hover {
  background: #dcfce7;
  border-color: #86efac;
}

.candidate-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.candidate-value {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1e293b;
  font-family: 'SF Mono', 'Cascadia Code', 'Consolas', monospace;
  word-break: break-all;
}

.candidate-label {
  font-size: 0.75rem;
  color: #64748b;
}

.candidate-item--id .candidate-value {
  color: #059669;
}

.candidate-item--id .candidate-label {
  color: #10b981;
}

.candidate-badge {
  width: 20px;
  height: 20px;
  color: #10b981;
  flex-shrink: 0;
}

.bottom-sheet-enter-active,
.bottom-sheet-leave-active {
  transition: all 0.3s ease;
}

.bottom-sheet-enter-from,
.bottom-sheet-leave-to {
  opacity: 0;
}

.bottom-sheet-enter-from .bottom-sheet,
.bottom-sheet-leave-to .bottom-sheet {
  transform: translateY(100%);
}

.bottom-sheet-enter-to .bottom-sheet,
.bottom-sheet-leave-from .bottom-sheet {
  transform: translateY(0);
}

.bottom-sheet-enter-from .bottom-sheet-overlay,
.bottom-sheet-leave-to .bottom-sheet-overlay {
  opacity: 0;
}

@media (min-width: 641px) {
  .bottom-sheet {
    border-radius: 1rem;
    margin-bottom: 2rem;
  }
}
</style>
