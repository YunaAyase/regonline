<script setup lang="ts">
const props = defineProps<{
  modelValue: string
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const isMobile = /Mobi|Android|iPhone|iPad|iPod/i.test(navigator.userAgent)

const showPicker = ref(false)
const showKeyboard = ref(false)
const currentMonth = ref(new Date().getMonth())
const currentYear = ref(new Date().getFullYear())
const viewMode = ref<'day' | 'year'>('day')

const weekDays = ['日', '一', '二', '三', '四', '五', '六']

const MIN_YEAR = 1970
const MIN_DATE = new Date(MIN_YEAR, 0, 1)
const MAX_DATE = new Date()

const kYear = ref('')
const kMonth = ref('')
const kDay = ref('')

const parsedDate = computed(() => {
  if (!props.modelValue) return null
  const parts = props.modelValue.split('-')
  if (parts.length !== 3) return null
  const y = parseInt(parts[0] as string)
  const m = parseInt(parts[1] as string)
  const d = parseInt(parts[2] as string)
  if (isNaN(y) || isNaN(m) || isNaN(d)) return null
  return { year: y, month: m - 1, day: d }
})

const isValidDate = computed(() => {
  if (!props.modelValue) return false
  const parts = props.modelValue.split('-')
  if (parts.length !== 3) return false
  const y = parseInt(parts[0] as string)
  const m = parseInt(parts[1] as string)
  const d = parseInt(parts[2] as string)
  if (isNaN(y) || isNaN(m) || isNaN(d)) return false
  if (y < MIN_YEAR || y > MAX_DATE.getFullYear()) return false
  if (m < 1 || m > 12) return false
  if (d < 1 || d > 31) return false
  const dateObj = new Date(y, m - 1, d)
  if (dateObj.getFullYear() !== y || dateObj.getMonth() !== m - 1 || dateObj.getDate() !== d) return false
  const inputDate = new Date(y, m - 1, d)
  return inputDate >= new Date(MIN_YEAR, 0, 1) && inputDate <= MAX_DATE
})

const calendarDays = computed(() => {
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)
  const startDay = firstDay.getDay()
  const totalDays = lastDay.getDate()

  const days: Array<{ date: number; currentMonth: boolean; outOfRange: boolean }> = []

  for (let i = 0; i < startDay; i++) {
    const prevMonthLastDay = new Date(currentYear.value, currentMonth.value, 0).getDate()
    const d = prevMonthLastDay - startDay + i + 1
    const dayDate = new Date(currentYear.value, currentMonth.value - 1, d)
    days.push({ date: d, currentMonth: false, outOfRange: isOutOfRange(dayDate) })
  }

  for (let i = 1; i <= totalDays; i++) {
    const dayDate = new Date(currentYear.value, currentMonth.value, i)
    days.push({ date: i, currentMonth: true, outOfRange: isOutOfRange(dayDate) })
  }

  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const dayDate = new Date(currentYear.value, currentMonth.value + 1, i)
    days.push({ date: i, currentMonth: false, outOfRange: isOutOfRange(dayDate) })
  }

  return days
})

const yearRange = computed(() => {
  const base = Math.floor(currentYear.value / 12) * 12
  const years = Array.from({ length: 12 }, (_, i) => base + i)
  return years.filter(y => y >= MIN_YEAR && y <= MAX_DATE.getFullYear())
})

const monthNames = [
  '一月', '二月', '三月', '四月', '五月', '六月',
  '七月', '八月', '九月', '十月', '十一月', '十二月',
]

function isOutOfRange(date: Date): boolean {
  const min = new Date(MIN_YEAR, 0, 1)
  const max = new Date(MAX_DATE.getFullYear(), MAX_DATE.getMonth(), MAX_DATE.getDate())
  return date < min || date > max
}

function canGoPrev(): boolean {
  if (viewMode.value === 'day') {
    if (currentMonth.value === 0 && currentYear.value <= MIN_YEAR) return false
    return currentYear.value > MIN_YEAR || currentMonth.value > 0
  } else if (viewMode.value === 'year') {
    return currentYear.value > MIN_YEAR
  }
  return true
}

