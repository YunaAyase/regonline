export interface ClassItem {
  id: number
  name: string
  description: string
  max_students: number
  min_age: number
  max_age: number
  enabled: boolean
  current_count: number
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface RegistrationRequest {
  name: string
  gender: string
  birth_date: string
  grade: string
  class_id: number
  parent_name: string
  parent_phone: string
  address: string
  id_number: string
}

export interface SiteSettings {
  site_name: string
  site_description: string
  icp_record: string
  copyright: string
}

export interface UpdateAccountRequest {
  username: string
  old_password: string
  new_password: string
}

export interface ServerInfo {
  go_version: string
  os: string
  arch: string
  cpu_cores: number
  goroutines: number
  memory_alloc: string
  uptime_seconds: number
  db_size: string
}

export interface RegistrationRecord {
  id: number
  name: string
  gender: string
  birth_date: string
  grade: string
  class_id: number
  parent_name: string
  parent_phone: string
  address: string
  id_number: string
  photo_path: string | null
  registration_time: string
  class?: {
    id: number
    name: string
  }
}

export interface OCRNumberCandidate {
  value: string
  label: string
  is_id_number: boolean
  length: number
}

export interface OCRRecognizeResult {
  id_number: string
  candidates: OCRNumberCandidate[]
}