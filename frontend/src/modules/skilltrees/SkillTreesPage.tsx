import { useState, useEffect } from 'react'
import {
  GitBranch, Loader2, ArrowLeft, ChevronRight,
  Clock, AlertTriangle, Star
} from 'lucide-react'

interface Step {
  order: number
  title: string
  description: string
}

interface Skill {
  id: string
  name: string
  difficulty: string
  time_estimate: string
  prerequisites?: string[]
  summary: string
  steps: Step[]
  tips?: string[]
}

interface SkillCategory {
  id: string
  name: string
  icon: string
  description: string
  skills: Skill[]
}

const difficultyStyles: Record<string, { bg: string; border: string; text: string }> = {
  beginner: { bg: 'bg-emerald-500/10', border: 'border-emerald-500/30', text: 'text-emerald-400' },
  intermediate: { bg: 'bg-amber-500/10', border: 'border-amber-500/30', text: 'text-amber-400' },
  advanced: { bg: 'bg-red-500/10', border: 'border-red-500/30', text: 'text-red-400' },
}

export function SkillTreesPage() {
  const [categories, setCategories] = useState<SkillCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [activeCategory, setActiveCategory] = useState<SkillCategory | null>(null)
  const [activeSkill, setActiveSkill] = useState<Skill | null>(null)

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const res = await fetch('/api/skills/categories')
        const data = await res.json()
        setCategories(data.categories || [])
      } catch (e) {
        console.error('Failed to fetch skill categories:', e)
      } finally {
        setLoading(false)
      }
    }
    fetchCategories()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <Loader2 className="w-8 h-8 text-emerald-500 animate-spin" />
      </div>
    )
  }

  // Skill detail view
  if (activeSkill) {
    const diff = difficultyStyles[activeSkill.difficulty] || difficultyStyles.beginner
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveSkill(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> Back to {activeCategory?.name}
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeSkill.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeSkill.summary}</p>
          <div className="flex items-center gap-3 mt-3">
            <span className={`px-2.5 py-0.5 rounded-full text-xs font-semibold uppercase ${diff.bg} ${diff.border} border ${diff.text}`}>
              {activeSkill.difficulty}
            </span>
            <span className="flex items-center gap-1 text-xs text-slate-500">
              <Clock className="w-3 h-3" /> {activeSkill.time_estimate}
            </span>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Prerequisites */}
          {activeSkill.prerequisites && activeSkill.prerequisites.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider mb-2">Prerequisites</h2>
              <div className="flex flex-wrap gap-2">
                {activeSkill.prerequisites.map((p, i) => (
                  <span key={i} className="px-2.5 py-1 rounded-lg bg-slate-800/50 border border-slate-700/30 text-xs text-slate-400">
                    {p}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* Steps */}
          <div>
            <h2 className="text-sm font-semibold text-slate-300 uppercase tracking-wider mb-3">Instructions</h2>
            <div className="space-y-3">
              {activeSkill.steps.map((step) => (
                <div key={step.order} className="flex gap-4 p-4 rounded-xl bg-slate-800/20 border border-slate-700/20">
                  <div className="w-8 h-8 rounded-full bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center shrink-0">
                    <span className="text-sm font-bold text-emerald-400">{step.order}</span>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-sm font-semibold text-slate-200">{step.title}</h3>
                    <p className="text-sm text-slate-400 mt-1 leading-relaxed">{step.description}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Tips */}
          {activeSkill.tips && activeSkill.tips.length > 0 && (
            <div>
              <h2 className="text-sm font-semibold text-amber-400 uppercase tracking-wider mb-3">
                <Star className="w-3.5 h-3.5 inline mr-1" /> Pro Tips
              </h2>
              <div className="space-y-2">
                {activeSkill.tips.map((tip, i) => (
                  <div key={i} className="flex gap-3 p-3 rounded-lg bg-amber-500/5 border border-amber-500/20">
                    <AlertTriangle className="w-4 h-4 text-amber-400 shrink-0 mt-0.5" />
                    <p className="text-sm text-amber-300/80 leading-relaxed">{tip}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    )
  }

  // Category skills list
  if (activeCategory) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-6 pb-4 border-b border-slate-800/50">
          <button
            onClick={() => setActiveCategory(null)}
            className="flex items-center gap-2 text-sm text-slate-400 hover:text-slate-200 mb-3 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" /> All Categories
          </button>
          <h1 className="text-xl font-semibold text-slate-100">{activeCategory.name}</h1>
          <p className="text-sm text-slate-400 mt-1">{activeCategory.description}</p>
        </div>

        <div className="flex-1 overflow-y-auto p-6 space-y-3">
          {activeCategory.skills.map(skill => {
            const diff = difficultyStyles[skill.difficulty] || difficultyStyles.beginner
            return (
              <button
                key={skill.id}
                onClick={() => setActiveSkill(skill)}
                className="w-full text-left p-4 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-emerald-500/30 hover:bg-emerald-500/5 transition-all group"
              >
                <div className="flex items-center justify-between">
                  <span className="text-sm font-semibold text-slate-200 group-hover:text-emerald-300">{skill.name}</span>
                  <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-emerald-400" />
                </div>
                <p className="text-xs text-slate-500 mt-1">{skill.summary}</p>
                <div className="flex items-center gap-3 mt-2">
                  <span className={`px-2 py-0.5 rounded-full text-[10px] font-semibold uppercase ${diff.bg} ${diff.border} border ${diff.text}`}>
                    {skill.difficulty}
                  </span>
                  <span className="text-[10px] text-slate-600 flex items-center gap-1">
                    <Clock className="w-3 h-3" /> {skill.time_estimate}
                  </span>
                  <span className="text-[10px] text-slate-600">{skill.steps.length} steps</span>
                </div>
              </button>
            )
          })}
        </div>
      </div>
    )
  }

  // Categories overview
  const totalSkills = categories.reduce((sum, c) => sum + c.skills.length, 0)

  return (
    <div className="flex flex-col h-full">
      <div className="p-6 pb-4 border-b border-slate-800/50">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-lg bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center">
            <GitBranch className="w-5 h-5 text-emerald-400" />
          </div>
          <div>
            <h1 className="text-xl font-semibold text-slate-100">Skill Trees</h1>
            <p className="text-sm text-slate-400">
              {categories.length} categories · {totalSkills} survival skills
            </p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {categories.map(cat => (
            <button
              key={cat.id}
              onClick={() => setActiveCategory(cat)}
              className="text-left p-5 rounded-xl bg-slate-800/30 border border-slate-700/30 hover:border-emerald-500/30 hover:bg-emerald-500/5 transition-all group"
            >
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-semibold text-slate-200 group-hover:text-emerald-300">
                  {cat.name}
                </h3>
                <ChevronRight className="w-4 h-4 text-slate-600 group-hover:text-emerald-400" />
              </div>
              <p className="text-xs text-slate-500 mb-3">{cat.description}</p>
              <span className="text-[10px] px-2 py-0.5 rounded-full bg-slate-700/50 text-slate-400">
                {cat.skills.length} {cat.skills.length === 1 ? 'skill' : 'skills'}
              </span>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
