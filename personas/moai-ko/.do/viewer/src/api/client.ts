const WORKER_URL = 'http://127.0.0.1:3778'

export interface Project {
  id: string
  path: string
  session_count: number
  last_activity: string
}

export interface Session {
  id: string
  user_name?: string
  project_id?: string
  started_at: string
  ended_at?: string
  updated_at?: string
  summary?: string
}

export interface Observation {
  id: number
  session_id: string
  created_at: string
  type: string
  content: string
  metadata?: Record<string, unknown>
}

export interface Plan {
  id: string
  title: string
  created_at: string
  status: string
  content: string
}

export interface TeamMember {
  name: string
  last_active: string
  sessions_count: number
}

export interface Summary {
  date: string
  sessions_count: number
  observations_count: number
  highlights?: string[]
}

async function fetchJson<T>(url: string): Promise<T> {
  const response = await fetch(url)
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }
  return response.json()
}

export const api = {
  // Projects
  getProjects: () =>
    fetchJson<Project[]>(`${WORKER_URL}/api/projects`),

  // Sessions
  getSessions: (projectId?: string) => {
    const params = projectId ? `?project_id=${encodeURIComponent(projectId)}` : ''
    return fetchJson<Session[]>(`${WORKER_URL}/api/sessions${params}`)
  },

  getSession: (id: string) =>
    fetchJson<Session>(`${WORKER_URL}/api/sessions/${id}`),

  // Observations
  getObservations: (params?: Record<string, string>) => {
    const query = params ? new URLSearchParams(params).toString() : ''
    return fetchJson<Observation[]>(`${WORKER_URL}/api/observations${query ? `?${query}` : ''}`)
  },

  searchObservations: (q: string) =>
    fetchJson<Observation[]>(`${WORKER_URL}/api/observations/search?q=${encodeURIComponent(q)}`),

  // Team
  getTeamContext: (projectPath: string) =>
    fetchJson<TeamMember[]>(`${WORKER_URL}/api/team/context?project_path=${encodeURIComponent(projectPath)}`),

  // Plans
  getPlans: () =>
    fetchJson<Plan[]>(`${WORKER_URL}/api/plans`),

  // Summaries
  getSummaries: (days?: number) =>
    fetchJson<Summary[]>(`${WORKER_URL}/api/summaries?days=${days || 7}`),

  // Health check
  health: () =>
    fetchJson<{ status: string }>(`${WORKER_URL}/health`),
}
