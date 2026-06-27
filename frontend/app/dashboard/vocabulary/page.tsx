'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import { cn, getDifficultyColor } from '@/lib/utils'
import { RotateCcw, CheckCircle, XCircle } from 'lucide-react'

type Vocab = {
  id: string
  word: string
  definition: string
  category: string
  example_sentence: string
  difficulty: string
  next_review?: string
}

type Mode = 'list' | 'review'

export default function VocabularyPage() {
  const [mode, setMode] = useState<Mode>('list')
  const [category, setCategory] = useState('')

  return (
    <div className="p-4 sm:p-8 max-w-4xl mx-auto">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6 sm:mb-8">
        <div>
          <h1 className="text-2xl sm:text-3xl font-black text-white tracking-tight">Vocabulary Builder</h1>
          <p className="text-white/45 text-sm mt-1 font-medium">IT Business English — Spaced Repetition (SM-2)</p>
        </div>
        <div className="flex gap-2 shrink-0">
          <button
            onClick={() => setMode('list')}
            className={cn('px-3 sm:px-4 py-2 rounded-xl text-sm font-bold transition-all',
              mode === 'list' ? 'btn-vibrant text-white' : 'glass text-white/50 hover:text-white')}
          >
            Browse
          </button>
          <button
            onClick={() => setMode('review')}
            className={cn('px-3 sm:px-4 py-2 rounded-xl text-sm font-bold transition-all',
              mode === 'review' ? 'btn-vibrant text-white' : 'glass text-white/50 hover:text-white')}
          >
            Review Session
          </button>
        </div>
      </div>

      {mode === 'list' ? (
        <VocabList category={category} setCategory={setCategory} />
      ) : (
        <ReviewSession />
      )}
    </div>
  )
}

function VocabList({ category, setCategory }: { category: string; setCategory: (c: string) => void }) {
  const { data, isLoading } = useQuery({
    queryKey: ['vocabulary', category],
    queryFn: () => api.get('/vocabulary', { params: { category } }).then((r) => r.data.data),
  })

  const categories = ['', 'Project Management', 'Technical Communication', 'Meeting & Email Phrases']

  if (isLoading) return <LoadingState />

  return (
    <>
      <div className="flex gap-2 mb-6 flex-wrap">
        {categories.map((c) => (
          <button
            key={c || 'all'}
            onClick={() => setCategory(c)}
            className={cn('px-3 py-1.5 rounded-lg text-sm font-semibold transition-all',
              category === c ? 'btn-vibrant text-white' : 'glass text-white/45 hover:text-white')}
          >
            {c || 'All Categories'}
          </button>
        ))}
      </div>

      <div className="space-y-3">
        {data?.map((v: Vocab) => (
          <div key={v.id} className="p-5 glass glass-hover rounded-2xl transition-all">
            <div className="flex items-start justify-between mb-2">
              <h3 className="text-white font-black text-lg tracking-tight">{v.word}</h3>
              <span className={cn('text-xs px-2 py-0.5 rounded-full font-semibold', getDifficultyColor(v.difficulty))}>
                {v.difficulty}
              </span>
            </div>
            <p className="text-white/70 text-sm mb-2">{v.definition}</p>
            <p className="text-white/35 text-sm italic">"{v.example_sentence}"</p>
          </div>
        ))}
      </div>
    </>
  )
}

