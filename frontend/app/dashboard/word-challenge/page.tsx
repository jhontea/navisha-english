'use client'

import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import { cn } from '@/lib/utils'
import {
  Sparkles, CheckCircle, XCircle, RefreshCw, BookOpen,
  ChevronRight, Lightbulb, AlertTriangle, ArrowRight, Briefcase
} from 'lucide-react'

type SentenceData = {
  challenge_id: string
  indonesian_sentence: string
  topic: string
  correct_answer: string
}

type CheckResult = {
  is_correct: boolean
  correct_answer: string
  user_answer: string
  explanation: string
  corrections: string
}

type HistoryItem = {
  indonesian_sentence: string
  correct_answer: string
  user_answer: string
  is_correct: boolean
  explanation: string
  corrections: string
  attempted_at: string
}

type PageState = 'idle' | 'challenge' | 'result'

// Parses inline **bold** and `code` markers into React elements
function RichText({ text, className }: { text: string; className?: string }) {
  const parts = text.split(/(\*\*[^*]+\*\*|`[^`]+`)/)
  return (
    <span className={className}>
      {parts.map((part, i) => {
        if (part.startsWith('**') && part.endsWith('**')) {
          return <strong key={i} className="text-white font-semibold">{part.slice(2, -2)}</strong>
        }
        if (part.startsWith('`') && part.endsWith('`')) {
          return (
            <code key={i} className="text-cyan-300 bg-cyan-500/10 px-1.5 py-0.5 rounded font-mono text-xs mx-0.5">
              {part.slice(1, -1)}
            </code>
          )
        }
        return <span key={i}>{part}</span>
      })}
    </span>
  )
}

// Splits text by newlines or "- " bullet patterns and renders each as a block
function RichBlock({ text, className }: { text: string; className?: string }) {
  const lines = text
    .split(/\n|(?=- )/)
    .map(l => l.replace(/^-\s*/, '').trim())
    .filter(Boolean)

  if (lines.length <= 1) {
    return <RichText text={text} className={cn('text-sm leading-relaxed', className)} />
  }

  return (
    <div className={cn('space-y-2', className)}>
      {lines.map((line, i) => (
        <div key={i} className="flex items-start gap-2.5">
          <span className="mt-2 w-1.5 h-1.5 rounded-full bg-current opacity-50 flex-shrink-0" />
          <RichText text={line} className="text-sm leading-relaxed" />
        </div>
      ))}
    </div>
  )
}

export default function WordChallengePage() {
  const [pageState, setPageState] = useState<PageState>('idle')
  const [currentSentence, setCurrentSentence] = useState<SentenceData | null>(null)
  const [userAnswer, setUserAnswer] = useState('')
  const [result, setResult] = useState<CheckResult | null>(null)
  const [showHistory, setShowHistory] = useState(false)

  const { data: history, refetch: refetchHistory } = useQuery({
    queryKey: ['word-challenge-history'],
    queryFn: () => api.get('/word-challenge/history').then((r) => r.data.data as HistoryItem[]),
    enabled: showHistory,
  })

  const generateMutation = useMutation({
    mutationFn: () => api.get('/word-challenge/generate').then((r) => r.data.data as SentenceData),
    onSuccess: (data) => {
      setCurrentSentence(data)
      setUserAnswer('')
      setResult(null)
      setPageState('challenge')
    },
    onError: () => toast.error('Gagal mendapatkan kalimat. Coba lagi.'),
  })

  const checkMutation = useMutation({
    mutationFn: (answer: string) =>
      api.post('/word-challenge/check', {
        challenge_id: currentSentence!.challenge_id,
        indonesian_sentence: currentSentence!.indonesian_sentence,
        correct_answer: currentSentence!.correct_answer,
        user_answer: answer,
      }).then((r) => r.data as CheckResult),
    onSuccess: (data) => {
      setResult(data)
      setPageState('result')
      refetchHistory()
    },
    onError: () => toast.error('Gagal memeriksa jawaban. Coba lagi.'),
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!userAnswer.trim()) {
      toast.error('Masukkan terjemahan terlebih dahulu.')
      return
    }
    checkMutation.mutate(userAnswer.trim())
  }

  return (
    <div className="p-8 max-w-2xl mx-auto">
      {/* Header */}
      <div className="mb-8 flex items-start justify-between">
        <div>
          <h1 className="text-3xl font-black text-white tracking-tight">Word Challenge</h1>
          <p className="text-white/45 text-sm mt-1 font-medium">
            Terjemahkan kalimat bisnis Indonesia ke Business English
          </p>
        </div>
        <button
          onClick={() => setShowHistory((v) => !v)}
          className="flex items-center gap-2 px-4 py-2 glass glass-hover rounded-xl text-white/60 hover:text-white text-sm font-semibold transition-all"
        >
          <BookOpen className="w-4 h-4" />
          Riwayat
        </button>
      </div>

      {showHistory ? (
        <HistoryPanel history={history ?? []} onClose={() => setShowHistory(false)} />
      ) : (
        <>
          {pageState === 'idle' && (
            <IdleState onStart={() => generateMutation.mutate()} isLoading={generateMutation.isPending} />
          )}
          {pageState === 'challenge' && currentSentence && (
            <ChallengeState
              sentence={currentSentence}
              userAnswer={userAnswer}
              setUserAnswer={setUserAnswer}
              onSubmit={handleSubmit}
              isLoading={checkMutation.isPending}
            />
          )}
          {pageState === 'result' && result && currentSentence && (
            <ResultState
              result={result}
              indonesianSentence={currentSentence.indonesian_sentence}
              topic={currentSentence.topic}
              onNext={() => generateMutation.mutate()}
              isLoading={generateMutation.isPending}
            />
          )}
        </>
      )}
    </div>
  )
}

function IdleState({ onStart, isLoading }: { onStart: () => void; isLoading: boolean }) {
  const features = [
    { icon: Briefcase, text: 'Kalimat dari konteks IT & software engineering' },
    { icon: AlertTriangle, text: 'Koreksi spesifik jika ada kesalahan' },
    { icon: Lightbulb, text: 'Penjelasan kenapa phrasing tertentu lebih tepat' },
  ]

  return (
    <div className="flex flex-col items-center gap-8 py-8">
      <div className="p-6 glass rounded-3xl">
        <Sparkles className="w-12 h-12 text-violet-400" />
      </div>

      <div className="text-center space-y-2">
        <h2 className="text-xl font-black text-white">Siap untuk berlatih?</h2>
        <p className="text-white/45 text-sm font-medium max-w-xs">
          AI memberikan kalimat bisnis Indonesia. Kamu terjemahkan ke Business English yang profesional.
        </p>
      </div>

      {/* Feature pills */}
      <div className="w-full space-y-2">
        {features.map(({ icon: Icon, text }, i) => (
          <div key={i} className="flex items-center gap-3 px-4 py-3 glass rounded-2xl">
            <Icon className="w-4 h-4 text-violet-400 flex-shrink-0" />
            <p className="text-white/60 text-sm font-medium">{text}</p>
          </div>
        ))}
      </div>

      <button
        onClick={onStart}
        disabled={isLoading}
        className="w-full btn-vibrant py-3.5 rounded-2xl text-white font-black tracking-tight flex items-center justify-center gap-2 disabled:opacity-50"
      >
        {isLoading ? <RefreshCw className="w-5 h-5 animate-spin" /> : <Sparkles className="w-5 h-5" />}
        Mulai Challenge
      </button>
    </div>
  )
}

function ChallengeState({
  sentence, userAnswer, setUserAnswer, onSubmit, isLoading,
}: {
  sentence: SentenceData
  userAnswer: string
  setUserAnswer: (v: string) => void
  onSubmit: (e: React.FormEvent) => void
  isLoading: boolean
}) {
  return (
    <div className="space-y-5">
      {/* Topic badge */}
      <div className="flex items-center gap-2">
        <span className="text-xs px-3 py-1 rounded-full bg-violet-500/15 text-violet-300 font-semibold border border-violet-500/20">
          {sentence.topic}
        </span>
      </div>

      {/* Sentence card */}
      <div className="relative p-8 glass rounded-3xl space-y-2 overflow-hidden">
        {/* Decorative gradient blob */}
        <div className="absolute -top-8 -right-8 w-32 h-32 bg-violet-600/10 rounded-full blur-2xl pointer-events-none" />
        <p className="text-white/35 text-xs uppercase tracking-widest font-semibold">
          Terjemahkan ke Business English
        </p>
        <p className="text-2xl font-bold text-white leading-relaxed">
          {sentence.indonesian_sentence}
        </p>
      </div>

      {/* Answer form */}
      <form onSubmit={onSubmit} className="space-y-3">
        <label className="block text-white/45 text-xs uppercase tracking-widest font-semibold">
          Terjemahan Business English kamu
        </label>
        <textarea
          value={userAnswer}
          onChange={(e) => setUserAnswer(e.target.value)}
          placeholder="Ketik terjemahan dalam Business English..."
          autoFocus
          rows={3}
          className="w-full px-4 py-3 glass rounded-2xl text-white placeholder-white/25 font-medium text-base focus:outline-none focus:ring-2 focus:ring-violet-500/50 transition-all resize-none"
        />
        <button
          type="submit"
          disabled={isLoading || !userAnswer.trim()}
          className="w-full btn-vibrant py-3.5 rounded-2xl text-white font-black tracking-tight flex items-center justify-center gap-2 disabled:opacity-50 transition-all"
        >
          {isLoading
            ? <RefreshCw className="w-5 h-5 animate-spin" />
            : <><span>Periksa Terjemahan</span><ChevronRight className="w-5 h-5" /></>
          }
        </button>
      </form>
    </div>
  )
}

function ResultState({
  result, indonesianSentence, topic, onNext, isLoading,
}: {
  result: CheckResult
  indonesianSentence: string
  topic: string
  onNext: () => void
  isLoading: boolean
}) {
  return (
    <div className="space-y-4">
      {/* Topic badge */}
      <div className="flex items-center gap-2">
        <span className="text-xs px-3 py-1 rounded-full bg-violet-500/15 text-violet-300 font-semibold border border-violet-500/20">
          {topic}
        </span>
      </div>

      {/* Verdict + original sentence */}
      <div className={cn(
        'rounded-3xl overflow-hidden border',
        result.is_correct ? 'border-emerald-500/25' : 'border-red-500/25'
      )}>
        {/* Verdict strip */}
        <div className={cn(
          'px-5 py-4 flex items-center gap-3',
          result.is_correct ? 'bg-emerald-500/10' : 'bg-red-500/10'
        )}>
          {result.is_correct
            ? <CheckCircle className="w-6 h-6 text-emerald-400 flex-shrink-0" />
            : <XCircle className="w-6 h-6 text-red-400 flex-shrink-0" />
          }
          <div>
            <p className={cn('font-black tracking-tight', result.is_correct ? 'text-emerald-300' : 'text-red-300')}>
              {result.is_correct ? 'Terjemahan tepat!' : 'Perlu diperbaiki'}
            </p>
            <p className="text-white/40 text-xs mt-0.5 italic">"{indonesianSentence}"</p>
          </div>
        </div>

        {/* Answer comparison rows */}
        <div className="divide-y divide-white/5">
          {/* Correct */}
          <div className="px-5 py-4 flex items-start gap-3">
            <CheckCircle className="w-4 h-4 text-emerald-400 mt-0.5 flex-shrink-0" />
            <div>
              <p className="text-emerald-400 text-xs font-bold uppercase tracking-widest mb-1">Terjemahan tepat</p>
              <p className="text-white text-sm font-semibold leading-relaxed">{result.correct_answer}</p>
            </div>
          </div>

          {/* User answer if wrong */}
          {!result.is_correct && (
            <div className="px-5 py-4 flex items-start gap-3">
              <XCircle className="w-4 h-4 text-red-400 mt-0.5 flex-shrink-0" />
              <div>
                <p className="text-red-400 text-xs font-bold uppercase tracking-widest mb-1">Jawaban kamu</p>
                <p className="text-white/50 text-sm leading-relaxed">{result.user_answer}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Corrections card */}
      {result.corrections && (
        <div className="rounded-2xl overflow-hidden border border-amber-500/20">
          <div className="px-4 py-3 bg-amber-500/10 border-b border-amber-500/15 flex items-center gap-2">
            <AlertTriangle className="w-4 h-4 text-amber-400" />
            <p className="text-amber-300 text-xs font-bold uppercase tracking-widest">Koreksi</p>
          </div>
          <div className="p-4 bg-amber-500/5">
            <RichBlock text={result.corrections} className="text-amber-200/85" />
          </div>
        </div>
      )}

      {/* Explanation card */}
      {result.explanation && (
        <div className="rounded-2xl overflow-hidden border border-violet-500/20">
          <div className="px-4 py-3 bg-violet-500/10 border-b border-violet-500/15 flex items-center gap-2">
            <Lightbulb className="w-4 h-4 text-violet-400" />
            <p className="text-violet-300 text-xs font-bold uppercase tracking-widest">Penjelasan</p>
          </div>
          <div className="p-4 bg-violet-500/5 space-y-4">
            <RichBlock text={result.explanation} className="text-white/80" />

            {/* Visual: "Correct answer" pinned at bottom as reference */}
            <div className="flex items-start gap-3 pt-2 border-t border-white/5">
              <ArrowRight className="w-3.5 h-3.5 text-violet-400 mt-0.5 flex-shrink-0" />
              <p className="text-violet-300/70 text-xs italic leading-relaxed">
                {result.correct_answer}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Next button */}
      <button
        onClick={onNext}
        disabled={isLoading}
        className="w-full btn-vibrant py-3.5 rounded-2xl text-white font-black tracking-tight flex items-center justify-center gap-2 disabled:opacity-50 transition-all"
      >
        {isLoading
          ? <RefreshCw className="w-5 h-5 animate-spin" />
          : <><span>Kalimat Berikutnya</span><Sparkles className="w-5 h-5" /></>
        }
      </button>
    </div>
  )
}

function HistoryPanel({ history, onClose }: { history: HistoryItem[]; onClose: () => void }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-white font-black text-lg">Riwayat Jawaban</h2>
        <button onClick={onClose} className="text-white/45 hover:text-white text-sm font-semibold transition-colors">
          Tutup
        </button>
      </div>

      {history.length === 0 ? (
        <div className="text-center py-12 text-white/35 text-sm font-medium">
          Belum ada riwayat. Mulai challenge dulu!
        </div>
      ) : (
        <div className="space-y-3">
          {history.map((item, i) => (
            <div key={i} className={cn(
              'rounded-2xl overflow-hidden border',
              item.is_correct ? 'border-emerald-500/15' : 'border-red-500/15'
            )}>
              {/* Header */}
              <div className={cn(
                'px-4 py-3 flex items-center gap-2',
                item.is_correct ? 'bg-emerald-500/8' : 'bg-red-500/8'
              )}>
                {item.is_correct
                  ? <CheckCircle className="w-3.5 h-3.5 text-emerald-400 flex-shrink-0" />
                  : <XCircle className="w-3.5 h-3.5 text-red-400 flex-shrink-0" />
                }
                <p className="text-white/60 text-xs italic truncate">"{item.indonesian_sentence}"</p>
              </div>

              {/* Body */}
              <div className="px-4 py-3 space-y-2 glass">
                <p className="text-emerald-300/80 text-xs font-semibold">{item.correct_answer}</p>
                {!item.is_correct && (
                  <p className="text-red-300/60 text-xs">Kamu: {item.user_answer}</p>
                )}
                {item.corrections && (
                  <p className="text-amber-300/60 text-xs leading-relaxed border-t border-white/5 pt-2">
                    {item.corrections}
                  </p>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
