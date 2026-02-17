import { useEffect, useState } from 'react'

interface UserPrompt {
  id: number
  session_id: string
  prompt_number: number
  prompt_text: string
  response?: string
  created_at: string
}

export default function UserPrompts() {
  const [prompts, setPrompts] = useState<UserPrompt[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedPrompt, setSelectedPrompt] = useState<UserPrompt | null>(null)

  useEffect(() => {
    async function loadPrompts() {
      setLoading(true)
      try {
        const response = await fetch('http://127.0.0.1:3778/api/prompts?limit=100')
        if (!response.ok) throw new Error(`HTTP ${response.status}`)
        const data = await response.json()
        setPrompts(Array.isArray(data) ? data : [])
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load prompts')
      } finally {
        setLoading(false)
      }
    }

    loadPrompts()
  }, [])

  // Group by session
  const groupedBySession = prompts.reduce((acc, p) => {
    if (!acc[p.session_id]) acc[p.session_id] = []
    acc[p.session_id].push(p)
    return acc
  }, {} as Record<string, UserPrompt[]>)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">User Prompts</h1>
        <p className="text-sm text-gray-500">{prompts.length} total</p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
          {error}
        </div>
      )}

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Total Prompts</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">{prompts.length}</p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Sessions</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">
              {Object.keys(groupedBySession).length}
            </p>
          )}
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">With Response</p>
          {loading ? (
            <div className="animate-pulse h-8 w-16 bg-gray-200 rounded mt-1" />
          ) : (
            <p className="text-2xl font-bold text-gray-900">
              {prompts.filter(p => p.response).length}
            </p>
          )}
        </div>
      </div>

      {/* Prompts List */}
      <div className="bg-white rounded-lg shadow">
        <div className="p-4 border-b">
          <h2 className="font-semibold">Recent Prompts (Q&A)</h2>
        </div>
        <div className="divide-y max-h-[600px] overflow-y-auto">
          {loading ? (
            [...Array(5)].map((_, i) => (
              <div key={i} className="p-4 animate-pulse">
                <div className="h-4 bg-gray-200 rounded w-1/3 mb-2" />
                <div className="h-3 bg-gray-100 rounded w-2/3" />
              </div>
            ))
          ) : prompts.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No prompts recorded yet
            </div>
          ) : (
            prompts.map((prompt) => (
              <button
                key={prompt.id}
                onClick={() => setSelectedPrompt(prompt)}
                className="w-full text-left p-4 hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-mono text-gray-400">
                    {prompt.session_id.slice(0, 8)}...
                  </span>
                  <div className="flex items-center gap-2">
                    {prompt.response && (
                      <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded">
                        응답 있음
                      </span>
                    )}
                    <span className="text-xs text-gray-400">
                      {new Date(prompt.created_at).toLocaleString('ko-KR')}
                    </span>
                  </div>
                </div>
                <p className="text-sm text-gray-900 font-medium line-clamp-2 mb-1">
                  Q: {prompt.prompt_text}
                </p>
                {prompt.response && (
                  <p className="text-sm text-gray-500 line-clamp-2">
                    A: {prompt.response.slice(0, 200)}...
                  </p>
                )}
              </button>
            ))
          )}
        </div>
      </div>

      {/* Detail Modal */}
      {selectedPrompt && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[85vh] flex flex-col">
            <div className="p-4 border-b flex items-center justify-between">
              <div>
                <h3 className="font-semibold">Prompt Detail</h3>
                <p className="text-xs text-gray-500 font-mono">
                  {selectedPrompt.session_id}
                </p>
              </div>
              <button
                onClick={() => setSelectedPrompt(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div className="p-4 overflow-y-auto flex-1 space-y-4">
              <p className="text-xs text-gray-400">
                {new Date(selectedPrompt.created_at).toLocaleString('ko-KR')}
              </p>

              {/* Question */}
              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2 flex items-center gap-2">
                  <span className="w-6 h-6 bg-blue-100 text-blue-700 rounded-full flex items-center justify-center text-xs font-bold">Q</span>
                  질문
                </h4>
                <pre className="text-sm text-gray-700 bg-blue-50 p-4 rounded-lg whitespace-pre-wrap">
                  {selectedPrompt.prompt_text}
                </pre>
              </div>

              {/* Response */}
              <div>
                <h4 className="text-sm font-medium text-gray-900 mb-2 flex items-center gap-2">
                  <span className="w-6 h-6 bg-green-100 text-green-700 rounded-full flex items-center justify-center text-xs font-bold">A</span>
                  응답
                  {selectedPrompt.response && (
                    <span className="text-xs text-gray-400">
                      ({selectedPrompt.response.length < 1000
                        ? `${selectedPrompt.response.length}B`
                        : `${Math.round(selectedPrompt.response.length / 1000)}KB`})
                    </span>
                  )}
                </h4>
                {selectedPrompt.response ? (
                  <pre className="text-sm text-gray-700 bg-green-50 p-4 rounded-lg whitespace-pre-wrap font-mono max-h-[400px] overflow-y-auto">
                    {selectedPrompt.response}
                  </pre>
                ) : (
                  <div className="text-center py-8 text-gray-400 bg-gray-50 rounded-lg">
                    응답이 저장되지 않았습니다
                  </div>
                )}
              </div>
            </div>

            <div className="p-4 border-t flex justify-end">
              <button
                onClick={() => setSelectedPrompt(null)}
                className="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
              >
                닫기
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
