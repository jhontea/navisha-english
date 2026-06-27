'use client'

import { useState } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import { cn, getDifficultyColor } from '@/lib/utils'
import { ArrowLeft, Sparkles, CheckCircle } from 'lucide-react'

type Exercise = {
  id: string
  title: string
  type: string
  context: string
  prompt: string
  template?: string
  key_phrases?: string[]
  difficulty: string
}

type Feedback = {
  score: number
  overall: string
  grammar?: string[]
  tone?: string
  vocabulary?: string[]
  improved_version?: string
  error?: string
}

export default function WritingPage() {
  const [selected, setSelected] = useState<Exercise | null>(null)

  return (
    <div className="p-4 sm:p-8 max-w-4xl mx-auto">
      <div className="mb-6 sm:mb-8">
        <h1 className="text-2xl sm:text-3xl font-black text-white tracking-tight">Business Writing</h1>
        <p className="text-white/45 text-sm mt-1 font-medium">Emails, proposals, and reports — with AI feedback</p>
      </div>
      {selected ? (
        <ExerciseView exercise={selected} onBack={() => setSelected(null)} />
      ) : (
        <ExerciseList onSelect={setSelected} />
      )}
    </div>
  )
}

function ExerciseList({ onSelect }: { onSelect: (e: Exercise) => void }) {
  const { data, isLoading } = useQuery({
    queryKey: ['writing-exercises'],
    queryFn: () => api.get('/writing/exercises').then((r) => r.data.data),
  })

  if (isLoading) return <LoadingState />

  const typeLabels: Record<string, string> = {
    email: 'Email',
    proposal: 'Proposal',
    report: 'Report',
  }

  return (
    <div className="space-y-3">
      {data?.map((e: Exercise) => (
        <button
          key={e.id}
          onClick={() => onSelect(e)}
          className="w-full text-left p-5 glass glass-hover rounded-2xl transition-all hover:-translate-y-0.5"
        >
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-white font-black tracking-tight">{e.title}</h3>
            <div className="flex items-center gap-2">
              <span className="text-xs px-2 py-0.5 rounded-full glass text-white/60 font-semibold">
                {typeLabels[e.type] ?? e.type}
              </span>
              <span className={cn('text-xs px-2 py-0.5 rounded-full font-bold', getDifficultyColor(e.difficulty))}>
                {e.difficulty}
              </span>
            </div>
          </div>
          <p className="text-white/45 text-sm line-clamp-2">{e.context}</p>
        </button>
      ))}
    </div>
  )
}

function ExerciseView({ exercise, onBack }: { exercise: Exercise; onBack: () => void }) {
  const [content, setContent] = useState(exercise.template ?? '')
  const [feedback, setFeedback] = useState<Feedback | null>(null)

  const submitMutation = useMutation({
    mutationFn: (text: string) =>
      api.post(`/writing/exercises/${exercise.id}/submit`, { content: text }),
    onSuccess: (res) => {
      setFeedback(res.data.feedback)
      toast.success('AI feedback received!')
    },
    onError: () => toast.error('Submission failed. Please try again.'),
  })

  return (
    <div>
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-white/45 hover:text-white text-sm mb-6 transition-colors font-semibold"
      >
        <ArrowLeft className="w-4 h-4" /> Back to exercises
      </button>

      <div className="mb-6 p-5 glass rounded-2xl">
        <h2 className="text-2xl font-black text-white mb-3 tracking-tight">{exercise.title}</h2>
        <p className="text-white/50 text-sm mb-4 leading-relaxed">{exercise.context}</p>
        <div className="p-4 bg-violet-500/10 border border-violet-500/20 rounded-xl">
          <p className="text-violet-300 text-sm font-bold">Task</p>
          <p className="text-violet-100 text-sm mt-1 leading-relaxed">{exercise.prompt}</p>
        </div>

        {exercise.key_phrases && exercise.key_phrases.length > 0 && (
          <div className="mt-4">
            <p className="text-white/30 text-xs mb-2 uppercase tracking-widest font-semibold">Key phrases to use</p>
            <div className="flex flex-wrap gap-2">
              {exercise.key_phrases.map((p) => (
                <span key={p} className="text-xs px-2.5 py-1 glass text-white/60 rounded-lg font-medium">
                  {p}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={14}
        placeholder="Write your response here..."
        className="w-full px-4 py-3 input-glass rounded-2xl resize-none text-sm font-mono"
      />

      <button
        onClick={() => submitMutation.mutate(content)}
        disabled={content.trim().length < 10 || submitMutation.isPending}
        className="mt-4 flex items-center gap-2 px-6 py-3 btn-vibrant disabled:opacity-40 disabled:cursor-not-allowed disabled:transform-none text-white font-black rounded-2xl transition-all"
      >
        <Sparkles className="w-4 h-4" />
        {submitMutation.isPending ? 'Getting AI feedback...' : 'Submit for AI Feedback'}
      </button>

      {feedback && <FeedbackPanel feedback={feedback} />}
    </div>
  )
}

function FeedbackPanel({ feedback }: { feedback: Feedback }) {
  const score = feedback.score ?? 0
  return (
    <div className="mt-6 space-y-4">
      {/* Score */}
      <div className={cn('p-5 rounded-2xl', score >= 80 ? 'panel-success' : score >= 50 ? 'panel-warning' : 'panel-danger')}>
        <div className="flex items-center justify-between mb-2">
          <p className="font-black tracking-tight">Overall Score</p>
          <p className="text-3xl font-black tracking-tight">{score}%</p>
        </div>
        <p className="text-sm opacity-80">{feedback.overall}</p>
      </div>

      {/* Tone */}
      {feedback.tone && (
        <div className="p-4 glass rounded-2xl">
          <p className="text-fg-subtle text-xs uppercase tracking-widest mb-1 font-semibold">Tone</p>
          <p className="text-fg text-sm">{feedback.tone}</p>
        </div>
      )}

      {/* Grammar issues */}
      {feedback.grammar && feedback.grammar.length > 0 && (
        <div className="p-4 glass rounded-2xl">
          <p className="text-fg-subtle text-xs uppercase tracking-widest mb-2 font-semibold">Grammar Issues</p>
          <ul className="space-y-1">
            {feedback.grammar.map((g, i) => (
              <li key={i} className="text-sm text-danger flex items-start gap-2">
                <span className="mt-0.5">•</span>{g}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Vocabulary suggestions */}
      {feedback.vocabulary && feedback.vocabulary.length > 0 && (
        <div className="p-4 glass rounded-2xl">
          <p className="text-fg-subtle text-xs uppercase tracking-widest mb-2 font-semibold">Vocabulary Suggestions</p>
          <ul className="space-y-1">
            {feedback.vocabulary.map((v, i) => (
              <li key={i} className="text-sm text-info flex items-start gap-2">
                <span className="mt-0.5">•</span>{v}
              </li>
            ))}
          </ul>
        </div>
      )}

      {/* Improved version */}
      {feedback.improved_version && (
        <div className="p-4 glass rounded-2xl">
          <p className="text-fg-subtle text-xs uppercase tracking-widest mb-2 flex items-center gap-1 font-semibold">
            <CheckCircle className="w-3 h-3" /> Improved Version
          </p>
          <p className="text-fg text-sm leading-relaxed whitespace-pre-wrap">{feedback.improved_version}</p>
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
