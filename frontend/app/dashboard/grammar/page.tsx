'use client'

import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import { cn, getDifficultyColor } from '@/lib/utils'
import { CheckCircle, XCircle, ArrowLeft } from 'lucide-react'

type Exercise = {
  id: string
  title: string
  topic: string
  instruction: string
  difficulty: string
  content?: {
    questions: Question[]
  }
}

type Question = {
  id: string
  text: string
  options: string[]
  correct_answer: string
  explanation: string
}

export default function GrammarPage() {
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(null)

  return (
    <div className="p-4 sm:p-8 max-w-4xl mx-auto">
      <div className="mb-6 sm:mb-8">
        <h1 className="text-2xl sm:text-3xl font-black text-white tracking-tight">Grammar</h1>
        <p className="text-white/45 text-sm mt-1 font-medium">Business IT context — Passive Voice, Conditionals, Modal Verbs</p>
      </div>

      {selectedExercise ? (
        <ExerciseView exercise={selectedExercise} onBack={() => setSelectedExercise(null)} />
      ) : (
        <ExerciseList onSelect={setSelectedExercise} />
      )}
    </div>
  )
}

function ExerciseList({ onSelect }: { onSelect: (e: Exercise) => void }) {
  const { data, isLoading } = useQuery({
    queryKey: ['grammar-exercises'],
    queryFn: () => api.get('/grammar/exercises').then((r) => r.data.data),
  })

  if (isLoading) return <LoadingState />

  return (
    <div className="space-y-3">
      {data?.map((e: Exercise) => (
        <button
          key={e.id}
          onClick={() => onSelect(e)}
          className="w-full text-left p-5 glass glass-hover rounded-2xl transition-all group hover:-translate-y-0.5"
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-white font-black tracking-tight">{e.title}</p>
              <p className="text-white/45 text-sm mt-0.5">{e.topic}</p>
            </div>
            <span className={cn('text-xs px-2 py-0.5 rounded-full font-bold', getDifficultyColor(e.difficulty))}>
              {e.difficulty}
            </span>
          </div>
        </button>
      ))}
    </div>
  )
}

function ExerciseView({ exercise, onBack }: { exercise: Exercise; onBack: () => void }) {
  const queryClient = useQueryClient()
  const [answers, setAnswers] = useState<Record<string, string>>({})
  const [results, setResults] = useState<any[] | null>(null)
  const [score, setScore] = useState<number | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['grammar-exercise', exercise.id],
    queryFn: () => api.get(`/grammar/exercises/${exercise.id}`).then((r) => r.data),
  })

  const submitMutation = useMutation({
    mutationFn: (ans: Record<string, string>) =>
      api.post(`/grammar/exercises/${exercise.id}/submit`, { answers: ans }),
    onSuccess: (res) => {
      setResults(res.data.results)
      setScore(res.data.score)
      queryClient.invalidateQueries({ queryKey: ['grammar-progress'] })
      toast.success(`Score: ${res.data.score}%`)
    },
  })

  if (isLoading) return <LoadingState />

  const questions: Question[] = data?.content?.questions ?? []

  return (
    <div>
      <button
        onClick={onBack}
        className="flex items-center gap-2 text-white/45 hover:text-white text-sm mb-6 transition-colors font-semibold"
      >
        <ArrowLeft className="w-4 h-4" /> Back to exercises
      </button>

      <div className="mb-6">
        <h2 className="text-2xl font-black text-white tracking-tight">{data?.title}</h2>
        <p className="text-white/45 text-sm mt-1">{data?.instruction}</p>
      </div>

      <div className="space-y-6">
        {questions.map((q, idx) => {
          const result = results?.find((r) => r.id === q.id)
          return (
            <div key={q.id} className="p-5 glass rounded-2xl">
              <p className="text-white mb-4 font-medium">
                <span className="text-white/30 mr-2 font-bold">{idx + 1}.</span>
                {q.text}
              </p>
              <div className="grid grid-cols-2 gap-2">
                {q.options.map((opt) => {
                  const selected = answers[q.id] === opt
                  const isCorrect = result?.correct_answer === opt
                  const isWrong = result && selected && !result.correct

                  return (
                    <button
                      key={opt}
                      disabled={!!results}
                      onClick={() => setAnswers({ ...answers, [q.id]: opt })}
                      className={cn(
                        'px-4 py-2.5 rounded-xl text-sm text-left transition-all border font-semibold',
                        !results && selected
                          ? 'bg-violet-600/20 border-violet-500 text-violet-700 dark:text-violet-200'
                          : !results
                          ? 'glass hover:border-violet-500/30 text-fg-muted hover:text-fg'
                          : isCorrect
                          ? 'bg-green-500/15 border-green-500/50 text-green-700 dark:text-green-300'
                          : isWrong
                          ? 'bg-red-500/15 border-red-500/50 text-red-700 dark:text-red-300'
                          : 'glass text-fg-subtle'
                      )}
                    >
                      {opt}
                    </button>
                  )
                })}
              </div>
              {result && (
                <div className={cn('mt-3 p-3 rounded-xl text-sm font-medium', result.correct ? 'bg-green-500/10 text-green-300' : 'bg-red-500/10 text-red-300')}>
                  {result.correct ? <CheckCircle className="inline w-4 h-4 mr-1" /> : <XCircle className="inline w-4 h-4 mr-1" />}
                  {q.explanation}
                </div>
              )}
            </div>
          )
        })}
      </div>

      {!results && (
        <button
          onClick={() => submitMutation.mutate(answers)}
          disabled={Object.keys(answers).length < questions.length || submitMutation.isPending}
          className="mt-8 w-full py-3 btn-vibrant disabled:opacity-40 disabled:cursor-not-allowed disabled:transform-none text-white font-black rounded-2xl transition-all"
        >
          {submitMutation.isPending ? 'Submitting...' : 'Submit Answers'}
        </button>
      )}

      {score !== null && (
        <div className={cn('mt-6 p-5 rounded-2xl text-center font-black text-lg tracking-tight', score >= 80 ? 'panel-success' : score >= 50 ? 'panel-warning' : 'panel-danger')}>
          Final Score: {score}%
        </div>
      )}
    </div>
  )
}

function LoadingState() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map((i) => (
        <div key={i} className="h-20 glass rounded-2xl animate-pulse" />
      ))}
    </div>
  )
}