function ReviewSession() {
  const queryClient = useQueryClient()
  const [cardIndex, setCardIndex] = useState(0)
  const [flipped, setFlipped] = useState(false)
  const [done, setDone] = useState(false)

  const { data: items, isLoading } = useQuery({
    queryKey: ['vocab-review'],
    queryFn: () => api.get('/vocabulary/review').then((r) => r.data.data),
  })

  const reviewMutation = useMutation({
    mutationFn: ({ id, quality }: { id: string; quality: number }) =>
      api.post(`/vocabulary/${id}/review`, { quality }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['vocab-progress'] })
    },
  })

  if (isLoading) return <LoadingState />
  if (!items?.length) return (
    <div className="text-center py-20">
      <CheckCircle className="w-12 h-12 text-green-400 mx-auto mb-4" />
      <p className="text-white font-black text-lg">All caught up!</p>
      <p className="text-white/45 text-sm mt-1">No vocabulary due for review right now.</p>
    </div>
  )

  if (done) return (
    <div className="text-center py-20">
      <CheckCircle className="w-12 h-12 text-violet-400 mx-auto mb-4" />
      <p className="text-white font-black text-lg">Session complete!</p>
      <p className="text-white/45 text-sm mt-1">You reviewed {items.length} words.</p>
      <button
        onClick={() => { setCardIndex(0); setFlipped(false); setDone(false); queryClient.invalidateQueries({ queryKey: ['vocab-review'] }) }}
        className="mt-6 px-6 py-2.5 btn-vibrant text-white rounded-xl text-sm font-bold transition-all"
      >
        Start New Session
      </button>
    </div>
  )

  const current: Vocab = items[cardIndex]

  const handleRate = async (quality: number) => {
    await reviewMutation.mutateAsync({ id: current.id, quality })
    toast.success(quality >= 3 ? 'Good job!' : 'Keep practicing!')
    if (cardIndex + 1 >= items.length) {
      setDone(true)
    } else {
      setCardIndex((i) => i + 1)
      setFlipped(false)
    }
  }

  return (
    <div className="max-w-lg mx-auto">
      <div className="text-center text-white/35 text-sm font-semibold mb-6 tracking-widest uppercase">
        Card {cardIndex + 1} of {items.length}
      </div>

      {/* Flashcard */}
      <div
        onClick={() => setFlipped(!flipped)}
        className="cursor-pointer p-8 glass glass-hover rounded-3xl min-h-48 flex flex-col items-center justify-center text-center mb-6 transition-all duration-300 hover:-translate-y-1"
      >
        {!flipped ? (
          <>
            <p className="text-4xl font-black text-white mb-3 tracking-tight">{current.word}</p>
            <p className={cn('text-xs px-3 py-1 rounded-full font-bold', getDifficultyColor(current.difficulty))}>
              {current.category}
            </p>
            <p className="text-white/30 text-sm mt-5">Tap to reveal</p>
          </>
        ) : (
          <>
            <p className="text-white text-base mb-3 font-medium leading-relaxed">{current.definition}</p>
            <p className="text-white/40 text-sm italic">"{current.example_sentence}"</p>
          </>
        )}
      </div>

      {/* Rating buttons — only show after flip */}
      {flipped && (
        <div className="space-y-3">
          <p className="text-center text-white/35 text-sm font-semibold mb-4">How well did you know this?</p>
            <div className="grid grid-cols-2 gap-3">
            <button
              onClick={() => handleRate(1)}
              className="flex items-center justify-center gap-2 py-3 panel-danger rounded-xl text-sm font-bold transition-all hover:opacity-90"
            >
              <XCircle className="w-4 h-4" /> Didn't know
            </button>
            <button
              onClick={() => handleRate(3)}
              className="flex items-center justify-center gap-2 py-3 panel-warning rounded-xl text-sm font-bold transition-all hover:opacity-90"
            >
              <RotateCcw className="w-4 h-4" /> Hard
            </button>
            <button
              onClick={() => handleRate(4)}
              className="flex items-center justify-center gap-2 py-3 panel-success rounded-xl text-sm font-bold transition-all hover:opacity-90"
            >
              <CheckCircle className="w-4 h-4" /> Good
            </button>
            <button
              onClick={() => handleRate(5)}
              className="flex items-center justify-center gap-2 py-3 bg-violet-500/10 border border-violet-500/30 text-violet-700 dark:text-violet-400 rounded-xl text-sm font-bold transition-all hover:opacity-90"
            >
              <CheckCircle className="w-4 h-4" /> Easy
            </button>
          </div>
        </div>
      )}
    </div>
  )
}

function LoadingState() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map((i) => (
        <div key={i} className="h-24 glass rounded-2xl animate-pulse" />
      ))}
    </div>
  )
}
