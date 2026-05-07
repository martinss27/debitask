function toDateKey(date) {
  return date.toISOString().slice(0, 10)
}

function isDayScheduled(daysMask, date) {
  const weekday = (date.getUTCDay() + 6) % 7
  return (daysMask & (1 << weekday)) !== 0
}

export default function HabitGrid({ habit, logs }) {
  const todayKey = toDateKey(new Date())
  const createdAtKey = habit.created_at ? habit.created_at.slice(0, 10) : todayKey

  const createdAt = new Date(createdAtKey)
  createdAt.setUTCHours(0, 0, 0, 0)

  const logMap = {}
  if (logs) {
    logs.forEach(l => { logMap[l.date.slice(0, 10)] = l.checked_in })
  }
  if (habit.checked_in_today) {
    logMap[todayKey] = true
  }

  // 7 rows (Mon-Sun) x 5 columns (weeks) = 35 cells
  // fill column by column, top to bottom
  const cells = []
  for (let row = 0; row < 7; row++) {
    for (let col = 0; col < 5; col++) {
      const index = col * 7 + row
      if (index >= 30) {
        // past day 30, render empty placeholder to keep grid shape
        cells.push(<div key={`empty-${row}-${col}`} className="habit-cell habit-cell--empty" />)
        continue
      }

      const d = new Date(createdAt)
      d.setUTCDate(createdAt.getUTCDate() + index)
      const key = toDateKey(d)
      const scheduled = isDayScheduled(habit.days, d)
      const checkedIn = logMap[key]
      const isToday = key === todayKey
      const isFuture = key > todayKey

      let state = 'skipped'
      if (!scheduled) state = 'skipped'
      else if (isFuture) state = 'future'
      else if (checkedIn === true) state = 'done'
      else if (isToday) state = 'today'
      else state = 'missed'

      cells.push(
        <div
          key={key}
          className={`habit-cell habit-cell--${state}`}
          title={`${key}${scheduled ? (checkedIn ? ' ✓' : isToday ? ' (today)' : ' ✗') : ''}`}
        />
      )
    }
  }

  return <div className="habit-grid">{cells}</div>
}
