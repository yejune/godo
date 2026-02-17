import { useEffect, useState } from 'react'
import { api, type Session, type Observation } from '../api/client'
import Timeline from '../components/Timeline'

interface Stats {
  sessions: number
  observations: number
  connected: boolean
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats>({ sessions: 0, observations: 0, connected: false })
  const [recentSessions, setRecentSessions] = useState<Session[]>([])
  const [recentObservations, setRecentObservations] = useState<Observation[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadData() {
      setLoading(true)
      setError(null)

      try {
        // Health check
        await api.health()

        // Load data in parallel
        const [sessions, observations] = await Promise.all([
          api.getSessions(),
          api.getObservations({ limit: '10' }),
        ])

        const sessionList = sessions || []
        const observationList = observations || []

        setStats({
          sessions: sessionList.length,
          observations: observationList.length,
          connected: true,
        })
        setRecentSessions(sessionList.slice(0, 5))
        setRecentObservations(observationList.slice(0, 10))
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to connect to Worker API')
        setStats(prev => ({ ...prev, connected: false }))
      } finally {
        setLoading(false)
      }
    }

    loadData()
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Do Memory Dashboard</h1>
        <div className={`flex items-center gap-2 text-sm ${stats.connected ? 'text-green-600' : 'text-red-600'}`}>
          <span className={`w-2 h-2 rounded-full ${stats.connected ? 'bg-green-500' : 'bg-red-500'}`} />
          {stats.connected ? 'Connected' : 'Disconnected'}
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          <p className="font-medium">Connection Error</p>
          <p className="text-sm">{error}</p>
          <p className="text-xs mt-2 text-red-500">
            Make sure the Worker API is running at localhost:3778
          </p>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard
          title="Sessions"
          value={stats.sessions}
          icon="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
          loading={loading}
        />
        <StatCard
          title="Observations"
          value={stats.observations}
          icon="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
          loading={loading}
        />
        <StatCard
          title="Status"
          value={stats.connected ? 'Online' : 'Offline'}
          icon="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
          loading={loading}
          valueColor={stats.connected ? 'text-green-600' : 'text-red-600'}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Sessions */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Sessions</h2>
          {loading ? (
            <div className="space-y-3">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="animate-pulse h-12 bg-gray-100 rounded" />
              ))}
            </div>
          ) : recentSessions.length === 0 ? (
            <p className="text-gray-500 text-center py-4">No sessions yet</p>
          ) : (
            <div className="space-y-3">
              {recentSessions.map((session) => (
                <div
                  key={session.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div>
                    <p className="font-medium text-sm truncate max-w-xs">
                      {session.project_id?.split('/').pop() || session.id || 'Unknown'}
                    </p>
                    <p className="text-xs text-gray-500">
                      {new Date(session.started_at).toLocaleString('ko-KR')}
                    </p>
                  </div>
                  <span className={`text-xs px-2 py-1 rounded ${
                    session.ended_at ? 'bg-gray-200 text-gray-600' : 'bg-green-100 text-green-700'
                  }`}>
                    {session.ended_at ? 'Ended' : 'Active'}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Recent Activity Timeline */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
          <Timeline items={recentObservations} loading={loading} />
        </div>
      </div>
    </div>
  )
}

interface StatCardProps {
  title: string
  value: number | string
  icon: string
  loading?: boolean
  valueColor?: string
}

function StatCard({ title, value, icon, loading, valueColor = 'text-gray-900' }: StatCardProps) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-center gap-4">
        <div className="p-3 bg-primary-100 rounded-lg">
          <svg className="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={icon} />
          </svg>
        </div>
        <div>
          <p className="text-sm text-gray-500">{title}</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className={`text-2xl font-bold ${valueColor}`}>{value}</p>
          )}
        </div>
      </div>
    </div>
  )
}
