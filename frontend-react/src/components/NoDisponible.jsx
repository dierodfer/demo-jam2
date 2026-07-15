export default function NoDisponible({ nombre }) {
  return (
    <div className="nodispo">
      <div className="nodispo-icon">🚧</div>
      <h2>Sección no disponible</h2>
      <p>«{nombre}» no está disponible en esta demo.</p>
    </div>
  );
}
