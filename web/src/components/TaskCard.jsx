export default function TaskCard({ task, onComplete, onDelete }) {
  const deadline = new Date(task.deadline)
  const formattedDeadline = deadline.toLocaleDateString('en-US', {
    month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit',
  })

  return (
    <div className={`task-card task-card--${task.status}`}>
      <div className="task-card__body">
        <p className="task-card__title">{task.title}</p>
        {task.description && <p className="task-card__description">{task.description}</p>}
        <p className="task-card__deadline">Due {formattedDeadline}</p>
      </div>
      <div className="task-card__actions">
        {task.status === 'pending' && (
          <button className="btn-complete" onClick={() => onComplete(task.id)}>Complete</button>
        )}
        <button className="btn-delete" onClick={() => onDelete(task.id)}>Delete</button>
      </div>
    </div>
  )
}
