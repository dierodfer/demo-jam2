// Logo "n" con la elipse, recreado en SVG para la demo.
export default function Logo({ className }) {
  return (
    <svg className={className} viewBox="0 0 160 110" role="img" aria-label="Nunegal Consulting">
      <defs>
        <linearGradient id="swoosh" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor="#1a237e" />
          <stop offset="60%" stopColor="#9fa8da" />
          <stop offset="100%" stopColor="#cfd8dc" />
        </linearGradient>
      </defs>
      <ellipse
        cx="80" cy="55" rx="70" ry="34"
        fill="none" stroke="url(#swoosh)" strokeWidth="7"
        transform="rotate(-14 80 55)"
      />
      <text
        x="80" y="78" textAnchor="middle"
        fontSize="72" fontWeight="700" fontFamily="Georgia, serif"
        fill="#1a237e"
      >
        n
      </text>
    </svg>
  );
}
