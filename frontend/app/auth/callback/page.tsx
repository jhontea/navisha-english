'use client'

import { Suspense, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { useAuthStore } from '@/lib/store/auth'
import toast from 'react-hot-toast'

function CallbackHandler() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const fetchMe = useAuthStore((s) => s.fetchMe)

  useEffect(() => {
    const accessToken = searchParams.get('access_token')
    const refreshToken = searchParams.get('refresh_token')
    const error = searchParams.get('error')

    if (error) {
      const messages: Record<string, string> = {
        invalid_state: 'Authentication failed. Please try again.',
        missing_code: 'Google sign-in was cancelled.',
        token_exchange_failed: 'Failed to authenticate with Google.',
        invalid_id_token: 'Invalid Google token. Please try again.',
        db_error: 'Server error. Please try again.',
      }
      toast.error(messages[error] ?? 'Sign-in failed. Please try again.')
      router.replace('/login')
      return
    }

    if (!accessToken || !refreshToken) {
      toast.error('Missing authentication tokens.')
      router.replace('/login')
      return
    }

    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('refresh_token', refreshToken)

    fetchMe().then(() => {
      toast.success('Signed in successfully!')
      router.replace('/dashboard')
    })
  }, [searchParams, router, fetchMe])

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center">
      <div className="text-center">
        <div className="w-10 h-10 border-2 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p className="text-slate-400 text-sm">Signing you in...</p>
      </div>
    </div>
  )
}

export default function AuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-slate-950 flex items-center justify-center">
        <div className="text-center">
          <div className="w-10 h-10 border-2 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-slate-400 text-sm">Loading...</p>
        </div>
      </div>
    }>
      <CallbackHandler />
    </Suspense>
  )
}
