import { ref } from 'vue'
import type { RegistrationRequest } from '~/types'

export function useRegistration() {
  const config = useRuntimeConfig()
  const submitting = ref(false)
  const successMessage = ref<string | null>(null)
  const errorMessage = ref<string | null>(null)
  const lastRegistration = ref<any>(null)

  async function submitRegistration(
    data: RegistrationRequest,
    photoFile?: File
  ) {
    submitting.value = true
    successMessage.value = null
    errorMessage.value = null

    try {
      if (photoFile) {
        const formData = new FormData()
        formData.append('name', data.name)
        formData.append('gender', data.gender)
        formData.append('birth_date', data.birth_date)
        formData.append('grade', data.grade)
        formData.append('class_id', String(data.class_id))
        formData.append('parent_name', data.parent_name)
        formData.append('parent_phone', data.parent_phone)
        formData.append('address', data.address)
        formData.append('id_number', data.id_number)
        formData.append('photo', photoFile)

        const response = await $fetch<any>(
          `${config.public.apiBase}/registrations`,
          {
            method: 'POST',
            body: formData,
          }
        )
        lastRegistration.value = response.data
        successMessage.value = response.data?.message || '报名成功！'
      } else {
        const params = new URLSearchParams()
        params.append('name', data.name)
        params.append('gender', data.gender)
        params.append('birth_date', data.birth_date)
        params.append('grade', data.grade)
        params.append('class_id', String(data.class_id))
        params.append('parent_name', data.parent_name)
        params.append('parent_phone', data.parent_phone)
        params.append('address', data.address)
        params.append('id_number', data.id_number)

        const response = await $fetch<any>(
          `${config.public.apiBase}/registrations`,
          {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: params,
          }
        )
        lastRegistration.value = response.data
        successMessage.value = response.data?.message || '报名成功！'
      }
    } catch (e: any) {
      if (e.data?.message) {
        errorMessage.value = e.data.message
      } else if (e.message) {
        errorMessage.value = e.message
      } else {
        errorMessage.value = '报名失败，请重试'
      }
    } finally {
      submitting.value = false
    }
  }

  return {
    submitting,
    successMessage,
    errorMessage,
    lastRegistration,
    submitRegistration,
  }
}
