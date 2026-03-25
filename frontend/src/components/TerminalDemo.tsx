import React, { useEffect, useRef, useState, useCallback } from 'react';

export type TerminalStep = {
  command: string;
  output: string[];
};

type Props = {
  steps: TerminalStep[];
  title?: string;
};

type AnimState = {
  stepIndex: number;
  typedChars: number;
  outputLines: number;
  phase: 'typing' | 'pause-after-cmd' | 'revealing' | 'pause-after-step' | 'done';
};

const CHAR_INTERVAL = 15;
const POST_CMD_PAUSE = 400;
const LINE_INTERVAL = 120;
const POST_STEP_PAUSE = 800;

const TerminalDemo: React.FC<Props> = ({ steps, title = 'goini — terminal' }) => {
  const initial: AnimState = {
    stepIndex: 0,
    typedChars: 0,
    outputLines: 0,
    phase: 'typing',
  };

  const [anim, setAnim] = useState<AnimState>(initial);
  const [tick, setTick] = useState(0);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const reset = useCallback(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    setAnim(initial);
    setTick(t => t + 1);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (anim.phase === 'done') return;

    const step = steps[anim.stepIndex];

    const schedule = (delay: number, fn: () => void) => {
      timerRef.current = setTimeout(fn, delay);
    };

    if (anim.phase === 'typing') {
      if (anim.typedChars < step.command.length) {
        schedule(CHAR_INTERVAL, () =>
          setAnim(s => ({ ...s, typedChars: s.typedChars + 1 }))
        );
      } else {
        schedule(POST_CMD_PAUSE, () =>
          setAnim(s => ({ ...s, phase: step.output.length > 0 ? 'revealing' : 'pause-after-step' }))
        );
      }
    } else if (anim.phase === 'revealing') {
      if (anim.outputLines < step.output.length) {
        schedule(LINE_INTERVAL, () =>
          setAnim(s => ({ ...s, outputLines: s.outputLines + 1 }))
        );
      } else {
        schedule(POST_STEP_PAUSE, () =>
          setAnim(s => ({ ...s, phase: 'pause-after-step' }))
        );
      }
    } else if (anim.phase === 'pause-after-step') {
      const nextIndex = anim.stepIndex + 1;
      if (nextIndex < steps.length) {
        setAnim({
          stepIndex: nextIndex,
          typedChars: 0,
          outputLines: 0,
          phase: 'typing',
        });
      } else {
        setAnim(s => ({ ...s, phase: 'done' }));
      }
    }

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  // tick forces re-run after reset
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [anim, steps, tick]);

  return (
    <div className="terminal-window">
      {/* Title bar */}
      <div className="terminal-titlebar">
        <span className="terminal-dot" style={{ background: '#ff5f57' }} />
        <span className="terminal-dot" style={{ background: '#febc2e' }} />
        <span className="terminal-dot" style={{ background: '#28c840' }} />
        <span style={{
          marginLeft: 10,
          fontSize: '0.75rem',
          color: 'var(--color-text-muted)',
          fontFamily: 'monospace',
          userSelect: 'none',
        }}>
          {title}
        </span>
      </div>

      {/* Body */}
      <div className="terminal-body">
        {steps.map((step, si) => {
          if (si > anim.stepIndex) return null;

          const isCurrentStep = si === anim.stepIndex;
          const fullyTyped = !isCurrentStep || anim.typedChars >= step.command.length;
          const displayCmd = isCurrentStep
            ? step.command.slice(0, anim.typedChars)
            : step.command;
          const visibleLines = isCurrentStep ? anim.outputLines : step.output.length;
          const showCursor = isCurrentStep && anim.phase === 'typing';

          return (
            <div key={si} style={{ marginBottom: si < steps.length - 1 || !isCurrentStep ? '0.75rem' : 0 }}>
              {/* Command line */}
              <div>
                <span className="terminal-prompt">$ </span>
                <span className="terminal-cmd">{displayCmd}</span>
                {showCursor && <span className="terminal-cursor" />}
                {fullyTyped && !showCursor && anim.phase !== 'done' && isCurrentStep && (
                  <span className="terminal-cursor" />
                )}
              </div>
              {/* Output lines */}
              {step.output.slice(0, visibleLines).map((line, li) => (
                <div key={li} className="terminal-output">{line}</div>
              ))}
            </div>
          );
        })}

        {/* Replay button */}
        {anim.phase === 'done' && (
          <div style={{ marginTop: '1rem' }}>
            <button
              type="button"
              onClick={reset}
              style={{
                background: 'none',
                border: '1px solid var(--color-border)',
                borderRadius: 'var(--radius-sm)',
                color: 'var(--color-text-secondary)',
                fontSize: '0.8125rem',
                padding: '0.3rem 0.75rem',
                cursor: 'pointer',
                fontFamily: 'monospace',
                transition: 'border-color 0.12s, color 0.12s',
              }}
              onMouseEnter={e => {
                (e.currentTarget as HTMLButtonElement).style.borderColor = 'var(--color-accent)';
                (e.currentTarget as HTMLButtonElement).style.color = 'var(--color-accent)';
              }}
              onMouseLeave={e => {
                (e.currentTarget as HTMLButtonElement).style.borderColor = 'var(--color-border)';
                (e.currentTarget as HTMLButtonElement).style.color = 'var(--color-text-secondary)';
              }}
            >
              ▶ Replay
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default TerminalDemo;
