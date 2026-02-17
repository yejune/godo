import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { ProjectProvider } from './context/ProjectContext'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Sessions from './pages/Sessions'
import Observations from './pages/Observations'
import Plans from './pages/Plans'
import Reports from './pages/Reports'
import UserPrompts from './pages/UserPrompts'

export default function App() {
  return (
    <ProjectProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="sessions" element={<Sessions />} />
            <Route path="observations" element={<Observations />} />
            <Route path="plans" element={<Plans />} />
            <Route path="reports" element={<Reports />} />
            <Route path="prompts" element={<UserPrompts />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ProjectProvider>
  )
}
