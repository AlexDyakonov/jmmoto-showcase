import { useState } from "react";

function App() {
  const [count, setCount] = useState(0);

  const increment = () => setCount((prev) => prev + 1);
  const decrement = () => setCount((prev) => prev - 1);
  const reset = () => setCount(0);

  return (
    <div className="min-h-screen bg-[var(--color-background)] text-[var(--color-foreground)] flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        <div className="bg-[var(--color-card)] rounded-[var(--radius)] p-8 shadow-lg border border-[var(--color-border)]">
          <h1 className="text-3xl font-bold text-center mb-2">
            Counter
          </h1>
          <p className="text-[var(--color-muted-foreground)] text-center mb-8">
            Simple app built with React + Vite + TypeScript
          </p>

          <div className="flex flex-col items-center gap-6">
            <div className="text-6xl font-bold font-mono text-[var(--color-primary)]">
              {count}
            </div>

            <div className="flex gap-3 w-full">
              <button
                onClick={decrement}
                className="flex-1 bg-[var(--color-secondary)] hover:bg-[var(--color-accent)] text-[var(--color-secondary-foreground)] font-semibold py-3 px-6 rounded-lg transition-colors duration-200 active:scale-95"
              >
                −
              </button>
              <button
                onClick={reset}
                className="flex-1 bg-[var(--color-muted)] hover:bg-[var(--color-accent)] text-[var(--color-muted-foreground)] font-semibold py-3 px-6 rounded-lg transition-colors duration-200 active:scale-95"
              >
                Reset
              </button>
              <button
                onClick={increment}
                className="flex-1 bg-[var(--color-primary)] hover:opacity-90 text-[var(--color-primary-foreground)] font-semibold py-3 px-6 rounded-lg transition-all duration-200 active:scale-95"
              >
                +
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;

