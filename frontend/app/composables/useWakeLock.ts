/**
 * 移动端后台保活机制
 * - Wake Lock API: 防止屏幕休眠（iOS Safari 16.4+, Chrome Android 84+）
 * - NoSleep fallback: 隐藏视频元素持续播放（兼容旧浏览器）
 * - Page Visibility: 页面恢复可见时重新获取锁
 * - Heartbeat: 定时 ping 后端保持 HTTP 连接活跃
 */
export function useWakeLock() {
  let wakeLock: WakeLockSentinel | null = null
  let noSleepVideo: HTMLVideoElement | null = null
  let heartbeatTimer: ReturnType<typeof setInterval> | null = null
  let visibilityHandler: (() => void) | null = null

  const isSupported = ref(false)
  const isActive = ref(false)
  const lockType = ref<'wake-lock' | 'nosleep' | 'none'>('none')

  async function requestWakeLock() {
    try {
      if ('wakeLock' in navigator) {
        wakeLock = await navigator.wakeLock.request('screen')
        isActive.value = true
        lockType.value = 'wake-lock'
        isSupported.value = true

        wakeLock.addEventListener('release', () => {
          isActive.value = false
        })

        return true
      }
    } catch (err) {
      console.warn('Wake Lock API 不可用，尝试 NoSleep 方式:', err)
    }

    return false
  }

  function startNoSleep() {
    try {
      noSleepVideo = document.createElement('video')
      noSleepVideo.setAttribute('playsinline', '')
      noSleepVideo.setAttribute('muted', '')
      noSleepVideo.setAttribute('loop', '')
      noSleepVideo.style.position = 'fixed'
      noSleepVideo.style.top = '0'
      noSleepVideo.style.left = '0'
      noSleepVideo.style.width = '2px'
      noSleepVideo.style.height = '2px'
      noSleepVideo.style.opacity = '0.01'
      noSleepVideo.style.pointerEvents = 'none'
      noSleepVideo.style.zIndex = '-1'

      const canvas = document.createElement('canvas')
      canvas.width = 2
      canvas.height = 2
      const ctx = canvas.getContext('2d')
      if (ctx) {
        ctx.fillStyle = '#000000'
        ctx.fillRect(0, 0, 2, 2)
      }
      noSleepVideo.src = canvas.toDataURL('video/webm', 0.5) || ''

      noSleepVideo.play().then(() => {
        isActive.value = true
        lockType.value = 'nosleep'
        isSupported.value = true
      }).catch(() => {
        console.warn('NoSleep 视频播放失败')
      })

      document.body.appendChild(noSleepVideo)
      return true
    } catch (err) {
      console.warn('NoSleep 方式不可用:', err)
      return false
    }
  }

  function stopNoSleep() {
    if (noSleepVideo) {
      noSleepVideo.pause()
      noSleepVideo.remove()
      noSleepVideo = null
    }
  }

  function startHeartbeat(apiBase: string) {
    stopHeartbeat()
    heartbeatTimer = setInterval(() => {
      fetch(`${apiBase}/server-info`, {
        method: 'GET',
        credentials: 'include',
        cache: 'no-cache',
      }).catch(() => {})
    }, 30000)
  }

  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  async function reacquireLock() {
    if (document.visibilityState === 'visible' && !isActive.value) {
      const granted = await requestWakeLock()
      if (!granted && !noSleepVideo) {
        startNoSleep()
      }
    }
  }

  function setupVisibilityListener() {
    visibilityHandler = () => {
      if (document.visibilityState === 'visible') {
        reacquireLock()
      }
    }
    document.addEventListener('visibilitychange', visibilityHandler)
  }

  function removeVisibilityListener() {
    if (visibilityHandler) {
      document.removeEventListener('visibilitychange', visibilityHandler)
      visibilityHandler = null
    }
  }

  async function acquire(apiBase?: string) {
    if (isActive.value) return

    const granted = await requestWakeLock()
    if (!granted) {
      startNoSleep()
    }

    setupVisibilityListener()

    if (apiBase) {
      startHeartbeat(apiBase)
    }
  }

  function release() {
    if (wakeLock) {
      wakeLock.release().catch(() => {})
      wakeLock = null
    }
    stopNoSleep()
    stopHeartbeat()
    removeVisibilityListener()
    isActive.value = false
    lockType.value = 'none'
  }

  return {
    isSupported,
    isActive,
    lockType,
    acquire,
    release,
  }
}