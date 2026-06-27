import Link from 'next/link'
import { ThemeToggle } from './components/ThemeToggle'

export default function HomePage() {
  return (
    <main className="min-h-screen flex flex-col">
      {/* Nav */}
      <nav className="flex items-center justify-between px-4 sm:px-8 py-4 sm:py-5 max-w-6xl mx-auto w-full">
        <span className="text-white font-extrabold text-lg sm:text-xl tracking-tight">
          Navisha<span className="gradient-text"> English</span>
        </span>
        <div className="flex gap-2 sm:gap-3 items-center">
          <ThemeToggle />
          <Link
            href="/login"
            className="px-3 sm:px-5 py-2 text-sm font-semibold text-white/70 hover:text-white transition-colors"
          >
            Sign In
          </Link>
          <Link
            href="/register"
            className="px-3 sm:px-5 py-2 btn-vibrant text-white font-bold text-sm rounded-xl"
          >
            Get Started
          </Link>
        </div>
      </nav>

      <div className="flex-1 max-w-6xl mx-auto px-4 sm:px-6 py-10 sm:py-16 w-full">
        {/* Hero */}
        <div className="text-center mb-12 sm:mb-20">
          <div className="inline-flex items-center gap-2 px-3 sm:px-4 py-1.5 rounded-full glass text-violet-300 text-xs sm:text-sm font-semibold mb-6 sm:mb-8 tracking-wide uppercase">
            ✦ Business English for IT Professionals
          </div>
          <h1 className="text-4xl sm:text-6xl md:text-7xl font-black text-white mb-5 sm:mb-6 leading-[1.05] tracking-tight">
            Speak the language<br />
            of{' '}
            <span className="gradient-text">global tech</span>
          </h1>
          <p className="text-base sm:text-lg text-white/55 max-w-2xl mx-auto mb-8 sm:mb-12 leading-relaxed px-2">
            Master Business English at B1–C1 level. Built for developers, engineers, and IT professionals who communicate with confidence.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center items-center">
            <Link
              href="/register"
              className="w-full sm:w-auto px-8 sm:px-9 py-3.5 sm:py-4 btn-vibrant text-white font-bold rounded-2xl text-base text-center"
            >
              Start Learning Free →
            </Link>
            <Link
              href="/login"
              className="w-full sm:w-auto px-8 sm:px-9 py-3.5 sm:py-4 glass glass-hover text-white font-bold rounded-2xl text-base transition-all text-center"
            >
              Sign In
            </Link>
          </div>

          {/* Social proof */}
          <p className="mt-6 sm:mt-8 text-white/30 text-xs sm:text-sm">No credit card required · Google login · Start in 30 seconds</p>
        </div>

        {/* Feature Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
          {features.map((f) => (
            <div
              key={f.title}
              className="group p-5 sm:p-7 rounded-3xl glass glass-hover transition-all duration-300 hover:-translate-y-1"
            >
              <div className={`w-11 h-11 sm:w-12 sm:h-12 rounded-2xl flex items-center justify-center text-xl sm:text-2xl mb-4 sm:mb-5 ${f.iconBg}`}>
                {f.icon}
              </div>
              <h3 className="text-lg sm:text-xl font-black text-white mb-2 tracking-tight">{f.title}</h3>
              <p className="text-white/50 text-sm leading-relaxed">{f.description}</p>
            </div>
          ))}
        </div>
      </div>
    </main>
  )
}

const features = [
  {
    icon: '✍️',
    iconBg: 'bg-violet-500/15',
    title: 'Business Writing',
    description: 'Write professional emails, technical proposals, sprint reports, and client communications with AI feedback on grammar and tone.',
  },
  {
    icon: '🃏',
    iconBg: 'bg-cyan-500/15',
    title: 'Vocabulary Builder',
    description: 'Learn 200+ IT business terms using spaced repetition (SM-2). From "scope creep" to "stakeholder" — in real IT context.',
  },
  {
    icon: '📐',
    iconBg: 'bg-green-500/15',
    title: 'Grammar for Professionals',
    description: 'Passive voice, conditionals, modal verbs — all taught through real IT and business scenarios, not generic textbook examples.',
  },
  {
    icon: '🎭',
    iconBg: 'bg-pink-500/15',
    title: 'Role-play Scenarios',
    description: 'Simulate real workplace conversations: presenting to a manager, negotiating deadlines, code reviews, and daily standups with AI.',
  },
]
