'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { ThemeProvider, useTheme } from 'next-themes'
import { useState } from 'react'

function ThemedToaster() {
  const { resolvedTheme } = useTheme()
  const isDark = resolvedTheme !== 'light'
  return (
    <Toaster
      position="top-right"
      toastOptions={{
        duration: 3000,
        style: {
          background: isDark ? '#1e1b4b' : '#ffffff',
          color: isDark ? '#f0f4ff' : '#0d0f1f',
          fontSize: '14px',
          border: isDark ? '1px solid rgba(139,92,246,0.25)' : '1px solid rgba(0,0,0,0.10)',
        },
      }}
    />
  )
}

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 1000 * 60 * 5,
        retry: 1,
      },
    },
  }))

  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="dark"
      enableSystem={false}
      disableTransitionOnChange={false}
    >
      <QueryClientProvider client={queryClient}>
        {children}
        <ThemedToaster />
      </QueryClientProvider>
    </ThemeProvider>
  )
}