function canGoNext(): boolean {
  if (viewMode.value === 'day') {
    const today = new Date()
    if (currentMonth.value === 11 && currentYear.value >= today.getFullYear()) return false
    return currentYear.value < today.getFullYear() || currentMonth.value < today.getMonth()
  } else if (viewMode.value === 'year') {
    const today = new Date()
    return currentYear.value + 11 < today.getFullYear() || currentYear.value <= today.getFullYear()
  }
  return true
}

function isToday(year: number, month: number, day: number): boolean {
  const today = new Date()
  return today.getFullYear() === year && today.getMonth() === month && today.getDate() === day
}

function isSelectedDay(day: number): boolean {
  if (!parsedDate.value) return false
  return parsedDate.value.year === currentYear.value && parsedDate.value.month === currentMonth.value && parsedDate.value.day === day
}

function isSelectedMonth(): boolean {
  if (!parsedDate.value) return false
  return parsedDate.value.year === currentYear.value
}

function isSelectedYear(): boolean {
  if (!parsedDate.value) return false
  return parsedDate.value.year === currentYear.value
}

function prevPage() {
  if (viewMode.value === 'day') {
    if (currentMonth.value === 0) {
      currentMonth.value = 11
      currentYear.value--
    } else {
      currentMonth.value--
    }
  } else if (viewMode.value === 'year') {
    currentYear.value -= 12
  }
}

function nextPage() {
  if (viewMode.value === 'day') {
    if (currentMonth.value === 11) {
      currentMonth.value = 0
      currentYear.value++
    } else {
      currentMonth.value++
    }
  } else if (viewMode.value === 'year') {
    currentYear.value += 12
  }
}

function selectDay(day: number, isCurrentMonth: boolean, outOfRange: boolean) {
  if (outOfRange) return
  if (!isCurrentMonth) {
    if (day > 15) {
      prevPage()
    } else {
      nextPage()
    }
    setTimeout(() => selectDay(day, true, false), 0)
    return
  }

  const year = currentYear.value
  const month = currentMonth.value + 1
  const dateStr = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
  emit('update:modelValue', dateStr)
}

function selectMonth(month: number) {
  currentMonth.value = month
  viewMode.value = 'day'
}

function selectYear(year: number) {
  currentYear.value = year
  viewMode.value = 'day'
}

function goToToday() {
  const today = new Date()
  currentYear.value = today.getFullYear()
  currentMonth.value = today.getMonth()
  viewMode.value = 'day'
  const dateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
  emit('update:modelValue', dateStr)
}

