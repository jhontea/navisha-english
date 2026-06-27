'use client'

import { useState, useRef, useEffect } from 'react'
import { useQuery, useMutation } from '@tanstack/react-query'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import { cn, getDifficultyColor } from '@/lib/utils'
import { ArrowLeft, Send } from 'lucide-react'

type Scenario = {
  id: string
  title: string
  context: string
  ai_role: string
  user_role: string
  difficulty: string
  tags: string[]
}

type Message = {
  role: 'user' | 'assistant'
  content: string
  created_at: string
}

export default function RoleplayPage() {
  const [sessionId, setSessionId] = useState<string | null>(null)
  const [scenario, setScenario] = useState<Scenario | null>(null)
  const [messages, setMessages] = useState<Message[]>([])

  const startSession = (s: Scenario, sid: string, msgs: Message[]) => {
    setScenario(s)
    setSessionId(sid)
    setMessages(msgs)
  }

  return (
    <div className="p-4 sm:p-8 max-w-4xl mx-auto">
      <div className="mb-6 sm:mb-8">
        <h1 className="text-2xl sm:text-3xl font-black text-white tracking-tight">Role-play Scenarios</h1>
        <p className="text-white/45 text-sm mt-1 font-medium">Practice real workplace conversations with AI</p>
      </div>

      {sessionId && scenario ? (
        <ChatView
          sessionId={sessionId}
          scenario={scenario}
          messages={messages}
          setMessages={setMessages}
          onBack={() => { setSessionId(null); setScenario(null); setMessages([]) }}
        />
      ) : (
        <ScenarioList onStart={startSession} />
      )}
    </div>
  )
}

function ScenarioList({ onStart }: {
  onStart: (s: Scenario, sid: string, msgs: Message[]) => void
}) {
  const { data, isLoading } = useQuery({
    queryKey: ['roleplay-scenarios'],
    queryFn: () => api.get('/roleplay/scenarios').then((r) => r.data.data),
  })

  const startMutation = useMutation({
    mutationFn: (scenario: Scenario) =>
      api.post(`/roleplay/scenarios/${scenario.id}/start`),
    onSuccess: (res, scenario) => {
      onStart(scenario, res.data.session_id, res.data.messages)
      toast.success('Scenario started!')
    },
    onError: () => toast.error('Failed to start scenario'),
  })

  if (isLoading) return <LoadingState />

  return (
    <div className="space-y-4">
      {data?.map((s: Scenario) => (
        <div key={s.id} className="p-5 glass glass-hover rounded-2xl transition-all">
          <div className="flex items-start justify-between mb-3">
            <div className="flex-1">
              <h3 className="text-white font-black mb-1 tracking-tight">{s.title}</h3>
              <p className="text-white/45 text-sm line-clamp-2">{s.context}</p>
            </div>
            <span className={cn('ml-4 shrink-0 text-xs px-2 py-0.5 rounded-full font-bold', getDifficultyColor(s.difficulty))}>
              {s.difficulty}
            </span>
          </div>

          <div className="flex items-center gap-4 mb-4 text-xs text-white/30 font-medium">
            <span>You: <span className="text-white/60">{s.user_role}</span></span>
            <span>AI: <span className="text-white/60">{s.ai_role}</span></span>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex gap-1.5 flex-wrap">
              {s.tags?.slice(0, 3).map((tag) => (
                <span key={tag} className="text-xs px-2 py-0.5 glass text-white/40 rounded-lg font-medium">
                  {tag}
                </span>
              ))}
            </div>
            <button
              onClick={() => startMutation.mutate(s)}
              disabled={startMutation.isPending}
              className="px-4 py-2 btn-vibrant disabled:opacity-50 text-white text-sm font-bold rounded-xl transition-all"
            >
              {startMutation.isPending ? 'Starting...' : 'Start Scenario'}
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}

function ChatView({ sessionId, scenario, messages, setMessages, onBack }: {
  sessionId: string
  scenario: Scenario
  messages: Message[]
  setMessages: (msgs: Message[]) => void
  onBack: () => void
}) {
  const [input, setInput] = useState('')
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const sendMutation = useMutation({
    mutationFn: (content: string) =>
      api.post(`/roleplay/sessions/${sessionId}/message`, { content }),
    onSuccess: (res) => {
      setMessages(res.data.messages)
      setInput('')
    },
    onError: () => toast.error('Failed to send message'),
  })

  const handleSend = () => {
    const trimmed = input.trim()
    if (!trimmed || sendMutation.isPending) return
    sendMutation.mutate(trimmed)
  }

  return (
    <div className="flex flex-col h-[calc(100vh-14rem)]">
      {/* Header */}
      <div className="flex items-center gap-4 mb-4">
        <button
          onClick={onBack}
          className="flex items-center gap-2 text-white/45 hover:text-white text-sm transition-colors font-semibold"
        >
          <ArrowLeft className="w-4 h-4" /> Back
        </button>
        <div className="flex-1 p-3 glass rounded-2xl">
          <p className="text-white text-sm font-black tracking-tight">{scenario.title}</p>
          <p className="text-white/35 text-xs mt-0.5 font-medium">
            You: {scenario.user_role} · AI: {scenario.ai_role}
          </p>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-y-auto space-y-4 mb-4 pr-1">
        {messages.map((msg, idx) => (
          <div key={idx} className={cn('flex', msg.role === 'user' ? 'justify-end' : 'justify-start')}>
            <div className={cn(
              'max-w-[75%] px-4 py-3 rounded-2xl text-sm leading-relaxed font-medium',
              msg.role === 'user'
                ? 'btn-vibrant rounded-br-sm'
                : 'glass text-fg rounded-bl-sm'
            )}>
              {msg.role === 'assistant' && (
                <p className="text-fg-subtle text-xs mb-1 font-bold uppercase tracking-wider">{scenario.ai_role}</p>
              )}
              {msg.content}
            </div>
          </div>
        ))}
        {sendMutation.isPending && (
          <div className="flex justify-start">
            <div className="glass text-fg-subtle px-4 py-3 rounded-2xl rounded-bl-sm text-sm font-medium">
              <span className="animate-pulse">Typing...</span>
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </div>

      {/* Input */}
      <div className="flex gap-3">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && handleSend()}
          placeholder={`Reply as ${scenario.user_role}...`}
          className="flex-1 px-4 py-3 input-glass rounded-2xl focus:outline-none text-sm font-medium"
        />
        <button
          onClick={handleSend}
          disabled={!input.trim() || sendMutation.isPending}
          className="px-4 py-3 btn-vibrant disabled:opacity-40 disabled:cursor-not-allowed disabled:transform-none text-white rounded-2xl transition-all"
        >
          <Send className="w-4 h-4" />
        </button>
      </div>
    </div>
  )
}

function LoadingState() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map((i) => (
        <div key={i} className="h-32 glass rounded-2xl animate-pulse" />
      ))}
    </div>
  )
}
