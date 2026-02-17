import { useEffect, useState, useCallback, useRef } from 'react'
import { api, type Observation } from '../api/client'
import Timeline from '../components/Timeline'

const PAGE_SIZE = 50

export default function Observations() {
  const [observations, setObservations] = useState<Observation[]>([])
  const [loading, setLoading] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState<string>('all')
  const [hasMore, setHasMore] = useState(true)
  const [offset, setOffset] = useState(0)
  const loaderRef = useRef<HTMLDivElement>(null)

  const loadObservations = useCallback(async (reset = true) => {
    if (reset) {
      setLoading(true)
      setOffset(0)
    } else {
      setLoadingMore(true)
    }
    setError(null)

    try {
      const currentOffset = reset ? 0 : offset
      const params: Record<string, string> = {
        limit: String(PAGE_SIZE),
        offset: String(currentOffset)
      }
      if (typeFilter !== 'all') {
        params.type = typeFilter
      }
      const data = await api.getObservations(params)
      const newData = data || []

      if (reset) {
        setObservations(newData)
        setOffset(PAGE_SIZE)
      } else {
        setObservations(prev => [...prev, ...newData])
        setOffset(prev => prev + PAGE_SIZE)
      }
      setHasMore(newData.length === PAGE_SIZE)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load observations')
    } finally {
      setLoading(false)
      setLoadingMore(false)
    }
  }, [typeFilter, offset])

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      loadObservations(true)
      return
    }

    setLoading(true)
    setError(null)
    setHasMore(false)
    try {
      const data = await api.searchObservations(searchQuery)
      setObservations(data || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Search failed')
    } finally {
      setLoading(false)
    }
  }

  // Initial load
  useEffect(() => {
    loadObservations(true)
  }, [typeFilter])

  // Infinite scroll with IntersectionObserver
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading && !loadingMore && !searchQuery) {
          loadObservations(false)
        }
      },
      { threshold: 0.1 }
    )

    if (loaderRef.current) {
      observer.observe(loaderRef.current)
    }

    return () => observer.disconnect()
  }, [hasMore, loading, loadingMore, searchQuery, loadObservations])

  const types = ['all', 'decision', 'error', 'success', 'insight', 'question']

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">Observations</h1>
        <p className="text-sm text-gray-500">{observations.length} found</p>
      </div>

      {/* Search and Filter */}
      <div className="bg-white rounded-lg shadow p-4">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1">
            <div className="relative">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                placeholder="Search observations..."
                className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              />
              <svg
                className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>
          <button
            onClick={handleSearch}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            Search
          </button>
        </div>

        <div className="mt-4 flex flex-wrap gap-2">
          {types.map((type) => (
            <button
              key={type}
              onClick={() => setTypeFilter(type)}
              className={`px-3 py-1 text-sm rounded-full transition-colors ${
                typeFilter === type
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {type === 'all' ? 'All' : type.charAt(0).toUpperCase() + type.slice(1)}
            </button>
          ))}
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      {/* Observations Timeline */}
      <div className="bg-white rounded-lg shadow p-6">
        <Timeline items={observations} loading={loading} />

        {/* Infinite scroll loader */}
        {!loading && hasMore && !searchQuery && (
          <div ref={loaderRef} className="py-4 text-center">
            {loadingMore ? (
              <div className="flex items-center justify-center gap-2 text-gray-500">
                <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
                <span>Loading more...</span>
              </div>
            ) : (
              <span className="text-gray-400 text-sm">Scroll for more</span>
            )}
          </div>
        )}

        {!loading && !hasMore && observations.length > 0 && (
          <div className="py-4 text-center text-gray-400 text-sm">
            No more observations
          </div>
        )}
      </div>
    </div>
  )
}