function onKeyboardInput() {
  if (kYear.value.length === 4 && kMonth.value.length >= 1 && kDay.value.length >= 1) {
    const y = parseInt(kYear.value)
    const m = parseInt(kMonth.value)
    const d = parseInt(kDay.value)
    if (y >= MIN_YEAR && y <= MAX_DATE.getFullYear() && m >= 1 && m <= 12 && d >= 1 && d <= 31) {
      const inputDate = new Date(y, m - 1, d)
      if (inputDate.getFullYear() === y && inputDate.getMonth() === m - 1 && inputDate.getDate() === d) {
        const dateStr = `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
        emit('update:modelValue', dateStr)
        currentYear.value = y
        currentMonth.value = m - 1
      }
    }
  }
}

function closePicker() {
  showPicker.value = false
  viewMode.value = 'day'
}

function openPicker() {
  showPicker.value = true
  viewMode.value = 'day'
}

function onYearInput(e: Event) {
  const v = (e.target as HTMLInputElement).value.replace(/\D/g, '')
  kYear.value = v
  onKeyboardInput()
}

function onMonthInput(e: Event) {
  const v = (e.target as HTMLInputElement).value.replace(/\D/g, '')
  kMonth.value = v
  onKeyboardInput()
}

function onDayInput(e: Event) {
  const v = (e.target as HTMLInputElement).value.replace(/\D/g, '')
  kDay.value = v
  onKeyboardInput()
}

function syncKeyboardFromModel() {
  if (isValidDate.value) {
    const parts = props.modelValue.split('-')
    kYear.value = parts[0] as string
    kMonth.value = String(parseInt(parts[1] as string))
    kDay.value = String(parseInt(parts[2] as string))
  } else {
    kYear.value = ''
    kMonth.value = ''
    kDay.value = ''
  }
}

watch(() => props.modelValue, () => {
  syncKeyboardFromModel()
}, { immediate: true })

watch(showPicker, (val) => {
  if (val) {
    syncKeyboardFromModel()
  }
})

const pickerRef = ref<HTMLElement | null>(null)
const wrapperRef = ref<HTMLElement | null>(null)
const pickerStyle = ref<Record<string, string>>({})

function updatePickerPosition() {
  if (!wrapperRef.value || isMobile) return
  const rect = wrapperRef.value.getBoundingClientRect()
  pickerStyle.value = {
    position: 'fixed',
    top: `${rect.bottom + 8}px`,
    left: `${rect.left}px`,
    width: `${rect.width}px`,
  }
}

function handleClickOutside(e: MouseEvent) {
  if (!isMobile && showPicker.value && pickerRef.value && !pickerRef.value.contains(e.target as Node)) {
    closePicker()
  }
}

onMounted(() => {
  if (!isMobile) {
    document.addEventListener('mousedown', handleClickOutside)
    window.addEventListener('scroll', updatePickerPosition, true)
    window.addEventListener('resize', updatePickerPosition)
  }
})

onUnmounted(() => {
  if (!isMobile) {
    document.removeEventListener('mousedown', handleClickOutside)
    window.removeEventListener('scroll', updatePickerPosition, true)
    window.removeEventListener('resize', updatePickerPosition)
  }
})

function formatDateDisplay(value: string): string {
  if (!value) return ''
  const parts = value.split('-')
  if (parts.length !== 3) return value
  return `${parts[0]}年${parseInt(parts[1] as string)}月${parseInt(parts[2] as string)}日`
}
</script>

<template>
  <div class="date-input-wrapper" ref="wrapperRef">
    <div class="date-input-display" @click="openPicker">
      <template v-if="isValidDate">
        <span class="date-input-value">{{ formatDateDisplay(modelValue) }}</span>
      </template>
      <template v-else>
        <span class="date-input-placeholder">{{ placeholder || '请选择日期' }}</span>
      </template>
      <UIcon name="i-heroicons-calendar" class="date-input-icon" />
    </div>

    <Transition :name="isMobile ? 'bottom-sheet' : 'calendar-fade'">
      <Teleport to="body">
        <div v-if="showPicker" :class="isMobile ? 'calendar-bottom-sheet-overlay' : 'calendar-picker-wrapper'" @click="isMobile ? closePicker() : undefined">
          <div :class="isMobile ? 'calendar-bottom-sheet' : 'calendar-picker'" ref="pickerRef" :style="isMobile ? {} : pickerStyle" @click.stop>
            <div class="calendar-header">
              <div class="calendar-header-top">
                <button type="button" class="calendar-nav-btn" :disabled="!canGoPrev()" @click="prevPage">
                  <UIcon name="i-heroicons-chevron-left" class="w-5 h-5" />
                </button>
                <div class="calendar-header-title">
                  <button
                    type="button"
                    class="calendar-title-part"
                    :class="{ 'calendar-title-part--dimmed': viewMode !== 'year' }"
                    @click="viewMode = viewMode === 'year' ? 'day' : 'year'"
                  >
                    {{ currentYear }}年
                  </button>
                  <button
                    type="button"
                    class="calendar-title-part"
                    :class="{ 'calendar-title-part--dimmed': viewMode !== 'day' }"
                    @click="viewMode = 'day'"
                  >
                    {{ monthNames[currentMonth] }}
                  </button>
                </div>
                <button type="button" class="calendar-nav-btn" :disabled="!canGoNext()" @click="nextPage">
                  <UIcon name="i-heroicons-chevron-right" class="w-5 h-5" />
                </button>
              </div>
              <button type="button" class="calendar-today-btn" @click="goToToday">
                <UIcon name="i-heroicons-clock" class="w-3.5 h-3.5" />
                今天
              </button>
            </div>

            <div class="calendar-mode-toggle">
              <button
                type="button"
                class="calendar-mode-btn"
                :class="{ 'calendar-mode-btn--active': !showKeyboard }"
                @click="showKeyboard = false"
              >
                日历
              </button>
              <button
                type="button"
                class="calendar-mode-btn"
                :class="{ 'calendar-mode-btn--active': showKeyboard }"
                @click="showKeyboard = true"
              >
                输入
              </button>
            </div>

            <div v-if="!showKeyboard && viewMode === 'day'" class="calendar-grid">
              <div v-for="day in weekDays" :key="day" class="calendar-weekday">
                {{ day }}
              </div>
              <div
                v-for="(day, index) in calendarDays"
                :key="index"
                class="calendar-day"
                :class="{
                  'calendar-day--current': day.currentMonth,
                  'calendar-day--today': isToday(currentYear, currentMonth, day.date) && day.currentMonth,
                  'calendar-day--selected': isSelectedDay(day.date) && day.currentMonth,
                  'calendar-day--disabled': day.outOfRange,
                }"
                @click="selectDay(day.date, day.currentMonth, day.outOfRange)"
              >
                {{ day.date }}
              </div>
            </div>

            <div v-else-if="!showKeyboard && viewMode === 'year'" class="calendar-grid calendar-grid--years">
              <div
                v-for="year in yearRange"
                :key="year"
                class="calendar-year-item"
                :class="{
                  'calendar-year-item--current': year === currentYear,
                  'calendar-year-item--today': year === new Date().getFullYear(),
                  'calendar-year-item--selected': isSelectedYear() && year === currentYear,
                  'calendar-year-item--disabled': year < MIN_YEAR || year > MAX_DATE.getFullYear(),
                }"
                @click="selectYear(year)"
              >
                {{ year }}年
              </div>
            </div>

            <div v-if="showKeyboard" class="calendar-keyboard">
              <div class="keyboard-input-group">
                <input
                  :value="kYear"
                  type="text"
                  class="input-glass keyboard-input-segment"
                  placeholder="XXXX"
                  maxlength="4"
                  @input="onYearInput"
                >
                <span class="keyboard-unit">年</span>
                <input
                  :value="kMonth"
                  type="text"
                  class="input-glass keyboard-input-segment"
                  placeholder="XX"
                  maxlength="2"
                  @input="onMonthInput"
                >
                <span class="keyboard-unit">月</span>
                <input
                  :value="kDay"
                  type="text"
                  class="input-glass keyboard-input-segment"
                  placeholder="XX"
                  maxlength="2"
                  @input="onDayInput"
                >
                <span class="keyboard-unit">日</span>
              </div>
            </div>
          </div>
        </div>
      </Teleport>
    </Transition>
  </div>
</template>

<style scoped>
.date-input-wrapper {
  position: relative;
}

.date-input-display {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--input-padding-y) var(--input-padding-x);
  border: 1px solid var(--input-border);
  border-radius: var(--input-radius);
  background: var(--input-bg-glass);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  font-size: var(--input-font-size);
  color: var(--input-text);
  transition: all var(--transition-fast);
  cursor: pointer;
  min-height: var(--select-display-height);
}

.date-input-display:hover {
  border-color: var(--input-border-hover);
  background: rgba(255, 255, 255, 0.95);
}

.date-input-display:focus-within {
  border-color: var(--input-border-focus);
  box-shadow: 0 0 0 3px var(--input-ring-focus);
  background: #ffffff;
}

.date-input-value {
  color: var(--input-text);
  font-weight: 500;
}

.date-input-placeholder {
  color: var(--select-display-placeholder);
}

.date-input-icon {
  width: 20px;
  height: 20px;
  color: #94a3b8;
  flex-shrink: 0;
}

.calendar-picker-wrapper {
  position: fixed;
  z-index: 100000;
}

.calendar-picker {
  background: #ffffff;
  border-radius: 1rem;
  border: 1px solid var(--input-border);
  box-shadow: 0 12px 32px -4px rgba(30, 58, 138, 0.15), 0 4px 12px -4px rgba(30, 58, 138, 0.08);
  padding: 1rem;
  overflow-y: auto;
  max-width: 400px;
}

.calendar-bottom-sheet-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 100000;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.calendar-bottom-sheet {
  width: 100%;
  max-width: 640px;
  max-height: 85vh;
  background: #ffffff;
  border-radius: 1rem 1rem 0 0;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  padding: 1rem;
  padding-bottom: calc(1rem + env(safe-area-inset-bottom, 0px));
}

.calendar-header {
  margin-bottom: 0.75rem;
}

.calendar-header-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.calendar-header-title {
  display: flex;
  align-items: baseline;
  gap: 0.25rem;
  flex: 1;
  justify-content: center;
}

.calendar-title-part {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
  background: transparent;
  border: none;
  padding: 0.25rem 0.5rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all var(--transition-fast);
  font-family: inherit;
}

.calendar-title-part:hover {
  background: #f1f5f9;
  color: #2563eb;
}

.calendar-title-part--dimmed {
  color: #94a3b8;
}

.calendar-title-part--dimmed:hover {
  background: #f1f5f9;
  color: #64748b;
}

.calendar-nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: #64748b;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.calendar-nav-btn:hover {
  background: #f1f5f9;
  color: #2563eb;
}

.calendar-today-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  margin-top: 0.5rem;
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #2563eb;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.calendar-today-btn:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.calendar-mode-toggle {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  padding: 0.25rem;
  background: #f1f5f9;
  border-radius: 9999px;
}

.calendar-mode-btn {
  flex: 1;
  padding: 0.375rem 0.75rem;
  border-radius: 9999px;
  border: none;
  background: transparent;
  color: #64748b;
  font-size: 0.8125rem;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.calendar-mode-btn--active {
  background: #ffffff;
  color: #1e293b;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.calendar-grid--years {
  grid-template-columns: repeat(3, 1fr);
  gap: 0.5rem;
}

.calendar-weekday {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem 0;
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
  text-align: center;
}

.calendar-weekday:first-child {
  color: #ef4444;
}

.calendar-day {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 36px;
  border-radius: 50%;
  font-size: 0.875rem;
  color: #cbd5e1;
  cursor: default;
  transition: all var(--transition-fast);
  position: relative;
}

.calendar-day--current {
  color: #1e293b;
  cursor: pointer;
}

.calendar-day--current:hover:not(.calendar-day--disabled) {
  background: #f1f5f9;
}

.calendar-day--today {
  color: #2563eb;
  font-weight: 700;
}

.calendar-day--today.calendar-day--disabled {
  color: #d1d5db;
  font-weight: 400;
}

.calendar-day--today::after {
  content: '';
  position: absolute;
  bottom: 4px;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #2563eb;
}

.calendar-day--selected {
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: #ffffff;
  font-weight: 600;
}

.calendar-day--selected:hover {
  background: linear-gradient(135deg, #1d4ed8, #0ea5e9);
}

.calendar-day--disabled {
  color: #d1d5db;
  cursor: not-allowed;
  pointer-events: none;
}

.calendar-year-item {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 48px;
  border-radius: 0.75rem;
  font-size: 0.875rem;
  color: #1e293b;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.calendar-year-item:hover {
  background: #f1f5f9;
}

.calendar-year-item--today {
  color: #2563eb;
  font-weight: 700;
}

.calendar-year-item--selected {
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: #ffffff;
  font-weight: 600;
}

.calendar-year-item--selected:hover {
  background: linear-gradient(135deg, #1d4ed8, #0ea5e9);
}

.calendar-year-item--disabled {
  color: #d1d5db;
  cursor: not-allowed;
  pointer-events: none;
}

.calendar-keyboard {
  padding: 0.5rem 0;
}

.keyboard-input-group {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
}

.keyboard-input-segment {
  width: 60px;
  text-align: center;
  font-size: 1rem;
  font-weight: 500;
  padding: 0.5rem 0.25rem;
}

.keyboard-unit {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
}

.calendar-fade-enter-active,
.calendar-fade-leave-active {
  transition: all var(--transition-normal);
}

.calendar-fade-enter-from,
.calendar-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.bottom-sheet-enter-active,
.bottom-sheet-leave-active {
  transition: all 0.3s ease;
}

.bottom-sheet-enter-from,
.bottom-sheet-leave-to {
  opacity: 0;
}

.bottom-sheet-enter-from .calendar-bottom-sheet,
.bottom-sheet-leave-to .calendar-bottom-sheet {
  transform: translateY(100%);
}

.bottom-sheet-enter-to .calendar-bottom-sheet,
.bottom-sheet-leave-from .calendar-bottom-sheet {
  transform: translateY(0);
}

.bottom-sheet-enter-from .calendar-bottom-sheet-overlay,
.bottom-sheet-leave-to .calendar-bottom-sheet-overlay {
  opacity: 0;
}

@media (max-width: 640px) {
  .calendar-picker {
    max-width: calc(100vw - 2rem);
  }
}
</style>
