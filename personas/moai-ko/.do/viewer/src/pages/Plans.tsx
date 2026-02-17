import { useEffect, useState } from 'react'
import { api, type Plan } from '../api/client'

export default function Plans() {
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null)

  useEffect(() => {
    async function loadPlans() {
      setLoading(true)
      try {
        const data = await api.getPlans()
        setPlans(data || [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load plans')
      } finally {
        setLoading(false)
      }
    }

    loadPlans()
  }, [])

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-600',
      active: 'bg-blue-100 text-blue-700',
      completed: 'bg-green-100 text-green-700',
      archived: 'bg-yellow-100 text-yellow-700',
    }
    return colors[status] || colors.draft
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Plans</h1>
        <p className="text-sm text-gray-500">{plans.length} total</p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Plans List */}
        <div className="lg:col-span-1 bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h2 className="font-semibold">All Plans</h2>
          </div>
          <div className="divide-y max-h-[600px] overflow-y-auto">
            {loading ? (
              [...Array(5)].map((_, i) => (
                <div key={i} className="p-4 animate-pulse">
                  <div className="h-4 bg-gray-200 rounded w-2/3 mb-2" />
                  <div className="h-3 bg-gray-100 rounded w-1/3" />
                </div>
              ))
            ) : plans.length === 0 ? (
              <div className="p-8 text-center text-gray-500">
                No plans found
              </div>
            ) : (
              plans.map((plan) => (
                <button
                  key={plan.id}
                  onClick={() => setSelectedPlan(plan)}
                  className={`w-full p-4 text-left hover:bg-gray-50 transition-colors ${
                    selectedPlan?.id === plan.id ? 'bg-primary-50' : ''
                  }`}
                >
                  <div className="flex items-start justify-between">
                    <p className="font-medium text-sm line-clamp-2">{plan.title}</p>
                    <span className={`text-xs px-2 py-0.5 rounded ml-2 flex-shrink-0 ${getStatusColor(plan.status)}`}>
                      {plan.status}
                    </span>
                  </div>
                  <p className="text-xs text-gray-500 mt-1">
                    {new Date(plan.created_at).toLocaleDateString('ko-KR')}
                  </p>
                </button>
              ))
            )}
          </div>
        </div>

        {/* Plan Content */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h2 className="font-semibold">Plan Details</h2>
          </div>
          <div className="p-4">
            {selectedPlan ? (
              <div className="space-y-4">
                <div className="flex items-start justify-between">
                  <h3 className="text-lg font-semibold">{selectedPlan.title}</h3>
                  <span className={`text-xs px-2 py-1 rounded ${getStatusColor(selectedPlan.status)}`}>
                    {selectedPlan.status}
                  </span>
                </div>
                <p className="text-xs text-gray-500">
                  Created: {new Date(selectedPlan.created_at).toLocaleString('ko-KR')}
                </p>
                <div className="border-t pt-4">
                  <pre className="text-sm text-gray-700 whitespace-pre-wrap font-sans">
                    {selectedPlan.content}
                  </pre>
                </div>
              </div>
            ) : (
              <p className="text-gray-500 text-center py-16">
                Select a plan to view details
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
