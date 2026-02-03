export default function AboutPage() {
  return (
    <section className="page">
      <h1>About</h1>
      <p>
        I design backend systems that survive contact with reality: data volume,
        strange edge cases, and teams that need observable, reliable primitives.
      </p>
      <p>
        This site is a public ledger of decisions, trade-offs, and the patterns
        I trust after operating software in production.
      </p>
      <div className="facts">
        <div>
          <span className="label">Focus</span>
          <span>Distributed systems, observability, pragmatic DDD.</span>
        </div>
        <div>
          <span className="label">Stack</span>
          <span>Go, Postgres, queues, metrics, sharp edges.</span>
        </div>
        <div>
          <span className="label">Current mode</span>
          <span>Writing, reviewing, shipping.</span>
        </div>
      </div>
    </section>
  )
}
