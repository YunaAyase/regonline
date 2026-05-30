<script setup lang="ts">
interface SelectOption {
  label: string
  value: string | number
  disabled?: boolean
  description?: string
  min_age?: number
  max_age?: number
}

const props = defineProps<{
  modelValue: string | number | null
  options: SelectOption[]
  placeholder?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const isOpen = ref(false)
const triggerRef = ref<HTMLElement | null>(null)
const dropdownStyle = ref<Record<string, string>>({})
const isMobile = ref(false)

const selectedLabel = computed(() => {
  const opt = props.options.find(o => o.value === props.modelValue)
  return opt ? opt.label : null
})

const enabledOptions = computed(() =>
  props.options.filter(o => !o.disabled)
)

function checkMobile() {
  isMobile.value = window.innerWidth < 768
}

function calcPosition() {
  if (!triggerRef.value) return
  const rect = triggerRef.value.getBoundingClientRect()
  dropdownStyle.value = {
    position: 'fixed',
    top: `${rect.bottom + 4}px`,
    left: `${rect.left}px`,
    width: `${rect.width}px`,
    zIndex: '100000',
  }
}

function open() {
  isOpen.value = true
  if (!isMobile.value) {
    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        calcPosition()
      })
    })
  }
}

function close() {
  isOpen.value = false
  dropdownStyle.value = {}
}

function toggle() {
  if (isOpen.value) {
    close()
  } else {
    open()
  }
}

function selectOption(opt: SelectOption) {
  if (opt.disabled) return
  emit('update:modelValue', opt.value)
  close()
}

function handleClickOutside(e: MouseEvent) {
  if (!isOpen.value || isMobile.value) return
  const triggerEl = triggerRef.value
  if (!triggerEl || triggerEl.contains(e.target as Node)) return
  close()
}

watch(isOpen, (val) => {
  if (val && !isMobile.value) {
    requestAnimationFrame(() => {
      document.addEventListener('click', handleClickOutside, true)
      window.addEventListener('scroll', calcPosition, true)
      window.addEventListener('resize', calcPosition)
    })
  } else {
    document.removeEventListener('click', handleClickOutside, true)
    window.removeEventListener('scroll', calcPosition, true)
    window.removeEventListener('resize', calcPosition)
  }
})

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside, true)
  window.removeEventListener('scroll', calcPosition, true)
  window.removeEventListener('resize', calcPosition)
  window.removeEventListener('resize', checkMobile)
})
</script>

<template>
  <div class="custom-select">
    <div ref="triggerRef" class="trigger-wrapper">
      <button
        type="button"
        class="select-trigger"
        :class="{ 'select-trigger--open': isOpen }"
        @click="toggle"
      >
        <span
          class="select-display"
          :class="{ 'select-display--placeholder': !selectedLabel }"
        >
          {{ selectedLabel || props.placeholder || '请选择' }}
        </span>
        <UIcon
          name="i-heroicons-chevron-down"
          class="select-arrow"
          :class="{ 'select-arrow--open': isOpen }"
        />
      </button>
    </div>

    <Teleport to="body">
      <Transition name="select-dropdown">
        <div
          v-if="isOpen && !isMobile"
          class="select-dropdown select-dropdown--desktop"
          :style="dropdownStyle"
          @click.stop
        >
          <div
            v-for="(opt, idx) in options"
            :key="idx"
            class="select-option"
            :class="{
              'select-option--active': opt.value === props.modelValue,
              'select-option--disabled': opt.disabled,
            }"
            @click="selectOption(opt)"
          >
            <div class="select-option-content">
              <span class="select-option-label">{{ opt.label }}</span>
              <span v-if="opt.description || (opt.min_age !== undefined && opt.max_age !== undefined)" class="select-option-desc">
                <span v-if="opt.description">{{ opt.description }}</span>
                <span v-if="opt.description && (opt.min_age !== undefined && opt.max_age !== undefined)"> | </span>
                <span v-if="opt.min_age !== undefined && opt.max_age !== undefined">{{ opt.min_age }} ~ {{ opt.max_age }}岁</span>
              </span>
            </div>
          </div>
        </div>
      </Transition>

      <Transition name="select-modal">
        <div
          v-if="isOpen && isMobile"
          class="select-modal"
          @click.stop
        >
          <div class="select-modal__overlay" @click="close" />
          <div class="select-modal__panel">
            <div class="select-modal__header">
              <button type="button" class="select-modal__cancel" @click="close">
                取消
              </button>
              <span class="select-modal__title">{{ props.placeholder || '请选择' }}</span>
              <button type="button" class="select-modal__confirm" @click="close">
                完成
              </button>
            </div>
            <div class="select-modal__options">
              <div
                v-for="(opt, idx) in options"
                :key="idx"
                class="select-modal__option"
                :class="{
                  'select-modal__option--active': opt.value === props.modelValue,
                  'select-modal__option--disabled': opt.disabled,
                }"
                @click="selectOption(opt)"
              >
                <div class="select-modal__option-content">
                  <span class="select-modal__option-label">{{ opt.label }}</span>
                  <span v-if="opt.description || (opt.min_age !== undefined && opt.max_age !== undefined)" class="select-modal__option-desc">
                    <span v-if="opt.description">{{ opt.description }}</span>
                    <span v-if="opt.description && (opt.min_age !== undefined && opt.max_age !== undefined)"> · </span>
                    <span v-if="opt.min_age !== undefined && opt.max_age !== undefined">{{ opt.min_age }} ~ {{ opt.max_age }}岁</span>
                  </span>
                </div>
                <UIcon
                  v-if="opt.value === props.modelValue"
                  name="i-heroicons-check"
                  class="select-modal__option-icon"
                />
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.custom-select {
  width: 100%;
}

