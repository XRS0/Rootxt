import { NavLink, Outlet } from 'react-router-dom'

export default function Layout() {
  return (
    <div className="shell">
      <header className="site-header">
        <div className="brand">SYSTEM JOURNAL</div>
        <nav className="site-nav">
          <NavLink to="/" className={({ isActive }) => (isActive ? 'active' : '')}>
            Home
          </NavLink>
          <NavLink to="/articles" className={({ isActive }) => (isActive ? 'active' : '')}>
            Articles
          </NavLink>
          <NavLink to="/about" className={({ isActive }) => (isActive ? 'active' : '')}>
            About
          </NavLink>
          <NavLink to="/contact" className={({ isActive }) => (isActive ? 'active' : '')}>
            Contact
          </NavLink>
        </nav>
      </header>
      <main className="site-main">
        <Outlet />
      </main>
      <footer className="site-footer">
        <div>Last sync: {new Date().getFullYear()}</div>
        <div>Systems over noise.</div>
      </footer>
    </div>
  )
}
