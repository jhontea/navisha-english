import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | Date) {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(new Date(date))
}

export function getDifficultyColor(difficulty: string) {
  // Uses colors that are readable in both dark and light mode
  switch (difficulty) {
    case 'A2': return 'text-emerald-600 dark:text-emerald-400 bg-emerald-500/10'
    case 'B1': return 'text-blue-600 dark:text-blue-400 bg-blue-500/10'
    case 'B2': return 'text-amber-600 dark:text-amber-400 bg-amber-500/10'
    case 'C1': return 'text-red-600 dark:text-red-400 bg-red-500/10'
    default:   return 'text-slate-600 dark:text-slate-400 bg-slate-500/10'
  }
}
