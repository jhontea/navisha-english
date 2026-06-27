'use client'

import { useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import Link from 'next/link'
import { useAuthStore } from '@/lib/store/auth'
import toast from 'react-hot-toast'
import { BookOpen, PenLine, Brain, MessageSquare, LogOut, TrendingUp, Menu, X, Languages } from 'lucide-react'
import { ThemeToggle } from '@/app/components/ThemeToggle'

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, logout, _hasHydrated } = useAuthStore()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  useEffect(() => {
    // Tunggu sampai Zustand selesai baca dari localStorage
    if (!_hasHydrated) return
    if (!user) router.push('/login')
  }, [user, _hasHydrated, router])

  // Close sidebar on route change (mobile)
  useEffect(() => {
    setSidebarOpen(false)
  }, [pathname])

  const handleLogout = async () => {
    await logout()
    toast.success('Signed out')
    router.push('/')
  }

  // Tampilkan loading spinner saat hydration belum selesai
  if (!_hasHydrated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="w-8 h-8 rounded-full border-2 border-violet-500/30 border-t-violet-500 animate-spin" />
      </div>
    )
  }

  if (!user) return null

  const SidebarContent = () => (
    <>
      <div className="p-5 border-b border-white/10 flex items-center justify-between">
        <div>
          <span className="text-lg font-black text-white tracking-tight">
            Navisha<span className="gradient-text"> English</span>
          </span>
          <p className="text-xs text-white/35 mt-0.5 font-medium">Business English for IT</p>
        </div>
        {/* Close button — mobile only */}
        <button
          onClick={() => setSidebarOpen(false)}
          className="lg:hidden p-1.5 text-white/45 hover:text-white transition-colors"
          aria-label="Close menu"
        >
          <X className="w-5 h-5" />
        </button>
      </div>

      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        {navItems.map((item) => {
          const isActive = item.href === '/dashboard'
            ? pathname === '/dashboard'
            : pathname.startsWith(item.href)
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all group ${
                isActive
                  ? 'bg-violet-500/20 border border-violet-500/30 text-white'
                  : 'text-white/45 hover:text-white hover:bg-white/8 border border-transparent'
              }`}
            >
              <item.icon className={`w-5 h-5 shrink-0 transition-colors ${
                isActive ? 'text-violet-400' : 'group-hover:text-violet-400'
              }`} />
              <span className={`text-sm font-semibold ${isActive ? 'font-black' : ''}`}>{item.label}</span>
              {isActive && (
                <span className="ml-auto w-1.5 h-1.5 rounded-full bg-violet-400 shrink-0" />
              )}
            </Link>
          )
        })}
      </nav>

      <div className="p-4 border-t border-white/10">
        <div className="flex items-center gap-3 px-3 py-2 mb-2">
          <div className="w-9 h-9 rounded-full btn-vibrant flex items-center justify-center text-white text-sm font-black shrink-0">
            {user.name.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-bold text-white truncate">{user.name}</p>
            <p className="text-xs text-white/35 font-medium">Level {user.level}</p>
          </div>
          <ThemeToggle />
        </div>
        <button
          onClick={handleLogout}
          className="flex items-center gap-2 w-full px-3 py-2 text-white/35 hover:text-red-400 hover:bg-red-500/10 rounded-xl transition-all text-sm font-semibold"
        >
          <LogOut className="w-4 h-4" />
          Sign out
        </button>
      </div>
    </>
  )

  return (
    <div className="min-h-screen flex">
      {/* ── Desktop sidebar (fixed) ── */}
      <aside className="hidden lg:flex w-64 glass border-r border-white/10 flex-col fixed h-full z-20">
        <SidebarContent />
      </aside>

      {/* ── Mobile overlay backdrop ── */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/60 z-30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* ── Mobile sidebar (slide-in drawer) ── */}
      <aside
        className={`
          fixed top-0 left-0 h-full w-72 glass flex flex-col z-40 lg:hidden
          transition-transform duration-300 ease-in-out
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        `}
      >
        <SidebarContent />
      </aside>

      {/* ── Main content ── */}
      <main className="flex-1 overflow-auto lg:ml-64 flex flex-col min-h-screen">
        {/* Mobile top bar */}
        <div className="lg:hidden flex items-center gap-3 px-4 py-3 glass border-b border-white/10 sticky top-0 z-20">
          <button
            onClick={() => setSidebarOpen(true)}
            className="p-2 text-white/60 hover:text-white transition-colors rounded-xl hover:bg-white/8"
            aria-label="Open menu"
          >
            <Menu className="w-5 h-5" />
          </button>
          <span className="flex-1 text-white font-black tracking-tight text-base">
            Navisha<span className="gradient-text"> English</span>
          </span>
          <ThemeToggle />
        </div>

        <div className="flex-1">
          {children}
        </div>
      </main>
    </div>
  )
}

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: TrendingUp },
  { href: '/dashboard/vocabulary', label: 'Vocabulary', icon: Brain },
  { href: '/dashboard/grammar', label: 'Grammar', icon: BookOpen },
  { href: '/dashboard/writing', label: 'Business Writing', icon: PenLine },
  { href: '/dashboard/roleplay', label: 'Role-play', icon: MessageSquare },
  { href: '/dashboard/word-challenge', label: 'Word Challenge', icon: Languages },
]