.trigger-wrapper {
  width: 100%;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: var(--input-padding-y) var(--input-padding-x);
  border: 1px solid var(--input-border);
  border-radius: var(--input-radius);
  background: var(--input-bg-glass);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  cursor: pointer;
  transition: all var(--transition-fast);
  outline: none;
  line-height: 1.5;
  font-size: var(--input-font-size);
  font-family: inherit;
  text-align: left;
}

.select-trigger:hover {
  border-color: var(--input-border-hover);
  background: rgba(255, 255, 255, 0.95);
}

.select-trigger--open {
  border-color: var(--input-border-focus);
  box-shadow: 0 0 0 3px var(--input-ring-focus);
  background: #ffffff;
}

.select-display {
  color: var(--select-display-text);
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-display--placeholder {
  color: var(--select-display-placeholder);
}

.select-arrow {
  width: 16px;
  height: 16px;
  color: #94a3b8;
  transition: transform 0.15s ease;
  flex-shrink: 0;
  margin-left: 0.5rem;
}

.select-arrow--open {
  transform: rotate(180deg);
  color: #3b82f6;
}

/* Desktop dropdown */
.select-dropdown--desktop {
  background: var(--select-dropdown-bg);
  border: 1px solid var(--select-dropdown-border);
  border-radius: var(--input-radius);
  box-shadow: var(--select-dropdown-shadow);
  max-height: 240px;
  overflow-y: auto;
  padding: 0.25rem 0;
  -webkit-overflow-scrolling: touch;
}

.select-option {
  padding: var(--select-option-padding-y) var(--select-option-padding-x);
  font-size: var(--input-font-size);
  font-family: inherit;
  color: var(--input-text);
  cursor: pointer;
  transition: background 0.1s ease;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-tap-highlight-color: transparent;
}

.select-option-content {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  overflow: hidden;
}

.select-option-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-option-desc {
  font-size: 0.75rem;
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-option:active:not(.select-option--disabled) {
  background: var(--select-option-hover-bg);
}

.select-option--active {
  background: var(--select-option-active-bg);
  color: #2563eb;
  font-weight: 500;
}

.select-option--disabled {
  color: #cbd5e1;
  pointer-events: none;
}

/* Mobile modal */
.select-modal {
  position: fixed;
  inset: 0;
  z-index: 999999;
}

.select-modal__overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.select-modal__panel {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: #ffffff;
  border-radius: 1rem 1rem 0 0;
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.12);
  max-height: 70vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.select-modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid #e5e7eb;
  flex-shrink: 0;
}

.select-modal__cancel {
  font-size: 0.875rem;
  color: #6b7280;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  font-family: inherit;
}

.select-modal__title {
  font-size: 0.9375rem;
  font-weight: 600;
  color: #1e293b;
}

.select-modal__confirm {
  font-size: 0.875rem;
  color: #2563eb;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  font-weight: 500;
  font-family: inherit;
}

.select-modal__options {
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  padding: 0.5rem 0;
  flex: 1;
}

.select-modal__option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1.25rem;
  cursor: pointer;
  transition: background 0.1s ease;
  -webkit-tap-highlight-color: transparent;
}

.select-modal__option-content {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  flex: 1;
  overflow: hidden;
  margin-right: 0.75rem;
}

.select-modal__option-desc {
  font-size: 0.8125rem;
  color: #94a3b8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-modal__option:active:not(.select-modal__option--disabled) {
  background: #f3f4f6;
}

.select-modal__option--active {
  background: #eff6ff;
}

.select-modal__option--disabled {
  color: #cbd5e1;
  pointer-events: none;
}

.select-modal__option-label {
  font-size: 0.9375rem;
  color: #1e293b;
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-modal__option--active .select-modal__option-label {
  color: #2563eb;
  font-weight: 500;
}

.select-modal__option-icon {
  width: 20px;
  height: 20px;
  color: #2563eb;
  flex-shrink: 0;
  margin-left: 0.75rem;
}

/* Transitions */
.select-dropdown-enter-active {
  transition: all 0.15s ease;
}

.select-dropdown-leave-active {
  transition: all 0.1s ease;
}

.select-dropdown-enter-from,
.select-dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.select-modal-enter-active .select-modal__panel {
  transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.select-modal-leave-active .select-modal__panel {
  transition: transform 0.2s ease;
}

.select-modal-enter-from .select-modal__panel,
.select-modal-leave-to .select-modal__panel {
  transform: translateY(100%);
}

.select-modal-enter-active .select-modal__overlay {
  transition: opacity 0.25s ease;
}

.select-modal-leave-active .select-modal__overlay {
  transition: opacity 0.2s ease;
}

.select-modal-enter-from .select-modal__overlay,
.select-modal-leave-to .select-modal__overlay {
  opacity: 0;
}
</style>
