'use client'

import { useQuery } from '@tanstack/react-query'
import Link from 'next/link'
import { useAuthStore } from '@/lib/store/auth'
import api from '@/lib/api'
import { Brain, BookOpen, PenLine, MessageSquare, ArrowRight } from 'lucide-react'

export default function DashboardPage() {
  const user = useAuthStore((s) => s.user)

  const { data: vocabProgress } = useQuery({
    queryKey: ['vocab-progress'],
    queryFn: () => api.get('/vocabulary/progress').then((r) => r.data),
  })

  const { data: grammarProgress } = useQuery({
    queryKey: ['grammar-progress'],
    queryFn: () => api.get('/grammar/progress').then((r) => r.data),
  })

  const { data: writingProgress } = useQuery({
    queryKey: ['writing-progress'],
    queryFn: () => api.get('/writing/progress').then((r) => r.data),
  })

  return (
    <div className="p-4 sm:p-8 max-w-5xl mx-auto">
      {/* Header */}
      <div className="mb-8 sm:mb-10">
        <h1 className="text-3xl sm:text-4xl font-black text-white tracking-tight">
          Good {getGreeting()},{' '}
          <span className="gradient-text">{user?.name?.split(' ')[0]}</span>
        </h1>
        <p className="text-white/45 mt-2 font-medium text-sm sm:text-base">Keep up the momentum — your next review is waiting.</p>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 sm:gap-4 mb-8 sm:mb-10">
        <StatCard label="Vocabulary Learned" value={vocabProgress?.learned ?? '—'} total={vocabProgress?.total} />
        <StatCard label="Due for Review" value={vocabProgress?.due_today ?? '—'} highlight />
        <StatCard label="Grammar Exercises" value={grammarProgress?.completed ?? '—'} total={grammarProgress?.total_exercises} />
      </div>

      {/* Module cards */}
      <h2 className="text-xs font-black text-white mb-4 uppercase tracking-widest">Your Modules</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
        {modules.map((mod) => (
          <Link
            key={mod.href}
            href={mod.href}
            className="group p-5 sm:p-6 glass glass-hover rounded-2xl transition-all duration-300 hover:-translate-y-1"
          >
            <div className="flex items-start justify-between mb-3">
              <div className={`p-2.5 rounded-xl ${mod.iconBg}`}>
                <mod.icon className={`w-5 h-5 ${mod.iconColor}`} />
              </div>
              <ArrowRight className="w-4 h-4 text-white/20 group-hover:text-violet-400 group-hover:translate-x-1 transition-all" />
            </div>
            <h3 className="text-white font-black mb-1 tracking-tight">{mod.title}</h3>
            <p className="text-white/45 text-sm">{mod.description}</p>
          </Link>
        ))}
      </div>
    </div>
  )
}

function StatCard({ label, value, total, highlight }: {
  label: string
  value: number | string
  total?: number
  highlight?: boolean
}) {
  return (
    <div className={`p-5 rounded-2xl glass ${highlight ? 'border-violet-500/40' : ''}`}>
      <p className="text-white/45 text-xs font-semibold uppercase tracking-wider mb-2">{label}</p>
      <p className={`text-3xl font-black tracking-tight ${highlight ? 'gradient-text' : 'text-white'}`}>
        {value}
        {total !== undefined && <span className="text-lg text-white/30 font-normal">/{total}</span>}
      </p>
    </div>
  )
}

function getGreeting() {
  const h = new Date().getHours()
  if (h < 12) return 'morning'
  if (h < 17) return 'afternoon'
  return 'evening'
}

const modules = [
  {
    href: '/dashboard/vocabulary',
    title: 'Vocabulary Builder',
    description: 'Review IT business terms with spaced repetition. 24 words across 3 categories.',
    icon: Brain,
    iconBg: 'bg-violet-500/15',
    iconColor: 'text-violet-400',
  },
  {
    href: '/dashboard/grammar',
    title: 'Grammar',
    description: 'Master passive voice, conditionals, and modal verbs in a business IT context.',
    icon: BookOpen,
    iconBg: 'bg-cyan-500/15',
    iconColor: 'text-cyan-400',
  },
  {
    href: '/dashboard/writing',
    title: 'Business Writing',
    description: 'Write emails, proposals, and reports with AI feedback on tone and grammar.',
    icon: PenLine,
    iconBg: 'bg-amber-500/15',
    iconColor: 'text-amber-400',
  },
  {
    href: '/dashboard/roleplay',
    title: 'Role-play Scenarios',
    description: 'Practice real conversations: standups, client calls, code reviews with AI.',
    icon: MessageSquare,
    iconBg: 'bg-pink-500/15',
    iconColor: 'text-pink-400',
  },
]
